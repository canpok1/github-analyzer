package github

import (
	"context"
	"fmt"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
)

// ListIssues は指定リポジトリのIssue一覧を取得する。
func (c *Client) ListIssues(ctx context.Context, owner, repo string, opts domain.ListIssuesOptions) ([]entity.Issue, error) {
	// Numbers指定がある場合は個別取得
	if len(opts.Numbers) > 0 {
		return c.getIssuesByNumbers(ctx, owner, repo, opts.Numbers)
	}

	apiState := resolveIssueAPIState(opts.Status)

	ghOpts := &gh.IssueListByRepoOptions{
		State:     apiState,
		Sort:      "updated",
		Direction: "desc",
		ListOptions: gh.ListOptions{
			PerPage: 100,
		},
	}

	if opts.Since != nil {
		ghOpts.Since = *opts.Since
	}

	var allIssues []entity.Issue

	for {
		issues, resp, err := c.client.Issues.ListByRepo(ctx, owner, repo, ghOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}

		for _, issue := range issues {
			// go-github の Issues.ListByRepo はPRも返すため除外
			if issue.IsPullRequest() {
				continue
			}

			// Sinceフィルタ: updatedで降順ソート済みのため、Since以前なら以降も全て古い
			if opts.Since != nil && issue.GetUpdatedAt().Before(*opts.Since) {
				return allIssues, nil
			}

			allIssues = append(allIssues, convertIssue(issue))
		}

		if resp.NextPage == 0 {
			break
		}
		ghOpts.Page = resp.NextPage
	}

	return allIssues, nil
}

// getIssuesByNumbers は指定番号のIssueを個別に取得する。
func (c *Client) getIssuesByNumbers(ctx context.Context, owner, repo string, numbers []int) ([]entity.Issue, error) {
	issues := make([]entity.Issue, 0, len(numbers))

	for _, num := range numbers {
		issue, _, err := c.client.Issues.Get(ctx, owner, repo, num)
		if err != nil {
			return nil, fmt.Errorf("failed to get issue #%d: %w", num, err)
		}

		// PRの場合はスキップ
		if issue.IsPullRequest() {
			continue
		}

		issues = append(issues, convertIssue(issue))
	}

	return issues, nil
}

// resolveIssueAPIState はIssueStateをGitHub API用のstate値に変換する。
func resolveIssueAPIState(status entity.IssueState) string {
	switch status {
	case entity.IssueStateOpen:
		return "open"
	case entity.IssueStateClosed:
		return "closed"
	default:
		return "all"
	}
}

// convertIssue はgo-githubのIssueをドメインエンティティに変換する。
func convertIssue(issue *gh.Issue) entity.Issue {
	result := entity.Issue{
		Number: issue.GetNumber(),
		Title:  issue.GetTitle(),
		Author: issue.GetUser().GetLogin(),
		URL:    issue.GetHTMLURL(),
	}

	if issue.GetState() == "closed" {
		result.State = entity.IssueStateClosed
	} else {
		result.State = entity.IssueStateOpen
	}

	if issue.CreatedAt != nil {
		result.CreatedAt = issue.CreatedAt.Time
	}
	if issue.UpdatedAt != nil {
		result.UpdatedAt = issue.UpdatedAt.Time
	}
	if issue.ClosedAt != nil {
		t := issue.ClosedAt.Time
		result.ClosedAt = &t
	}

	for _, label := range issue.Labels {
		result.Labels = append(result.Labels, label.GetName())
	}

	return result
}

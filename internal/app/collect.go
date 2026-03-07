package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// CollectedData は収集したデータを集約する構造体。
type CollectedData struct {
	PullRequests []entity.PullRequest
	Issues       []entity.Issue
	Comments     map[int][]entity.Comment
	Timeline     map[int][]entity.TimelineEvent
}

// CollectData はQueryに応じてGitHub APIからデータを収集する。
func CollectData(ctx context.Context, gh domain.GitHubRepository, query entity.Query) (*CollectedData, error) {
	owner, repo, err := parseRepo(query.Repo)
	if err != nil {
		return nil, err
	}

	data := &CollectedData{
		Comments: make(map[int][]entity.Comment),
		Timeline: make(map[int][]entity.TimelineEvent),
	}

	if query.Since == nil && query.PR == nil && query.Issue == nil {
		return nil, fmt.Errorf("no target specified: specify --pr, --issue, --since, or --today")
	}

	if query.Since != nil {
		if err := collectByPeriod(ctx, gh, owner, repo, query, data); err != nil {
			return nil, err
		}
	}

	if query.Issue != nil {
		issues, err := gh.ListIssues(ctx, owner, repo, domain.ListIssuesOptions{
			Numbers: []int{*query.Issue},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}
		data.Issues = issues

		if err := collectDetails(ctx, gh, owner, repo, *query.Issue, false, data); err != nil {
			return nil, err
		}
	}

	if query.PR != nil {
		prs, err := gh.ListPullRequests(ctx, owner, repo, domain.ListPullRequestsOptions{
			Numbers: []int{*query.PR},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests: %w", err)
		}
		data.PullRequests = prs

		if err := collectDetails(ctx, gh, owner, repo, *query.PR, true, data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// collectDetails はコメントとタイムラインを収集する。isPRがtrueの場合はレビューコメントも取得する。
func collectDetails(ctx context.Context, gh domain.GitHubRepository, owner, repo string, number int, isPR bool, data *CollectedData) error {
	comments, err := gh.ListIssueComments(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to list issue comments for #%d: %w", number, err)
	}

	if isPR {
		reviewComments, err := gh.ListPullRequestComments(ctx, owner, repo, number)
		if err != nil {
			return fmt.Errorf("failed to list review comments for #%d: %w", number, err)
		}
		comments = append(comments, reviewComments...)
	}
	data.Comments[number] = comments

	timeline, err := gh.ListTimelineEvents(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to list timeline events for #%d: %w", number, err)
	}
	data.Timeline[number] = timeline

	return nil
}

// collectByPeriod は期間指定時にPR/Issueを一括取得し、各詳細を収集する。
func collectByPeriod(ctx context.Context, gh domain.GitHubRepository, owner, repo string, query entity.Query, data *CollectedData) error {
	prs, err := gh.ListPullRequests(ctx, owner, repo, domain.ListPullRequestsOptions{
		Since: query.Since,
	})
	if err != nil {
		return fmt.Errorf("failed to list pull requests: %w", err)
	}
	data.PullRequests = prs

	for _, pr := range prs {
		if err := collectDetails(ctx, gh, owner, repo, pr.Number, true, data); err != nil {
			return err
		}
	}

	issues, err := gh.ListIssues(ctx, owner, repo, domain.ListIssuesOptions{
		Since: query.Since,
	})
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}
	data.Issues = issues

	for _, issue := range issues {
		if err := collectDetails(ctx, gh, owner, repo, issue.Number, false, data); err != nil {
			return err
		}
	}

	return nil
}

// parseRepo は "owner/repo" 形式の文字列をownerとrepoに分割する。
func parseRepo(repoStr string) (string, string, error) {
	parts := strings.SplitN(repoStr, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo format %q: expected owner/repo", repoStr)
	}
	return parts[0], parts[1], nil
}

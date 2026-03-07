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

	if query.Issue != nil {
		issues, err := gh.ListIssues(ctx, owner, repo, domain.ListIssuesOptions{
			Numbers: []int{*query.Issue},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}
		data.Issues = issues

		if err := collectIssueDetails(ctx, gh, owner, repo, *query.Issue, data); err != nil {
			return nil, err
		}
	}

	if query.PR != nil {
		prs, err := gh.ListPullRequests(ctx, owner, repo, domain.ListPullRequestsOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests: %w", err)
		}
		// 指定PRのみフィルタ
		for _, pr := range prs {
			if pr.Number == *query.PR {
				data.PullRequests = append(data.PullRequests, pr)
			}
		}

		if err := collectPRDetails(ctx, gh, owner, repo, *query.PR, data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// collectPRDetails はPRのコメントとタイムラインを収集する。
func collectPRDetails(ctx context.Context, gh domain.GitHubRepository, owner, repo string, number int, data *CollectedData) error {
	issueComments, err := gh.ListIssueComments(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to list issue comments for #%d: %w", number, err)
	}

	reviewComments, err := gh.ListPullRequestComments(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to list review comments for #%d: %w", number, err)
	}

	var comments []entity.Comment
	comments = append(comments, issueComments...)
	comments = append(comments, reviewComments...)
	data.Comments[number] = comments

	timeline, err := gh.ListTimelineEvents(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to list timeline events for #%d: %w", number, err)
	}
	data.Timeline[number] = timeline

	return nil
}

// collectIssueDetails はIssueのコメントとタイムラインを収集する。
func collectIssueDetails(ctx context.Context, gh domain.GitHubRepository, owner, repo string, number int, data *CollectedData) error {
	comments, err := gh.ListIssueComments(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to list comments for issue #%d: %w", number, err)
	}
	data.Comments[number] = comments

	timeline, err := gh.ListTimelineEvents(ctx, owner, repo, number)
	if err != nil {
		return fmt.Errorf("failed to list timeline events for issue #%d: %w", number, err)
	}
	data.Timeline[number] = timeline

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

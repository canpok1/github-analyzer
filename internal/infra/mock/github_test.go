package mock

import (
	"context"
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain"
)

func TestGitHubRepository_ListPullRequests(t *testing.T) {
	repo := &GitHubRepository{}
	prs, err := repo.ListPullRequests(context.Background(), "owner", "repo", domain.ListPullRequestsOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(prs) == 0 {
		t.Error("should return at least one pull request")
	}
	for _, pr := range prs {
		if pr.Number == 0 {
			t.Error("PR number should not be zero")
		}
		if pr.Title == "" {
			t.Error("PR title should not be empty")
		}
		if pr.Author == "" {
			t.Error("PR author should not be empty")
		}
	}
}

func TestGitHubRepository_ListIssues(t *testing.T) {
	repo := &GitHubRepository{}
	issues, err := repo.ListIssues(context.Background(), "owner", "repo", domain.ListIssuesOptions{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Error("should return at least one issue")
	}
	for _, issue := range issues {
		if issue.Number == 0 {
			t.Error("Issue number should not be zero")
		}
		if issue.Title == "" {
			t.Error("Issue title should not be empty")
		}
		if issue.Author == "" {
			t.Error("Issue author should not be empty")
		}
	}
}

func TestGitHubRepository_ListIssueComments(t *testing.T) {
	repo := &GitHubRepository{}
	comments, err := repo.ListIssueComments(context.Background(), "owner", "repo", 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) == 0 {
		t.Error("should return at least one comment")
	}
}

func TestGitHubRepository_ListPullRequestComments(t *testing.T) {
	repo := &GitHubRepository{}
	comments, err := repo.ListPullRequestComments(context.Background(), "owner", "repo", 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) == 0 {
		t.Error("should return at least one comment")
	}
}

func TestGitHubRepository_ListTimelineEvents(t *testing.T) {
	repo := &GitHubRepository{}
	events, err := repo.ListTimelineEvents(context.Background(), "owner", "repo", 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(events) == 0 {
		t.Error("should return at least one event")
	}
}

package domain

import (
	"context"
	"testing"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// mockGitHubRepository はテスト用のモック実装。
type mockGitHubRepository struct{}

func (m *mockGitHubRepository) ListPullRequests(ctx context.Context, owner, repo string, opts ListPullRequestsOptions) ([]entity.PullRequest, error) {
	return nil, nil
}

func (m *mockGitHubRepository) ListIssues(ctx context.Context, owner, repo string, opts ListIssuesOptions) ([]entity.Issue, error) {
	return nil, nil
}

func TestGitHubRepository_InterfaceImplementation(t *testing.T) {
	var _ GitHubRepository = &mockGitHubRepository{}
}

func TestListIssuesOptions_HasExpectedFields(t *testing.T) {
	now := time.Now()
	opts := ListIssuesOptions{
		Since:   &now,
		Status:  entity.IssueStateOpen,
		Numbers: []int{1, 2, 3},
	}

	if opts.Since == nil || !opts.Since.Equal(now) {
		t.Error("Since field not set correctly")
	}
	if opts.Status != entity.IssueStateOpen {
		t.Errorf("Status = %q, want %q", opts.Status, entity.IssueStateOpen)
	}
	if len(opts.Numbers) != 3 || opts.Numbers[0] != 1 {
		t.Errorf("Numbers = %v, want [1 2 3]", opts.Numbers)
	}
}

func TestListPullRequestsOptions_HasExpectedFields(t *testing.T) {
	now := time.Now()
	opts := ListPullRequestsOptions{
		Since:  &now,
		Status: entity.PRStateOpen,
	}

	if opts.Since == nil || !opts.Since.Equal(now) {
		t.Error("Since field not set correctly")
	}
	if opts.Status != entity.PRStateOpen {
		t.Errorf("Status = %q, want %q", opts.Status, entity.PRStateOpen)
	}
}

func TestListPullRequestsOptions_SinceNil(t *testing.T) {
	opts := ListPullRequestsOptions{
		Status: entity.PRStateMerged,
	}

	if opts.Since != nil {
		t.Error("Since should be nil when not set")
	}
}

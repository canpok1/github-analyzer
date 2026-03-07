package domain

import (
	"context"
	"testing"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// TODO: 正常系: GitHubRepositoryインターフェースを実装できることの確認
// TODO: 正常系: ListPullRequestsOptionsの全フィールドが設定できる
// TODO: 境界値: ListPullRequestsOptionsのSinceがnilの場合

// mockGitHubRepository はテスト用のモック実装。
type mockGitHubRepository struct{}

func (m *mockGitHubRepository) ListPullRequests(ctx context.Context, owner, repo string, opts ListPullRequestsOptions) ([]entity.PullRequest, error) {
	return nil, nil
}

func TestGitHubRepository_InterfaceImplementation(t *testing.T) {
	var _ GitHubRepository = &mockGitHubRepository{}
}

func TestListPullRequestsOptions_HasExpectedFields(t *testing.T) {
	now := time.Now()
	opts := ListPullRequestsOptions{
		Since:  &now,
		Status: "open",
	}

	if opts.Since == nil || !opts.Since.Equal(now) {
		t.Error("Since field not set correctly")
	}
	if opts.Status != "open" {
		t.Errorf("Status = %q, want %q", opts.Status, "open")
	}
}

func TestListPullRequestsOptions_SinceNil(t *testing.T) {
	opts := ListPullRequestsOptions{
		Status: "merged",
	}

	if opts.Since != nil {
		t.Error("Since should be nil when not set")
	}
}

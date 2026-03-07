package app

import (
	"context"
	"testing"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// テストリスト: データ収集オーケストレーション
//
// TODO: 単一PR指定 -- PRのデータとコメント・タイムラインを取得すること
// TODO: 単一Issue指定 -- Issueのデータとコメント・タイムラインを取得すること
// TODO: 期間指定(since) -- 期間内の全PR/Issueとコメント・タイムラインを取得すること
// TODO: 期間指定(today) -- 当日の全PR/Issueとコメント・タイムラインを取得すること
// TODO: PR取得エラー時にエラーを返すこと
// TODO: Issue取得エラー時にエラーを返すこと
// TODO: コメント取得エラー時にエラーを返すこと
// TODO: タイムライン取得エラー時にエラーを返すこと
// TODO: 対象指定なしの場合エラーを返すこと

// mockGitHubRepository はテスト用のモック実装。
type mockGitHubRepository struct {
	listPullRequests        func(ctx context.Context, owner, repo string, opts domain.ListPullRequestsOptions) ([]entity.PullRequest, error)
	listIssues              func(ctx context.Context, owner, repo string, opts domain.ListIssuesOptions) ([]entity.Issue, error)
	listIssueComments       func(ctx context.Context, owner, repo string, number int) ([]entity.Comment, error)
	listPullRequestComments func(ctx context.Context, owner, repo string, number int) ([]entity.Comment, error)
	listTimelineEvents      func(ctx context.Context, owner, repo string, number int) ([]entity.TimelineEvent, error)
}

func (m *mockGitHubRepository) ListPullRequests(ctx context.Context, owner, repo string, opts domain.ListPullRequestsOptions) ([]entity.PullRequest, error) {
	if m.listPullRequests != nil {
		return m.listPullRequests(ctx, owner, repo, opts)
	}
	return nil, nil
}

func (m *mockGitHubRepository) ListIssues(ctx context.Context, owner, repo string, opts domain.ListIssuesOptions) ([]entity.Issue, error) {
	if m.listIssues != nil {
		return m.listIssues(ctx, owner, repo, opts)
	}
	return nil, nil
}

func (m *mockGitHubRepository) ListIssueComments(ctx context.Context, owner, repo string, number int) ([]entity.Comment, error) {
	if m.listIssueComments != nil {
		return m.listIssueComments(ctx, owner, repo, number)
	}
	return nil, nil
}

func (m *mockGitHubRepository) ListPullRequestComments(ctx context.Context, owner, repo string, number int) ([]entity.Comment, error) {
	if m.listPullRequestComments != nil {
		return m.listPullRequestComments(ctx, owner, repo, number)
	}
	return nil, nil
}

func (m *mockGitHubRepository) ListTimelineEvents(ctx context.Context, owner, repo string, number int) ([]entity.TimelineEvent, error) {
	if m.listTimelineEvents != nil {
		return m.listTimelineEvents(ctx, owner, repo, number)
	}
	return nil, nil
}

func TestCollectData_ReturnsCollectedData(t *testing.T) {
	mock := &mockGitHubRepository{}
	pr := 1
	query := entity.Query{
		PR:   &pr,
		Repo: "owner/repo",
	}

	result, err := CollectData(context.Background(), mock, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("result should not be nil")
	}
}

// suppress unused import warnings
var _ = time.Now

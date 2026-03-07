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
// DONE: 期間指定(today) -- sinceと同じロジックパス（cmd層でSinceにセットされる）
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

func TestCollectData_SinglePR(t *testing.T) {
	expectedPR := entity.PullRequest{
		Number: 1024,
		Title:  "Test PR",
		State:  entity.PRStateOpen,
		Author: "testuser",
	}
	expectedIssueComments := []entity.Comment{
		{ID: 1, Body: "issue comment", Type: entity.CommentTypeIssue},
	}
	expectedReviewComments := []entity.Comment{
		{ID: 2, Body: "review comment", Type: entity.CommentTypeReview},
	}
	expectedTimeline := []entity.TimelineEvent{
		{ID: 3, Event: "labeled", Label: "bug"},
	}

	mock := &mockGitHubRepository{
		listPullRequests: func(_ context.Context, owner, repo string, opts domain.ListPullRequestsOptions) ([]entity.PullRequest, error) {
			if owner != "myowner" || repo != "myrepo" {
				t.Errorf("unexpected owner/repo: %s/%s", owner, repo)
			}
			return []entity.PullRequest{expectedPR}, nil
		},
		listIssueComments: func(_ context.Context, _, _ string, number int) ([]entity.Comment, error) {
			if number != 1024 {
				t.Errorf("unexpected number: %d", number)
			}
			return expectedIssueComments, nil
		},
		listPullRequestComments: func(_ context.Context, _, _ string, number int) ([]entity.Comment, error) {
			if number != 1024 {
				t.Errorf("unexpected number: %d", number)
			}
			return expectedReviewComments, nil
		},
		listTimelineEvents: func(_ context.Context, _, _ string, number int) ([]entity.TimelineEvent, error) {
			if number != 1024 {
				t.Errorf("unexpected number: %d", number)
			}
			return expectedTimeline, nil
		},
	}

	pr := 1024
	query := entity.Query{
		PR:   &pr,
		Repo: "myowner/myrepo",
	}

	result, err := CollectData(context.Background(), mock, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.PullRequests) != 1 || result.PullRequests[0].Number != 1024 {
		t.Errorf("PullRequests = %v, want PR #1024", result.PullRequests)
	}

	comments, ok := result.Comments[1024]
	if !ok {
		t.Fatal("Comments[1024] not found")
	}
	if len(comments) != 2 {
		t.Errorf("len(Comments[1024]) = %d, want 2", len(comments))
	}

	timeline, ok := result.Timeline[1024]
	if !ok {
		t.Fatal("Timeline[1024] not found")
	}
	if len(timeline) != 1 {
		t.Errorf("len(Timeline[1024]) = %d, want 1", len(timeline))
	}
}

func TestCollectData_SingleIssue(t *testing.T) {
	expectedIssue := entity.Issue{
		Number: 42,
		Title:  "Test Issue",
		State:  entity.IssueStateOpen,
		Author: "testuser",
	}
	expectedComments := []entity.Comment{
		{ID: 10, Body: "comment on issue", Type: entity.CommentTypeIssue},
	}
	expectedTimeline := []entity.TimelineEvent{
		{ID: 20, Event: "assigned", Assignee: "dev1"},
	}

	mock := &mockGitHubRepository{
		listIssues: func(_ context.Context, owner, repo string, opts domain.ListIssuesOptions) ([]entity.Issue, error) {
			if owner != "myowner" || repo != "myrepo" {
				t.Errorf("unexpected owner/repo: %s/%s", owner, repo)
			}
			if len(opts.Numbers) != 1 || opts.Numbers[0] != 42 {
				t.Errorf("unexpected Numbers: %v, want [42]", opts.Numbers)
			}
			return []entity.Issue{expectedIssue}, nil
		},
		listIssueComments: func(_ context.Context, _, _ string, number int) ([]entity.Comment, error) {
			if number != 42 {
				t.Errorf("unexpected number: %d", number)
			}
			return expectedComments, nil
		},
		listTimelineEvents: func(_ context.Context, _, _ string, number int) ([]entity.TimelineEvent, error) {
			if number != 42 {
				t.Errorf("unexpected number: %d", number)
			}
			return expectedTimeline, nil
		},
	}

	issue := 42
	query := entity.Query{
		Issue: &issue,
		Repo:  "myowner/myrepo",
	}

	result, err := CollectData(context.Background(), mock, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Issues) != 1 || result.Issues[0].Number != 42 {
		t.Errorf("Issues = %v, want Issue #42", result.Issues)
	}

	comments, ok := result.Comments[42]
	if !ok {
		t.Fatal("Comments[42] not found")
	}
	if len(comments) != 1 {
		t.Errorf("len(Comments[42]) = %d, want 1", len(comments))
	}

	timeline, ok := result.Timeline[42]
	if !ok {
		t.Fatal("Timeline[42] not found")
	}
	if len(timeline) != 1 {
		t.Errorf("len(Timeline[42]) = %d, want 1", len(timeline))
	}
}

func TestCollectData_Since(t *testing.T) {
	since := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	pr1 := entity.PullRequest{Number: 10, Title: "PR 10"}
	pr2 := entity.PullRequest{Number: 11, Title: "PR 11"}
	issue1 := entity.Issue{Number: 20, Title: "Issue 20"}

	mock := &mockGitHubRepository{
		listPullRequests: func(_ context.Context, _, _ string, opts domain.ListPullRequestsOptions) ([]entity.PullRequest, error) {
			if opts.Since == nil || !opts.Since.Equal(since) {
				t.Errorf("unexpected Since: %v", opts.Since)
			}
			return []entity.PullRequest{pr1, pr2}, nil
		},
		listIssues: func(_ context.Context, _, _ string, opts domain.ListIssuesOptions) ([]entity.Issue, error) {
			if opts.Since == nil || !opts.Since.Equal(since) {
				t.Errorf("unexpected Since: %v", opts.Since)
			}
			return []entity.Issue{issue1}, nil
		},
		listIssueComments: func(_ context.Context, _, _ string, _ int) ([]entity.Comment, error) {
			return []entity.Comment{{ID: 100, Body: "c"}}, nil
		},
		listPullRequestComments: func(_ context.Context, _, _ string, _ int) ([]entity.Comment, error) {
			return nil, nil
		},
		listTimelineEvents: func(_ context.Context, _, _ string, _ int) ([]entity.TimelineEvent, error) {
			return []entity.TimelineEvent{{ID: 200, Event: "e"}}, nil
		},
	}

	query := entity.Query{
		Since: &since,
		Repo:  "owner/repo",
	}

	result, err := CollectData(context.Background(), mock, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.PullRequests) != 2 {
		t.Errorf("len(PullRequests) = %d, want 2", len(result.PullRequests))
	}
	if len(result.Issues) != 1 {
		t.Errorf("len(Issues) = %d, want 1", len(result.Issues))
	}

	// PR 10, 11 と Issue 20 のコメント・タイムラインが取得されていること
	for _, num := range []int{10, 11, 20} {
		if _, ok := result.Comments[num]; !ok {
			t.Errorf("Comments[%d] not found", num)
		}
		if _, ok := result.Timeline[num]; !ok {
			t.Errorf("Timeline[%d] not found", num)
		}
	}
}

func TestCollectData_NoTargetSpecified(t *testing.T) {
	mock := &mockGitHubRepository{}
	query := entity.Query{
		Repo: "owner/repo",
	}

	_, err := CollectData(context.Background(), mock, query)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

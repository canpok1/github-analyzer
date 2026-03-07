package domain

import (
	"context"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// ListPullRequestsOptions はPR一覧取得のオプション。
type ListPullRequestsOptions struct {
	Since   *time.Time
	Status  entity.PRState
	Numbers []int
}

// ListIssuesOptions はIssue一覧取得のオプション。
type ListIssuesOptions struct {
	Since   *time.Time
	Status  entity.IssueState
	Numbers []int
}

// GitHubRepository はGitHub APIへのアクセスを抽象化するインターフェース。
type GitHubRepository interface {
	// ListPullRequests は指定リポジトリのPR一覧を取得する。
	ListPullRequests(ctx context.Context, owner, repo string, opts ListPullRequestsOptions) ([]entity.PullRequest, error)
	// ListIssues は指定リポジトリのIssue一覧を取得する。
	ListIssues(ctx context.Context, owner, repo string, opts ListIssuesOptions) ([]entity.Issue, error)
	// ListIssueComments は指定Issue/PRの一般コメント一覧を取得する。
	ListIssueComments(ctx context.Context, owner, repo string, number int) ([]entity.Comment, error)
	// ListPullRequestComments は指定PRのレビューコメント一覧を取得する。
	ListPullRequestComments(ctx context.Context, owner, repo string, number int) ([]entity.Comment, error)
	// ListTimelineEvents は指定Issue/PRのタイムラインイベント一覧を取得する。
	ListTimelineEvents(ctx context.Context, owner, repo string, number int) ([]entity.TimelineEvent, error)
}

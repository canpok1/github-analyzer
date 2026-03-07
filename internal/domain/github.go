package domain

import (
	"context"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// ListPullRequestsOptions はPR一覧取得のオプション。
type ListPullRequestsOptions struct {
	Since  *time.Time
	Status entity.PRState
}

// GitHubRepository はGitHub APIへのアクセスを抽象化するインターフェース。
type GitHubRepository interface {
	// ListPullRequests は指定リポジトリのPR一覧を取得する。
	ListPullRequests(ctx context.Context, owner, repo string, opts ListPullRequestsOptions) ([]entity.PullRequest, error)
}

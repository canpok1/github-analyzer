package github

import (
	"context"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
)

// Client はGitHub APIクライアント。
type Client struct {
	client *gh.Client
}

// NewClient は新しいGitHub APIクライアントを生成する。
// tokenが空の場合は認証なしのクライアントを返す。
func NewClient(token string) *Client {
	var client *gh.Client
	if token != "" {
		client = gh.NewClient(nil).WithAuthToken(token)
	} else {
		client = gh.NewClient(nil)
	}
	return &Client{client: client}
}

// ListPullRequests は指定リポジトリのPR一覧を取得する。
func (c *Client) ListPullRequests(_ context.Context, _, _ string, _ domain.ListPullRequestsOptions) ([]entity.PullRequest, error) {
	return nil, nil
}

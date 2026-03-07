package github

import (
	gh "github.com/google/go-github/v68/github"
)

// Client はGitHub APIクライアント。
type Client struct {
	client *gh.Client
}

// NewClient は新しいGitHub APIクライアントを生成する。
// tokenが空の場合は認証なしのクライアントを返す。
func NewClient(token string) *Client {
	client := gh.NewClient(nil)
	if token != "" {
		client = client.WithAuthToken(token)
	}
	return &Client{client: client}
}

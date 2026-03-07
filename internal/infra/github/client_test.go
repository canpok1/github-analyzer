package github

import (
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain"
)

// TODO: 正常系: NewClientでGitHubRepositoryインターフェースを満たす
// TODO: 正常系: トークンありでクライアント生成
// TODO: 正常系: トークンなしでクライアント生成

func TestNewClient_ImplementsGitHubRepository(t *testing.T) {
	c := NewClient("dummy-token")
	var _ domain.GitHubRepository = c
}

func TestNewClient_WithToken(t *testing.T) {
	c := NewClient("test-token")
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
}

func TestNewClient_EmptyToken(t *testing.T) {
	c := NewClient("")
	if c == nil {
		t.Fatal("NewClient with empty token returned nil")
	}
}

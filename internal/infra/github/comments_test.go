package github

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
)

// === テストリスト ===
// DONE: 正常系: Issueコメント一覧を取得できる
// DONE: 正常系: PRレビューコメント一覧を取得できる
// DONE: 正常系: 空のコメント一覧を返す
// DONE: 異常系: IssueコメントAPIエラー時にエラーを返す
// DONE: 異常系: PRレビューコメントAPIエラー時にエラーを返す

func TestListIssueComments_ReturnsComments(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	updatedAt := now.Add(1 * time.Hour)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/issues/1/comments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		comments := []*gh.IssueComment{
			{
				ID:        gh.Ptr(int64(100)),
				Body:      gh.Ptr("テストコメント"),
				User:      &gh.User{Login: gh.Ptr("testuser")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: updatedAt},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/issues/1#issuecomment-100"),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	})
	defer server.Close()

	results, err := c.ListIssueComments(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("got %d comments, want 1", len(results))
	}

	comment := results[0]
	if comment.ID != 100 {
		t.Errorf("ID = %d, want 100", comment.ID)
	}
	if comment.Body != "テストコメント" {
		t.Errorf("Body = %q, want %q", comment.Body, "テストコメント")
	}
	if comment.Author != "testuser" {
		t.Errorf("Author = %q, want %q", comment.Author, "testuser")
	}
	if !comment.CreatedAt.Equal(now) {
		t.Error("CreatedAt not set correctly")
	}
	if comment.UpdatedAt == nil || !comment.UpdatedAt.Equal(updatedAt) {
		t.Error("UpdatedAt not set correctly")
	}
	if comment.Type != entity.CommentTypeIssue {
		t.Errorf("Type = %q, want %q", comment.Type, entity.CommentTypeIssue)
	}
	if comment.Path != "" {
		t.Errorf("Path = %q, want empty", comment.Path)
	}
	if comment.URL != "https://github.com/owner/repo/issues/1#issuecomment-100" {
		t.Errorf("URL = %q, want %q", comment.URL, "https://github.com/owner/repo/issues/1#issuecomment-100")
	}
}

func TestListPullRequestComments_ReturnsReviewComments(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	updatedAt := now.Add(1 * time.Hour)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/pulls/5/comments" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		comments := []*gh.PullRequestComment{
			{
				ID:        gh.Ptr(int64(200)),
				Body:      gh.Ptr("レビューコメント"),
				User:      &gh.User{Login: gh.Ptr("reviewer")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: updatedAt},
				Path:      gh.Ptr("main.go"),
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/pull/5#discussion_r200"),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	})
	defer server.Close()

	results, err := c.ListPullRequestComments(context.Background(), "owner", "repo", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("got %d comments, want 1", len(results))
	}

	comment := results[0]
	if comment.ID != 200 {
		t.Errorf("ID = %d, want 200", comment.ID)
	}
	if comment.Body != "レビューコメント" {
		t.Errorf("Body = %q, want %q", comment.Body, "レビューコメント")
	}
	if comment.Author != "reviewer" {
		t.Errorf("Author = %q, want %q", comment.Author, "reviewer")
	}
	if comment.Type != entity.CommentTypeReview {
		t.Errorf("Type = %q, want %q", comment.Type, entity.CommentTypeReview)
	}
	if comment.Path != "main.go" {
		t.Errorf("Path = %q, want %q", comment.Path, "main.go")
	}
	if comment.URL != "https://github.com/owner/repo/pull/5#discussion_r200" {
		t.Errorf("URL = %q, want %q", comment.URL, "https://github.com/owner/repo/pull/5#discussion_r200")
	}
}

func TestListIssueComments_Empty(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*gh.IssueComment{})
	})
	defer server.Close()

	results, err := c.ListIssueComments(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("got %d comments, want 0", len(results))
	}
}

func TestListIssueComments_APIError(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	_, err := c.ListIssueComments(context.Background(), "owner", "repo", 1)
	if err == nil {
		t.Error("expected error for API failure")
	}
}

func TestListPullRequestComments_APIError(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	_, err := c.ListPullRequestComments(context.Background(), "owner", "repo", 1)
	if err == nil {
		t.Error("expected error for API failure")
	}
}

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain"
	gh "github.com/google/go-github/v68/github"
)

// TODO: 正常系: PR一覧を取得できる
// TODO: 正常系: Sinceフィルタが適用される
// TODO: 正常系: Statusフィルタ (open) が適用される
// TODO: 正常系: Statusフィルタ (merged) が適用される - closedで取得しmergedAtでフィルタ
// TODO: 正常系: 空のPR一覧を取得
// TODO: 異常系: APIエラー時にエラーを返す

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := gh.NewClient(nil)
	url := server.URL + "/"
	client.BaseURL, _ = client.BaseURL.Parse(url)
	return &Client{client: client}, server
}

func TestListPullRequests_ReturnsPRs(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	mergedAt := now.Add(-1 * time.Hour)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/pulls" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		prs := []*gh.PullRequest{
			{
				Number:    gh.Ptr(1),
				Title:     gh.Ptr("Test PR"),
				State:     gh.Ptr("closed"),
				User:      &gh.User{Login: gh.Ptr("testuser")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: now},
				MergedAt:  &gh.Timestamp{Time: mergedAt},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/pull/1"),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prs)
	})
	defer server.Close()

	results, err := c.ListPullRequests(context.Background(), "owner", "repo", domain.ListPullRequestsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("got %d PRs, want 1", len(results))
	}

	pr := results[0]
	if pr.Number != 1 {
		t.Errorf("Number = %d, want 1", pr.Number)
	}
	if pr.Title != "Test PR" {
		t.Errorf("Title = %q, want %q", pr.Title, "Test PR")
	}
	if pr.Author != "testuser" {
		t.Errorf("Author = %q, want %q", pr.Author, "testuser")
	}
	if pr.MergedAt == nil {
		t.Error("MergedAt should not be nil")
	}
	if pr.URL != "https://github.com/owner/repo/pull/1" {
		t.Errorf("URL = %q, want %q", pr.URL, "https://github.com/owner/repo/pull/1")
	}
}

func TestListPullRequests_StatusOpen(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		if state != "open" {
			t.Errorf("state query param = %q, want %q", state, "open")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*gh.PullRequest{})
	})
	defer server.Close()

	_, err := c.ListPullRequests(context.Background(), "owner", "repo", domain.ListPullRequestsOptions{
		Status: "open",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListPullRequests_StatusMerged(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	mergedAt := now.Add(-1 * time.Hour)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// merged は GitHub API では closed + mergedAt != nil
		state := r.URL.Query().Get("state")
		if state != "closed" {
			t.Errorf("state query param = %q, want %q for merged filter", state, "closed")
		}

		prs := []*gh.PullRequest{
			{
				Number:    gh.Ptr(1),
				Title:     gh.Ptr("Merged PR"),
				State:     gh.Ptr("closed"),
				User:      &gh.User{Login: gh.Ptr("user1")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: now},
				MergedAt:  &gh.Timestamp{Time: mergedAt},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/pull/1"),
			},
			{
				Number:    gh.Ptr(2),
				Title:     gh.Ptr("Closed PR (not merged)"),
				State:     gh.Ptr("closed"),
				User:      &gh.User{Login: gh.Ptr("user2")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: now},
				MergedAt:  nil,
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/pull/2"),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prs)
	})
	defer server.Close()

	results, err := c.ListPullRequests(context.Background(), "owner", "repo", domain.ListPullRequestsOptions{
		Status: "merged",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// mergedフィルタではマージ済みのPRのみ返る
	if len(results) != 1 {
		t.Fatalf("got %d PRs, want 1 (only merged)", len(results))
	}
	if results[0].Number != 1 {
		t.Errorf("expected merged PR #1, got #%d", results[0].Number)
	}
}

func TestListPullRequests_SinceFilter(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	since := now.Add(-24 * time.Hour)
	oldTime := now.Add(-72 * time.Hour)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		prs := []*gh.PullRequest{
			{
				Number:    gh.Ptr(1),
				Title:     gh.Ptr("Recent PR"),
				State:     gh.Ptr("open"),
				User:      &gh.User{Login: gh.Ptr("user1")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: now},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/pull/1"),
			},
			{
				Number:    gh.Ptr(2),
				Title:     gh.Ptr("Old PR"),
				State:     gh.Ptr("open"),
				User:      &gh.User{Login: gh.Ptr("user2")},
				CreatedAt: &gh.Timestamp{Time: oldTime},
				UpdatedAt: &gh.Timestamp{Time: oldTime},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/pull/2"),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prs)
	})
	defer server.Close()

	results, err := c.ListPullRequests(context.Background(), "owner", "repo", domain.ListPullRequestsOptions{
		Since: &since,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Sinceフィルタにより最近のPRのみ返る
	if len(results) != 1 {
		t.Fatalf("got %d PRs, want 1 (only recent)", len(results))
	}
	if results[0].Number != 1 {
		t.Errorf("expected recent PR #1, got #%d", results[0].Number)
	}
}

func TestListPullRequests_Empty(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*gh.PullRequest{})
	})
	defer server.Close()

	results, err := c.ListPullRequests(context.Background(), "owner", "repo", domain.ListPullRequestsOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("got %d PRs, want 0", len(results))
	}
}

func TestListPullRequests_APIError(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	_, err := c.ListPullRequests(context.Background(), "owner", "repo", domain.ListPullRequestsOptions{})
	if err == nil {
		t.Error("expected error for API failure")
	}
}

package github

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
)

// === テストリスト ===
// DONE: 正常系: Issue一覧を取得できる（PRを除外する）
// DONE: 正常系: ステータスフィルタ（open）が適用される
// DONE: 正常系: Sinceフィルタで古いIssueが除外される
// DONE: 正常系: 空の結果を返す
// DONE: 異常系: APIエラー時にエラーを返す
// DONE: 正常系: Numbers指定で個別Issueを取得する
// DONE: 正常系: Numbers指定でAPIエラー時にエラーを返す

func TestListIssues_ReturnsIssuesExcludingPRs(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	closedAt := now.Add(-1 * time.Hour)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/issues" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		issues := []*gh.Issue{
			{
				Number:    gh.Ptr(1),
				Title:     gh.Ptr("Test Issue"),
				State:     gh.Ptr("closed"),
				User:      &gh.User{Login: gh.Ptr("testuser")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: now},
				ClosedAt:  &gh.Timestamp{Time: closedAt},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/issues/1"),
				Labels: []*gh.Label{
					{Name: gh.Ptr("bug")},
					{Name: gh.Ptr("help wanted")},
				},
			},
			{
				// PRはIssuesAPIでも返るので除外されるべき
				Number:           gh.Ptr(2),
				Title:            gh.Ptr("Test PR via Issues API"),
				State:            gh.Ptr("open"),
				User:             &gh.User{Login: gh.Ptr("pruser")},
				CreatedAt:        &gh.Timestamp{Time: now},
				UpdatedAt:        &gh.Timestamp{Time: now},
				HTMLURL:          gh.Ptr("https://github.com/owner/repo/pull/2"),
				PullRequestLinks: &gh.PullRequestLinks{URL: gh.Ptr("https://api.github.com/repos/owner/repo/pulls/2")},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(issues)
	})
	defer server.Close()

	results, err := c.ListIssues(context.Background(), "owner", "repo", domain.ListIssuesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("got %d issues, want 1 (PR should be excluded)", len(results))
	}

	issue := results[0]
	if issue.Number != 1 {
		t.Errorf("Number = %d, want 1", issue.Number)
	}
	if issue.Title != "Test Issue" {
		t.Errorf("Title = %q, want %q", issue.Title, "Test Issue")
	}
	if issue.State != entity.IssueStateClosed {
		t.Errorf("State = %q, want %q", issue.State, entity.IssueStateClosed)
	}
	if issue.Author != "testuser" {
		t.Errorf("Author = %q, want %q", issue.Author, "testuser")
	}
	if issue.ClosedAt == nil {
		t.Error("ClosedAt should not be nil")
	}
	if issue.URL != "https://github.com/owner/repo/issues/1" {
		t.Errorf("URL = %q, want %q", issue.URL, "https://github.com/owner/repo/issues/1")
	}
	if len(issue.Labels) != 2 || issue.Labels[0] != "bug" || issue.Labels[1] != "help wanted" {
		t.Errorf("Labels = %v, want [bug, help wanted]", issue.Labels)
	}
}

func TestListIssues_StatusOpen(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		if state != "open" {
			t.Errorf("state query param = %q, want %q", state, "open")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*gh.Issue{})
	})
	defer server.Close()

	_, err := c.ListIssues(context.Background(), "owner", "repo", domain.ListIssuesOptions{
		Status: entity.IssueStateOpen,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListIssues_SinceFilter(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	since := now.Add(-24 * time.Hour)
	oldTime := now.Add(-72 * time.Hour)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		issues := []*gh.Issue{
			{
				Number:    gh.Ptr(1),
				Title:     gh.Ptr("Recent Issue"),
				State:     gh.Ptr("open"),
				User:      &gh.User{Login: gh.Ptr("user1")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: now},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/issues/1"),
			},
			{
				Number:    gh.Ptr(2),
				Title:     gh.Ptr("Old Issue"),
				State:     gh.Ptr("open"),
				User:      &gh.User{Login: gh.Ptr("user2")},
				CreatedAt: &gh.Timestamp{Time: oldTime},
				UpdatedAt: &gh.Timestamp{Time: oldTime},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/issues/2"),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(issues)
	})
	defer server.Close()

	results, err := c.ListIssues(context.Background(), "owner", "repo", domain.ListIssuesOptions{
		Since: &since,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("got %d issues, want 1 (only recent)", len(results))
	}
	if results[0].Number != 1 {
		t.Errorf("expected recent issue #1, got #%d", results[0].Number)
	}
}

func TestListIssues_Empty(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*gh.Issue{})
	})
	defer server.Close()

	results, err := c.ListIssues(context.Background(), "owner", "repo", domain.ListIssuesOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("got %d issues, want 0", len(results))
	}
}

func TestListIssues_APIError(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	_, err := c.ListIssues(context.Background(), "owner", "repo", domain.ListIssuesOptions{})
	if err == nil {
		t.Error("expected error for API failure")
	}
}

func TestListIssues_ByNumbers(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Issues.Get は /repos/owner/repo/issues/{number} にアクセスする
		switch r.URL.Path {
		case "/repos/owner/repo/issues/10":
			issue := &gh.Issue{
				Number:    gh.Ptr(10),
				Title:     gh.Ptr("Issue Ten"),
				State:     gh.Ptr("open"),
				User:      &gh.User{Login: gh.Ptr("user1")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: now},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/issues/10"),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(issue)
		case "/repos/owner/repo/issues/20":
			issue := &gh.Issue{
				Number:    gh.Ptr(20),
				Title:     gh.Ptr("Issue Twenty"),
				State:     gh.Ptr("closed"),
				User:      &gh.User{Login: gh.Ptr("user2")},
				CreatedAt: &gh.Timestamp{Time: now},
				UpdatedAt: &gh.Timestamp{Time: now},
				HTMLURL:   gh.Ptr("https://github.com/owner/repo/issues/20"),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(issue)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	defer server.Close()

	results, err := c.ListIssues(context.Background(), "owner", "repo", domain.ListIssuesOptions{
		Numbers: []int{10, 20},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("got %d issues, want 2", len(results))
	}
	if results[0].Number != 10 {
		t.Errorf("first issue Number = %d, want 10", results[0].Number)
	}
	if results[1].Number != 20 {
		t.Errorf("second issue Number = %d, want 20", results[1].Number)
	}
}

func TestListIssues_ByNumbers_APIError(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	defer server.Close()

	_, err := c.ListIssues(context.Background(), "owner", "repo", domain.ListIssuesOptions{
		Numbers: []int{999},
	})
	if err == nil {
		t.Error("expected error for API failure")
	}
}

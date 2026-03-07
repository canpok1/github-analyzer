package github

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	gh "github.com/google/go-github/v68/github"
)

// === テストリスト ===
// DONE: 正常系: タイムラインイベント一覧を取得できる
// DONE: 正常系: 空のタイムラインイベント一覧を返す
// DONE: 正常系: ラベル・アサインイベントが変換される
// DONE: 異常系: APIエラー時にエラーを返す

func TestListTimelineEvents_ReturnsEvents(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/issues/1/timeline" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// mockaccept header check
		accept := r.Header.Get("Accept")
		if accept == "" {
			t.Error("Accept header should be set")
		}

		events := []*gh.Timeline{
			{
				ID:        gh.Ptr(int64(1000)),
				Event:     gh.Ptr("closed"),
				Actor:     &gh.User{Login: gh.Ptr("closer")},
				CreatedAt: &gh.Timestamp{Time: now},
				CommitID:  gh.Ptr("abc123def"),
				URL:       gh.Ptr("https://api.github.com/repos/owner/repo/issues/events/1000"),
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})
	defer server.Close()

	results, err := c.ListTimelineEvents(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("got %d events, want 1", len(results))
	}

	event := results[0]
	if event.ID != 1000 {
		t.Errorf("ID = %d, want 1000", event.ID)
	}
	if event.Event != "closed" {
		t.Errorf("Event = %q, want %q", event.Event, "closed")
	}
	if event.Actor != "closer" {
		t.Errorf("Actor = %q, want %q", event.Actor, "closer")
	}
	if !event.CreatedAt.Equal(now) {
		t.Error("CreatedAt not set correctly")
	}
	if event.CommitID != "abc123def" {
		t.Errorf("CommitID = %q, want %q", event.CommitID, "abc123def")
	}
}

func TestListTimelineEvents_Empty(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]*gh.Timeline{})
	})
	defer server.Close()

	results, err := c.ListTimelineEvents(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("got %d events, want 0", len(results))
	}
}

func TestListTimelineEvents_LabelAndAssignEvents(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		events := []*gh.Timeline{
			{
				ID:        gh.Ptr(int64(2000)),
				Event:     gh.Ptr("labeled"),
				Actor:     &gh.User{Login: gh.Ptr("labeler")},
				CreatedAt: &gh.Timestamp{Time: now},
				Label:     &gh.Label{Name: gh.Ptr("enhancement")},
			},
			{
				ID:        gh.Ptr(int64(2001)),
				Event:     gh.Ptr("assigned"),
				Actor:     &gh.User{Login: gh.Ptr("assigner")},
				CreatedAt: &gh.Timestamp{Time: now},
				Assignee:  &gh.User{Login: gh.Ptr("assignee-user")},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})
	defer server.Close()

	results, err := c.ListTimelineEvents(context.Background(), "owner", "repo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("got %d events, want 2", len(results))
	}

	// ラベルイベント
	if results[0].Event != "labeled" {
		t.Errorf("Event = %q, want %q", results[0].Event, "labeled")
	}
	if results[0].Label != "enhancement" {
		t.Errorf("Label = %q, want %q", results[0].Label, "enhancement")
	}

	// アサインイベント
	if results[1].Event != "assigned" {
		t.Errorf("Event = %q, want %q", results[1].Event, "assigned")
	}
	if results[1].Assignee != "assignee-user" {
		t.Errorf("Assignee = %q, want %q", results[1].Assignee, "assignee-user")
	}
}

func TestListTimelineEvents_APIError(t *testing.T) {
	c, server := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	defer server.Close()

	_, err := c.ListTimelineEvents(context.Background(), "owner", "repo", 1)
	if err == nil {
		t.Error("expected error for API failure")
	}
}

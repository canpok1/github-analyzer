package entity

import (
	"testing"
	"time"
)

// === テストリスト ===
// DONE: TimelineEvent構造体に期待フィールドが設定できる
// DONE: TimelineEvent の Label が空の場合
// DONE: TimelineEvent の Assignee が空の場合

func TestTimelineEvent_HasExpectedFields(t *testing.T) {
	now := time.Now()
	event := TimelineEvent{
		ID:        456,
		Event:     "labeled",
		Actor:     "testuser",
		CreatedAt: now,
		Label:     "bug",
		Assignee:  "assigneduser",
		CommitID:  "abc123",
		URL:       "https://github.com/owner/repo/issues/1#event-456",
	}

	if event.ID != 456 {
		t.Errorf("ID = %d, want 456", event.ID)
	}
	if event.Event != "labeled" {
		t.Errorf("Event = %q, want %q", event.Event, "labeled")
	}
	if event.Actor != "testuser" {
		t.Errorf("Actor = %q, want %q", event.Actor, "testuser")
	}
	if !event.CreatedAt.Equal(now) {
		t.Error("CreatedAt not set correctly")
	}
	if event.Label != "bug" {
		t.Errorf("Label = %q, want %q", event.Label, "bug")
	}
	if event.Assignee != "assigneduser" {
		t.Errorf("Assignee = %q, want %q", event.Assignee, "assigneduser")
	}
	if event.CommitID != "abc123" {
		t.Errorf("CommitID = %q, want %q", event.CommitID, "abc123")
	}
	if event.URL != "https://github.com/owner/repo/issues/1#event-456" {
		t.Errorf("URL = %q, want %q", event.URL, "https://github.com/owner/repo/issues/1#event-456")
	}
}

func TestTimelineEvent_EmptyLabel(t *testing.T) {
	event := TimelineEvent{
		Event: "closed",
		Actor: "user",
		Label: "",
	}

	if event.Label != "" {
		t.Errorf("Label = %q, want empty string", event.Label)
	}
}

func TestTimelineEvent_EmptyAssignee(t *testing.T) {
	event := TimelineEvent{
		Event:    "merged",
		Actor:    "user",
		Assignee: "",
	}

	if event.Assignee != "" {
		t.Errorf("Assignee = %q, want empty string", event.Assignee)
	}
}

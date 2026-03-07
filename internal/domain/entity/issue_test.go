package entity

import (
	"testing"
	"time"
)

// === テストリスト ===
// DONE: IssueState の値が正しい（open, closed）
// DONE: Issue構造体に期待フィールドが設定できる
// DONE: Issue の ClosedAt が nil の場合（オープンIssue）
// DONE: Issue の Labels が空の場合

func TestIssueState_Values(t *testing.T) {
	tests := []struct {
		state IssueState
		want  string
	}{
		{IssueStateOpen, "open"},
		{IssueStateClosed, "closed"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if string(tt.state) != tt.want {
				t.Errorf("IssueState = %q, want %q", tt.state, tt.want)
			}
		})
	}
}

func TestIssue_HasExpectedFields(t *testing.T) {
	now := time.Now()
	closedAt := now.Add(-1 * time.Hour)
	issue := Issue{
		Number:    42,
		Title:     "テストIssue",
		State:     IssueStateClosed,
		Author:    "testuser",
		CreatedAt: now,
		UpdatedAt: now,
		ClosedAt:  &closedAt,
		URL:       "https://github.com/owner/repo/issues/42",
		Labels:    []string{"bug", "high-priority"},
	}

	if issue.Number != 42 {
		t.Errorf("Number = %d, want 42", issue.Number)
	}
	if issue.Title != "テストIssue" {
		t.Errorf("Title = %q, want %q", issue.Title, "テストIssue")
	}
	if issue.State != IssueStateClosed {
		t.Errorf("State = %q, want %q", issue.State, IssueStateClosed)
	}
	if issue.Author != "testuser" {
		t.Errorf("Author = %q, want %q", issue.Author, "testuser")
	}
	if !issue.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt not set correctly")
	}
	if !issue.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt not set correctly")
	}
	if issue.ClosedAt == nil || !issue.ClosedAt.Equal(closedAt) {
		t.Errorf("ClosedAt not set correctly")
	}
	if issue.URL != "https://github.com/owner/repo/issues/42" {
		t.Errorf("URL = %q, want %q", issue.URL, "https://github.com/owner/repo/issues/42")
	}
	if len(issue.Labels) != 2 || issue.Labels[0] != "bug" || issue.Labels[1] != "high-priority" {
		t.Errorf("Labels = %v, want [bug high-priority]", issue.Labels)
	}
}

func TestIssue_ClosedAtNil(t *testing.T) {
	issue := Issue{
		Number:    1,
		Title:     "オープンIssue",
		State:     IssueStateOpen,
		Author:    "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ClosedAt:  nil,
		URL:       "https://github.com/owner/repo/issues/1",
	}

	if issue.ClosedAt != nil {
		t.Error("ClosedAt should be nil for open issue")
	}
}

func TestIssue_EmptyLabels(t *testing.T) {
	issue := Issue{
		Number: 1,
		Labels: []string{},
	}

	if len(issue.Labels) != 0 {
		t.Errorf("Labels should be empty, got %v", issue.Labels)
	}
}

package entity

import (
	"testing"
	"time"
)

// TODO: 正常系: PullRequestの全フィールドが正しく設定される
// TODO: 正常系: PRStatusの文字列値が正しい
// TODO: 境界値: MergedAtがnilの場合（未マージPR）

func TestPullRequest_HasExpectedFields(t *testing.T) {
	now := time.Now()
	mergedAt := now.Add(-1 * time.Hour)
	pr := PullRequest{
		Number:    1,
		Title:     "テストPR",
		State:     PRStateOpen,
		Author:    "testuser",
		CreatedAt: now,
		UpdatedAt: now,
		MergedAt:  &mergedAt,
		URL:       "https://github.com/owner/repo/pull/1",
	}

	if pr.Number != 1 {
		t.Errorf("Number = %d, want 1", pr.Number)
	}
	if pr.Title != "テストPR" {
		t.Errorf("Title = %q, want %q", pr.Title, "テストPR")
	}
	if pr.State != PRStateOpen {
		t.Errorf("State = %q, want %q", pr.State, PRStateOpen)
	}
	if pr.Author != "testuser" {
		t.Errorf("Author = %q, want %q", pr.Author, "testuser")
	}
	if !pr.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt not set correctly")
	}
	if !pr.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt not set correctly")
	}
	if pr.MergedAt == nil || !pr.MergedAt.Equal(mergedAt) {
		t.Errorf("MergedAt not set correctly")
	}
	if pr.URL != "https://github.com/owner/repo/pull/1" {
		t.Errorf("URL = %q, want %q", pr.URL, "https://github.com/owner/repo/pull/1")
	}
}

func TestPRState_Values(t *testing.T) {
	tests := []struct {
		state PRState
		want  string
	}{
		{PRStateOpen, "open"},
		{PRStateClosed, "closed"},
		{PRStateMerged, "merged"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if string(tt.state) != tt.want {
				t.Errorf("PRState = %q, want %q", tt.state, tt.want)
			}
		})
	}
}

func TestPullRequest_MergedAtNil(t *testing.T) {
	pr := PullRequest{
		Number:    2,
		Title:     "未マージPR",
		State:     PRStateOpen,
		Author:    "testuser",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		MergedAt:  nil,
		URL:       "https://github.com/owner/repo/pull/2",
	}

	if pr.MergedAt != nil {
		t.Error("MergedAt should be nil for unmerged PR")
	}
}

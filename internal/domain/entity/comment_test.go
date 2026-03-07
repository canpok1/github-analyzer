package entity

import (
	"testing"
	"time"
)

// === テストリスト ===
// DONE: CommentType の値が正しい（issue_comment, review_comment）
// DONE: Comment構造体に期待フィールドが設定できる
// DONE: Comment の UpdatedAt が nil の場合
// DONE: Comment の Path が空の場合（一般コメント）

func TestCommentType_Values(t *testing.T) {
	tests := []struct {
		ct   CommentType
		want string
	}{
		{CommentTypeIssue, "issue_comment"},
		{CommentTypeReview, "review_comment"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if string(tt.ct) != tt.want {
				t.Errorf("CommentType = %q, want %q", tt.ct, tt.want)
			}
		})
	}
}

func TestComment_HasExpectedFields(t *testing.T) {
	now := time.Now()
	updatedAt := now.Add(1 * time.Hour)
	comment := Comment{
		ID:        123,
		Body:      "テストコメント",
		Author:    "testuser",
		CreatedAt: now,
		UpdatedAt: &updatedAt,
		Type:      CommentTypeReview,
		Path:      "main.go",
		URL:       "https://github.com/owner/repo/pull/1#discussion_r123",
	}

	if comment.ID != 123 {
		t.Errorf("ID = %d, want 123", comment.ID)
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
	if comment.Type != CommentTypeReview {
		t.Errorf("Type = %q, want %q", comment.Type, CommentTypeReview)
	}
	if comment.Path != "main.go" {
		t.Errorf("Path = %q, want %q", comment.Path, "main.go")
	}
	if comment.URL != "https://github.com/owner/repo/pull/1#discussion_r123" {
		t.Errorf("URL = %q, want %q", comment.URL, "https://github.com/owner/repo/pull/1#discussion_r123")
	}
}

func TestComment_UpdatedAtNil(t *testing.T) {
	comment := Comment{
		ID:        1,
		Body:      "コメント",
		Author:    "user",
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		Type:      CommentTypeIssue,
	}

	if comment.UpdatedAt != nil {
		t.Error("UpdatedAt should be nil")
	}
}

func TestComment_EmptyPath(t *testing.T) {
	comment := Comment{
		ID:     1,
		Type:   CommentTypeIssue,
		Path:   "",
		Author: "user",
	}

	if comment.Path != "" {
		t.Errorf("Path = %q, want empty string", comment.Path)
	}
}

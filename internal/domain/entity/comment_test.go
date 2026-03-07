package entity

import (
	"testing"
)

// === テストリスト ===
// TODO: CommentType の値が正しい（issue_comment, review_comment）
// TODO: Comment構造体に期待フィールドが設定できる
// TODO: Comment の UpdatedAt が nil の場合
// TODO: Comment の Path が空の場合（一般コメント）

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

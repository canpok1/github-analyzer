package entity

import "time"

// CommentType はコメントの種類を表す。
type CommentType string

const (
	// CommentTypeIssue はIssue/PRの一般コメント。
	CommentTypeIssue CommentType = "issue_comment"
	// CommentTypeReview はPRのレビューコメント。
	CommentTypeReview CommentType = "review_comment"
)

// Comment はGitHubのコメントを表すエンティティ。
// Issueコメントとレビューコメントを統一的に扱う。
type Comment struct {
	ID        int64
	Body      string
	Author    string
	CreatedAt time.Time
	UpdatedAt *time.Time
	Type      CommentType
	Path      string
	URL       string
}

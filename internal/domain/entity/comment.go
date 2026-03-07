package entity

// CommentType はコメントの種類を表す。
type CommentType string

const (
	// CommentTypeIssue はIssue/PRの一般コメント。
	CommentTypeIssue CommentType = "issue_comment"
	// CommentTypeReview はPRのレビューコメント。
	CommentTypeReview CommentType = "review_comment"
)

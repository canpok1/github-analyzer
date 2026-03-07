package entity

import "time"

// IssueState はIssueの状態を表す。
type IssueState string

const (
	// IssueStateOpen はオープン状態のIssue。
	IssueStateOpen IssueState = "open"
	// IssueStateClosed はクローズ状態のIssue。
	IssueStateClosed IssueState = "closed"
)

// Issue はGitHubのIssueを表すエンティティ。
type Issue struct {
	Number    int
	Title     string
	State     IssueState
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  *time.Time
	URL       string
	Labels    []string
}

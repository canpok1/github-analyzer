package entity

import "time"

// PRState はPRの状態を表す。
type PRState string

const (
	// PRStateOpen はオープン状態のPR。
	PRStateOpen PRState = "open"
	// PRStateClosed はクローズ状態のPR。
	PRStateClosed PRState = "closed"
	// PRStateMerged はマージ済みのPR。
	PRStateMerged PRState = "merged"
)

// PullRequest はGitHubのPRを表すエンティティ。
type PullRequest struct {
	Number    int
	Title     string
	State     PRState
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	MergedAt  *time.Time
	URL       string
}

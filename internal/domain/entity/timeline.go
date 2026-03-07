package entity

import "time"

// TimelineEvent はGitHubのタイムラインイベントを表すエンティティ。
// コミット、ラベル変更、アサイン変更等のイベントを扱う。
type TimelineEvent struct {
	ID        int64
	Event     string
	Actor     string
	CreatedAt time.Time
	Label     string
	Assignee  string
	CommitID  string
	URL       string
}

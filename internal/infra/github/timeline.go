package github

import (
	"context"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// ListTimelineEvents は指定Issue/PRのタイムラインイベント一覧を取得する。
func (c *Client) ListTimelineEvents(ctx context.Context, owner, repo string, number int) ([]entity.TimelineEvent, error) {
	return nil, nil
}

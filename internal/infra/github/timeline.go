package github

import (
	"context"
	"fmt"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
)

// ListTimelineEvents は指定Issue/PRのタイムラインイベント一覧を取得する。
func (c *Client) ListTimelineEvents(ctx context.Context, owner, repo string, number int) ([]entity.TimelineEvent, error) {
	opts := &gh.ListOptions{
		PerPage: 100,
	}

	var allEvents []entity.TimelineEvent

	for {
		events, resp, err := c.client.Issues.ListIssueTimeline(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list timeline events: %w", err)
		}

		for _, event := range events {
			allEvents = append(allEvents, convertTimelineEvent(event))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allEvents, nil
}

// convertTimelineEvent はgo-githubのTimelineをドメインエンティティに変換する。
func convertTimelineEvent(event *gh.Timeline) entity.TimelineEvent {
	result := entity.TimelineEvent{
		ID:       event.GetID(),
		Event:    event.GetEvent(),
		Actor:    event.GetActor().GetLogin(),
		CommitID: event.GetCommitID(),
		URL:      event.GetURL(),
	}

	if event.CreatedAt != nil {
		result.CreatedAt = event.CreatedAt.Time
	}

	if event.Label != nil {
		result.Label = event.Label.GetName()
	}

	if event.Assignee != nil {
		result.Assignee = event.Assignee.GetLogin()
	}

	return result
}

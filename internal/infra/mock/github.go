package mock

import (
	"context"
	"time"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

var baseTime = time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)

// GitHubRepository は domain.GitHubRepository のモック実装。
// 固定のダミーデータを返す。
type GitHubRepository struct{}

func (r *GitHubRepository) ListPullRequests(_ context.Context, _, _ string, _ domain.ListPullRequestsOptions) ([]entity.PullRequest, error) {
	mergedAt := baseTime.Add(48 * time.Hour)
	return []entity.PullRequest{
		{
			Number:    1,
			Title:     "[Mock] 機能追加のPR",
			State:     entity.PRStateMerged,
			Author:    "mock-user",
			CreatedAt: baseTime,
			UpdatedAt: baseTime.Add(24 * time.Hour),
			MergedAt:  &mergedAt,
			URL:       "https://github.com/example/repo/pull/1",
		},
		{
			Number:    2,
			Title:     "[Mock] バグ修正のPR",
			State:     entity.PRStateOpen,
			Author:    "mock-reviewer",
			CreatedAt: baseTime.Add(72 * time.Hour),
			UpdatedAt: baseTime.Add(96 * time.Hour),
			URL:       "https://github.com/example/repo/pull/2",
		},
	}, nil
}

func (r *GitHubRepository) ListIssues(_ context.Context, _, _ string, _ domain.ListIssuesOptions) ([]entity.Issue, error) {
	closedAt := baseTime.Add(72 * time.Hour)
	return []entity.Issue{
		{
			Number:    10,
			Title:     "[Mock] 機能要望",
			State:     entity.IssueStateClosed,
			Author:    "mock-user",
			CreatedAt: baseTime,
			UpdatedAt: baseTime.Add(48 * time.Hour),
			ClosedAt:  &closedAt,
			URL:       "https://github.com/example/repo/issues/10",
			Labels:    []string{"enhancement"},
		},
		{
			Number:    11,
			Title:     "[Mock] バグ報告",
			State:     entity.IssueStateOpen,
			Author:    "mock-reporter",
			CreatedAt: baseTime.Add(24 * time.Hour),
			UpdatedAt: baseTime.Add(48 * time.Hour),
			URL:       "https://github.com/example/repo/issues/11",
			Labels:    []string{"bug"},
		},
	}, nil
}

func (r *GitHubRepository) ListIssueComments(_ context.Context, _, _ string, _ int) ([]entity.Comment, error) {
	return []entity.Comment{
		{
			ID:        1,
			Body:      "[Mock] 対応方針について確認しました。",
			Author:    "mock-reviewer",
			CreatedAt: baseTime.Add(2 * time.Hour),
			Type:      entity.CommentTypeIssue,
		},
	}, nil
}

func (r *GitHubRepository) ListPullRequestComments(_ context.Context, _, _ string, _ int) ([]entity.Comment, error) {
	return []entity.Comment{
		{
			ID:        2,
			Body:      "[Mock] コードレビューコメントです。",
			Author:    "mock-reviewer",
			CreatedAt: baseTime.Add(26 * time.Hour),
			Type:      entity.CommentTypeReview,
			Path:      "main.go",
		},
	}, nil
}

func (r *GitHubRepository) ListTimelineEvents(_ context.Context, _, _ string, _ int) ([]entity.TimelineEvent, error) {
	return []entity.TimelineEvent{
		{
			ID:        1,
			Event:     "labeled",
			Actor:     "mock-user",
			CreatedAt: baseTime.Add(1 * time.Hour),
			Label:     "enhancement",
		},
		{
			ID:        2,
			Event:     "assigned",
			Actor:     "mock-user",
			CreatedAt: baseTime.Add(1 * time.Hour),
			Assignee:  "mock-developer",
		},
	}, nil
}

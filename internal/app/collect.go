package app

import (
	"context"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// CollectedData は収集したデータを集約する構造体。
type CollectedData struct {
	PullRequests []entity.PullRequest
	Issues       []entity.Issue
	Comments     map[int][]entity.Comment
	Timeline     map[int][]entity.TimelineEvent
}

// CollectData はQueryに応じてGitHub APIからデータを収集する。
func CollectData(_ context.Context, _ domain.GitHubRepository, _ entity.Query) (*CollectedData, error) {
	return &CollectedData{
		Comments: make(map[int][]entity.Comment),
		Timeline: make(map[int][]entity.TimelineEvent),
	}, nil
}

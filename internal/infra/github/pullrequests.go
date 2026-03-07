package github

import (
	"context"
	"fmt"

	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
)

// ListPullRequests は指定リポジトリのPR一覧を取得する。
func (c *Client) ListPullRequests(ctx context.Context, owner, repo string, opts domain.ListPullRequestsOptions) ([]entity.PullRequest, error) {
	apiState := resolveAPIState(opts.Status)

	ghOpts := &gh.PullRequestListOptions{
		State:     apiState,
		Sort:      "updated",
		Direction: "desc",
		ListOptions: gh.ListOptions{
			PerPage: 100,
		},
	}

	var allPRs []entity.PullRequest

	for {
		prs, resp, err := c.client.PullRequests.List(ctx, owner, repo, ghOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list pull requests: %w", err)
		}

		for _, pr := range prs {
			// Sinceフィルタ: updatedで降順ソート済みのため、Since以前なら以降も全て古い
			if opts.Since != nil && pr.GetUpdatedAt().Time.Before(*opts.Since) {
				return allPRs, nil
			}

			// mergedフィルタ: MergedAtがnilのPR(クローズされたがマージされていない)は除外
			if opts.Status == entity.PRStateMerged && pr.MergedAt == nil {
				continue
			}

			allPRs = append(allPRs, convertPullRequest(pr))
		}

		if resp.NextPage == 0 {
			break
		}
		ghOpts.Page = resp.NextPage
	}

	return allPRs, nil
}

// resolveAPIState はStatusフィルタをGitHub API用のstate値に変換する。
func resolveAPIState(status entity.PRState) string {
	switch status {
	case entity.PRStateOpen:
		return "open"
	case entity.PRStateClosed:
		return "closed"
	case entity.PRStateMerged:
		// GitHub APIにはmerged stateがないため、closedで取得してフィルタする
		return "closed"
	default:
		return "all"
	}
}

// convertPullRequest はgo-githubのPullRequestをドメインエンティティに変換する。
func convertPullRequest(pr *gh.PullRequest) entity.PullRequest {
	result := entity.PullRequest{
		Number: pr.GetNumber(),
		Title:  pr.GetTitle(),
		Author: pr.GetUser().GetLogin(),
		URL:    pr.GetHTMLURL(),
	}

	if pr.CreatedAt != nil {
		result.CreatedAt = pr.CreatedAt.Time
	}
	if pr.UpdatedAt != nil {
		result.UpdatedAt = pr.UpdatedAt.Time
	}
	if pr.MergedAt != nil {
		t := pr.MergedAt.Time
		result.MergedAt = &t
		result.State = entity.PRStateMerged
	} else if pr.GetState() == "closed" {
		result.State = entity.PRStateClosed
	} else {
		result.State = entity.PRStateOpen
	}

	return result
}

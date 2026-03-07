package github

import (
	"context"
	"fmt"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
	gh "github.com/google/go-github/v68/github"
)

// ListIssueComments は指定Issue/PRの一般コメント一覧を取得する。
func (c *Client) ListIssueComments(ctx context.Context, owner, repo string, number int) ([]entity.Comment, error) {
	opts := &gh.IssueListCommentsOptions{
		ListOptions: gh.ListOptions{
			PerPage: 100,
		},
	}

	var allComments []entity.Comment

	for {
		comments, resp, err := c.client.Issues.ListComments(ctx, owner, repo, number, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issue comments: %w", err)
		}

		for _, comment := range comments {
			allComments = append(allComments, convertIssueComment(comment))
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allComments, nil
}

// ListPullRequestComments は指定PRのレビューコメント一覧を取得する。
func (c *Client) ListPullRequestComments(ctx context.Context, owner, repo string, number int) ([]entity.Comment, error) {
	return nil, nil
}

// convertIssueComment はgo-githubのIssueCommentをドメインエンティティに変換する。
func convertIssueComment(comment *gh.IssueComment) entity.Comment {
	result := entity.Comment{
		ID:     comment.GetID(),
		Body:   comment.GetBody(),
		Author: comment.GetUser().GetLogin(),
		Type:   entity.CommentTypeIssue,
		URL:    comment.GetHTMLURL(),
	}

	if comment.CreatedAt != nil {
		result.CreatedAt = comment.CreatedAt.Time
	}
	if comment.UpdatedAt != nil {
		t := comment.UpdatedAt.Time
		result.UpdatedAt = &t
	}

	return result
}

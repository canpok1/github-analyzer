package gemini

import (
	"fmt"
	"strings"

	"github.com/canpok1/github-analyzer/internal/app"
	"github.com/canpok1/github-analyzer/internal/domain"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// DefaultPrompt はユーザーがプロンプトを指定しなかった場合に使用するデフォルトプロンプト。
const DefaultPrompt = "チームの開発プロセスの健全性を分析し、改善点を提案してください。"

// systemPromptTemplate はレポート構成を指示するシステムプロンプト。
const systemPromptTemplate = `あなたはソフトウェア開発チームのプロセスアナリストです。
以下のGitHubアクティビティデータを分析し、次の構成でレポートを作成してください。

## Overview
活動の「空気感」を1行で要約してください。

## Process Insights
指定されたプロンプトに対する分析結果を記述してください。

## Potential Risks
停滞や対立など、注意すべきリスクを指摘してください。

## Manager's Hint
マネージャーへの具体的なアクション提案を記述してください。
`

// BuildPrompt はCollectedDataとユーザープロンプトからAnalysisRequestを構築する。
func BuildPrompt(data *app.CollectedData, userPrompt string) domain.AnalysisRequest {
	prompt := userPrompt
	if prompt == "" {
		prompt = DefaultPrompt
	}

	fullPrompt := systemPromptTemplate + "\nユーザーからの分析指示: " + prompt

	dataStr := buildDataString(data)

	return domain.AnalysisRequest{
		Prompt: fullPrompt,
		Data:   dataStr,
	}
}

// buildDataString はCollectedDataを構造化テキストに変換する。
func buildDataString(data *app.CollectedData) string {
	var sb strings.Builder

	if len(data.PullRequests) > 0 {
		sb.WriteString("# Pull Requests\n\n")
		for _, pr := range data.PullRequests {
			writePR(&sb, &pr)
		}
	}

	if len(data.Issues) > 0 {
		sb.WriteString("# Issues\n\n")
		for _, issue := range data.Issues {
			writeIssue(&sb, &issue)
		}
	}

	if len(data.Comments) > 0 {
		sb.WriteString("# Comments\n\n")
		for number, comments := range data.Comments {
			fmt.Fprintf(&sb, "## Comments for #%d\n\n", number)
			for _, comment := range comments {
				writeComment(&sb, &comment)
			}
		}
	}

	if len(data.Timeline) > 0 {
		sb.WriteString("# Timeline Events\n\n")
		for number, events := range data.Timeline {
			fmt.Fprintf(&sb, "## Timeline for #%d\n\n", number)
			for _, event := range events {
				writeTimelineEvent(&sb, &event)
			}
		}
	}

	return sb.String()
}

// writePR はPR情報をフォーマットして書き込む。
func writePR(sb *strings.Builder, pr *entity.PullRequest) {
	fmt.Fprintf(sb, "## #%d: %s\n", pr.Number, pr.Title)
	fmt.Fprintf(sb, "- State: %s\n", pr.State)
	fmt.Fprintf(sb, "- Author: %s\n", pr.Author)
	fmt.Fprintf(sb, "- Created: %s\n", pr.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(sb, "- Updated: %s\n", pr.UpdatedAt.Format("2006-01-02 15:04:05"))
	if pr.MergedAt != nil {
		fmt.Fprintf(sb, "- Merged: %s\n", pr.MergedAt.Format("2006-01-02 15:04:05"))
	}
	sb.WriteString("\n")
}

// writeIssue はIssue情報をフォーマットして書き込む。
func writeIssue(sb *strings.Builder, issue *entity.Issue) {
	fmt.Fprintf(sb, "## #%d: %s\n", issue.Number, issue.Title)
	fmt.Fprintf(sb, "- State: %s\n", issue.State)
	fmt.Fprintf(sb, "- Author: %s\n", issue.Author)
	fmt.Fprintf(sb, "- Created: %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(sb, "- Updated: %s\n", issue.UpdatedAt.Format("2006-01-02 15:04:05"))
	if issue.ClosedAt != nil {
		fmt.Fprintf(sb, "- Closed: %s\n", issue.ClosedAt.Format("2006-01-02 15:04:05"))
	}
	if len(issue.Labels) > 0 {
		fmt.Fprintf(sb, "- Labels: %s\n", strings.Join(issue.Labels, ", "))
	}
	sb.WriteString("\n")
}

// writeComment はコメント情報をフォーマットして書き込む。
func writeComment(sb *strings.Builder, comment *entity.Comment) {
	fmt.Fprintf(sb, "- [%s] %s (%s): %s\n",
		comment.CreatedAt.Format("2006-01-02 15:04:05"),
		comment.Author,
		comment.Type,
		comment.Body,
	)
}

// writeTimelineEvent はタイムラインイベント情報をフォーマットして書き込む。
func writeTimelineEvent(sb *strings.Builder, event *entity.TimelineEvent) {
	fmt.Fprintf(sb, "- [%s] %s by %s",
		event.CreatedAt.Format("2006-01-02 15:04:05"),
		event.Event,
		event.Actor,
	)
	if event.Label != "" {
		fmt.Fprintf(sb, " (label: %s)", event.Label)
	}
	if event.Assignee != "" {
		fmt.Fprintf(sb, " (assignee: %s)", event.Assignee)
	}
	if event.CommitID != "" {
		fmt.Fprintf(sb, " (commit: %s)", event.CommitID)
	}
	sb.WriteString("\n")
}

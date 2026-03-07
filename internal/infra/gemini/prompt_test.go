package gemini

import (
	"strings"
	"testing"
	"time"

	"github.com/canpok1/github-analyzer/internal/app"
	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// DONE: BuildPrompt: userPromptが空の場合、デフォルトプロンプトが使用される
// DONE: BuildPrompt: userPromptが指定された場合、そのプロンプトが使用される
// DONE: BuildPrompt: システムプロンプトにレポート構成指示が含まれる
// DONE: BuildPrompt: PRデータが構造化されてDataに含まれる
// DONE: BuildPrompt: Issueデータが構造化されてDataに含まれる
// DONE: BuildPrompt: コメントデータが構造化されてDataに含まれる
// DONE: BuildPrompt: タイムラインデータが構造化されてDataに含まれる
// DONE: BuildPrompt: データが空の場合でもエラーにならない
// DONE: BuildPrompt: 複数のPR/Issueが含まれる場合の出力

func TestBuildPrompt_EmptyData(t *testing.T) {
	data := &app.CollectedData{}

	result := BuildPrompt(data, "")

	if result.Prompt == "" {
		t.Error("Prompt should not be empty")
	}
	if !strings.Contains(result.Prompt, "Overview") {
		t.Error("Prompt should contain report structure with Overview")
	}
}

func TestBuildPrompt_DefaultPrompt(t *testing.T) {
	data := &app.CollectedData{}

	result := BuildPrompt(data, "")

	if !strings.Contains(result.Prompt, DefaultPrompt) {
		t.Errorf("Prompt should contain default prompt %q, got %q", DefaultPrompt, result.Prompt)
	}
}

func TestBuildPrompt_UserPrompt(t *testing.T) {
	data := &app.CollectedData{}
	userPrompt := "レビュー速度を分析してください"

	result := BuildPrompt(data, userPrompt)

	if !strings.Contains(result.Prompt, userPrompt) {
		t.Errorf("Prompt should contain user prompt %q, got %q", userPrompt, result.Prompt)
	}
	if strings.Contains(result.Prompt, DefaultPrompt) {
		t.Error("Prompt should not contain default prompt when user prompt is specified")
	}
}

func TestBuildPrompt_ReportStructure(t *testing.T) {
	data := &app.CollectedData{}

	result := BuildPrompt(data, "")

	sections := []string{"Overview", "Process Insights", "Potential Risks", "Manager's Hint"}
	for _, section := range sections {
		if !strings.Contains(result.Prompt, section) {
			t.Errorf("Prompt should contain section %q", section)
		}
	}
}

func TestBuildPrompt_PRData(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	data := &app.CollectedData{
		PullRequests: []entity.PullRequest{
			{
				Number:    42,
				Title:     "Add feature X",
				State:     entity.PRStateOpen,
				Author:    "alice",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}

	result := BuildPrompt(data, "")

	if !strings.Contains(result.Data, "#42") {
		t.Error("Data should contain PR number")
	}
	if !strings.Contains(result.Data, "Add feature X") {
		t.Error("Data should contain PR title")
	}
	if !strings.Contains(result.Data, "alice") {
		t.Error("Data should contain PR author")
	}
	if !strings.Contains(result.Data, "open") {
		t.Error("Data should contain PR state")
	}
}

func TestBuildPrompt_IssueData(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	data := &app.CollectedData{
		Issues: []entity.Issue{
			{
				Number:    10,
				Title:     "Bug report",
				State:     entity.IssueStateOpen,
				Author:    "bob",
				CreatedAt: now,
				UpdatedAt: now,
				Labels:    []string{"bug", "priority:high"},
			},
		},
	}

	result := BuildPrompt(data, "")

	if !strings.Contains(result.Data, "#10") {
		t.Error("Data should contain Issue number")
	}
	if !strings.Contains(result.Data, "Bug report") {
		t.Error("Data should contain Issue title")
	}
	if !strings.Contains(result.Data, "bob") {
		t.Error("Data should contain Issue author")
	}
	if !strings.Contains(result.Data, "bug") {
		t.Error("Data should contain Issue labels")
	}
}

func TestBuildPrompt_CommentData(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	data := &app.CollectedData{
		Comments: map[int][]entity.Comment{
			42: {
				{
					ID:        1,
					Body:      "LGTM! Great work.",
					Author:    "carol",
					CreatedAt: now,
					Type:      entity.CommentTypeIssue,
				},
			},
		},
	}

	result := BuildPrompt(data, "")

	if !strings.Contains(result.Data, "#42") {
		t.Error("Data should contain comment's associated number")
	}
	if !strings.Contains(result.Data, "LGTM! Great work.") {
		t.Error("Data should contain comment body")
	}
	if !strings.Contains(result.Data, "carol") {
		t.Error("Data should contain comment author")
	}
}

func TestBuildPrompt_TimelineData(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	data := &app.CollectedData{
		Timeline: map[int][]entity.TimelineEvent{
			42: {
				{
					ID:        1,
					Event:     "labeled",
					Actor:     "dave",
					CreatedAt: now,
					Label:     "bug",
				},
			},
		},
	}

	result := BuildPrompt(data, "")

	if !strings.Contains(result.Data, "#42") {
		t.Error("Data should contain timeline's associated number")
	}
	if !strings.Contains(result.Data, "labeled") {
		t.Error("Data should contain timeline event type")
	}
	if !strings.Contains(result.Data, "dave") {
		t.Error("Data should contain timeline actor")
	}
}

func TestBuildPrompt_MultiplePRsAndIssues(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC)
	mergedAt := time.Date(2025, 1, 16, 10, 0, 0, 0, time.UTC)
	closedAt := time.Date(2025, 1, 17, 10, 0, 0, 0, time.UTC)
	data := &app.CollectedData{
		PullRequests: []entity.PullRequest{
			{
				Number:    1,
				Title:     "First PR",
				State:     entity.PRStateMerged,
				Author:    "alice",
				CreatedAt: now,
				UpdatedAt: now,
				MergedAt:  &mergedAt,
			},
			{
				Number:    2,
				Title:     "Second PR",
				State:     entity.PRStateOpen,
				Author:    "bob",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Issues: []entity.Issue{
			{
				Number:    3,
				Title:     "First Issue",
				State:     entity.IssueStateClosed,
				Author:    "carol",
				CreatedAt: now,
				UpdatedAt: now,
				ClosedAt:  &closedAt,
			},
			{
				Number:    4,
				Title:     "Second Issue",
				State:     entity.IssueStateOpen,
				Author:    "dave",
				CreatedAt: now,
				UpdatedAt: now,
				Labels:    []string{"enhancement"},
			},
		},
		Comments: map[int][]entity.Comment{
			1: {
				{ID: 1, Body: "Comment on PR 1", Author: "eve", CreatedAt: now, Type: entity.CommentTypeReview},
			},
			3: {
				{ID: 2, Body: "Comment on Issue 3", Author: "frank", CreatedAt: now, Type: entity.CommentTypeIssue},
			},
		},
		Timeline: map[int][]entity.TimelineEvent{
			1: {
				{ID: 1, Event: "merged", Actor: "alice", CreatedAt: mergedAt},
			},
		},
	}

	result := BuildPrompt(data, "カスタム分析")

	// PRs
	if !strings.Contains(result.Data, "#1") || !strings.Contains(result.Data, "First PR") {
		t.Error("Data should contain first PR")
	}
	if !strings.Contains(result.Data, "#2") || !strings.Contains(result.Data, "Second PR") {
		t.Error("Data should contain second PR")
	}
	if !strings.Contains(result.Data, "Merged:") {
		t.Error("Data should contain merged date for merged PR")
	}

	// Issues
	if !strings.Contains(result.Data, "#3") || !strings.Contains(result.Data, "First Issue") {
		t.Error("Data should contain first issue")
	}
	if !strings.Contains(result.Data, "#4") || !strings.Contains(result.Data, "Second Issue") {
		t.Error("Data should contain second issue")
	}
	if !strings.Contains(result.Data, "Closed:") {
		t.Error("Data should contain closed date for closed issue")
	}
	if !strings.Contains(result.Data, "enhancement") {
		t.Error("Data should contain labels")
	}

	// Comments
	if !strings.Contains(result.Data, "Comment on PR 1") {
		t.Error("Data should contain PR comment")
	}
	if !strings.Contains(result.Data, "Comment on Issue 3") {
		t.Error("Data should contain issue comment")
	}
	if !strings.Contains(result.Data, "review_comment") {
		t.Error("Data should contain comment type")
	}

	// Timeline
	if !strings.Contains(result.Data, "merged") {
		t.Error("Data should contain timeline event")
	}

	// Prompt
	if !strings.Contains(result.Prompt, "カスタム分析") {
		t.Error("Prompt should contain user prompt")
	}
	if strings.Contains(result.Prompt, DefaultPrompt) {
		t.Error("Prompt should not contain default prompt")
	}

	// Sections
	if !strings.Contains(result.Data, "# Pull Requests") {
		t.Error("Data should have Pull Requests section header")
	}
	if !strings.Contains(result.Data, "# Issues") {
		t.Error("Data should have Issues section header")
	}
	if !strings.Contains(result.Data, "# Comments") {
		t.Error("Data should have Comments section header")
	}
	if !strings.Contains(result.Data, "# Timeline Events") {
		t.Error("Data should have Timeline Events section header")
	}
}

package mock

import (
	"context"
	"strings"
	"testing"

	"github.com/canpok1/github-analyzer/internal/domain"
)

func TestAnalyzer_Analyze(t *testing.T) {
	analyzer := &Analyzer{}
	resp, err := analyzer.Analyze(context.Background(), domain.AnalysisRequest{
		Prompt: "test prompt",
		Data:   "test data",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response should not be nil")
	}
	if resp.Content == "" {
		t.Error("response content should not be empty")
	}
}

func TestAnalyzer_ReportContainsAllSections(t *testing.T) {
	analyzer := &Analyzer{}
	resp, err := analyzer.Analyze(context.Background(), domain.AnalysisRequest{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sections := []string{
		"## Overview",
		"## Process Insights",
		"## Potential Risks",
		"## Manager's Hint",
	}
	for _, section := range sections {
		if !strings.Contains(resp.Content, section) {
			t.Errorf("response should contain %q", section)
		}
	}
}

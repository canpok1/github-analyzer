package gemini

import (
	"fmt"
	"strings"

	"github.com/canpok1/github-analyzer/internal/domain/entity"
)

// sectionHeaders はレポートの各セクション見出し。
var sectionHeaders = []string{
	"## Overview",
	"## Process Insights",
	"## Potential Risks",
	"## Manager's Hint",
}

// ParseReport はGemini APIの生テキスト出力をReportエンティティにパースする。
func ParseReport(content string) (*entity.Report, error) {
	if content == "" {
		return nil, fmt.Errorf("empty content")
	}

	sections := make(map[string]string, len(sectionHeaders))

	for i, header := range sectionHeaders {
		idx := strings.Index(content, header)
		if idx == -1 {
			return nil, fmt.Errorf("missing section: %s", header)
		}

		start := idx + len(header)
		var end int
		if i+1 < len(sectionHeaders) {
			nextIdx := strings.Index(content, sectionHeaders[i+1])
			if nextIdx == -1 {
				end = len(content)
			} else {
				end = nextIdx
			}
		} else {
			end = len(content)
		}

		body := strings.TrimSpace(content[start:end])
		sections[header] = body
	}

	return &entity.Report{
		Overview:        sections["## Overview"],
		ProcessInsights: sections["## Process Insights"],
		PotentialRisks:  sections["## Potential Risks"],
		ManagersHint:    sections["## Manager's Hint"],
	}, nil
}

package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Query は分析対象の検索条件を表す。
type Query struct {
	Since  *time.Time
	PR     *int
	Issue  *int
	Status string
	Prompt string
	Repo   string
	Model  string
}

// ParseDuration は "7d", "2w" のような期間文字列をtime.Durationに変換する。
func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("duration string is empty")
	}

	s = strings.TrimSpace(s)
	unit := s[len(s)-1:]
	numStr := s[:len(s)-1]

	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q: %w", s, err)
	}

	switch unit {
	case "d":
		return time.Duration(num) * 24 * time.Hour, nil
	case "w":
		return time.Duration(num) * 7 * 24 * time.Hour, nil
	case "m":
		return time.Duration(num) * 30 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown duration unit %q in %q", unit, s)
	}
}

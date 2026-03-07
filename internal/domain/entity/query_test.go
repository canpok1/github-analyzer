package entity

import (
	"testing"
	"time"
)

func TestQuery_HasExpectedFields(t *testing.T) {
	now := time.Now()
	q := Query{
		Since:  &now,
		PR:     intPtr(123),
		Issue:  intPtr(456),
		Status: "open",
		Prompt: "テスト分析",
		Repo:   "owner/repo",
	}

	if q.Since == nil || !q.Since.Equal(now) {
		t.Error("Since field not set correctly")
	}
	if q.PR == nil || *q.PR != 123 {
		t.Error("PR field not set correctly")
	}
	if q.Issue == nil || *q.Issue != 456 {
		t.Error("Issue field not set correctly")
	}
	if q.Status != "open" {
		t.Errorf("Status = %q, want %q", q.Status, "open")
	}
	if q.Prompt != "テスト分析" {
		t.Errorf("Prompt = %q, want %q", q.Prompt, "テスト分析")
	}
	if q.Repo != "owner/repo" {
		t.Errorf("Repo = %q, want %q", q.Repo, "owner/repo")
	}
}

func TestParseDuration_7d(t *testing.T) {
	d, err := ParseDuration("7d")
	if err != nil {
		t.Fatalf("ParseDuration(\"7d\") returned error: %v", err)
	}
	expected := 7 * 24 * time.Hour
	if d != expected {
		t.Errorf("ParseDuration(\"7d\") = %v, want %v", d, expected)
	}
}

func TestParseDuration_2w(t *testing.T) {
	d, err := ParseDuration("2w")
	if err != nil {
		t.Fatalf("ParseDuration(\"2w\") returned error: %v", err)
	}
	expected := 2 * 7 * 24 * time.Hour
	if d != expected {
		t.Errorf("ParseDuration(\"2w\") = %v, want %v", d, expected)
	}
}

func TestParseDuration_1m(t *testing.T) {
	d, err := ParseDuration("1m")
	if err != nil {
		t.Fatalf("ParseDuration(\"1m\") returned error: %v", err)
	}
	expected := 30 * 24 * time.Hour
	if d != expected {
		t.Errorf("ParseDuration(\"1m\") = %v, want %v", d, expected)
	}
}

func TestParseDuration_InvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"no number", "d"},
		{"unknown unit", "7x"},
		{"completely invalid", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseDuration(tt.input)
			if err == nil {
				t.Errorf("ParseDuration(%q) expected error, got nil", tt.input)
			}
		})
	}
}

func intPtr(n int) *int {
	return &n
}

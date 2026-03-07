package github

import (
	"testing"
)

// TODO: 正常系: HTTPS形式のリモートURLからowner/repoを抽出
// TODO: 正常系: SSH形式のリモートURLからowner/repoを抽出
// TODO: 正常系: .git拡張子付きURLからowner/repoを抽出
// TODO: 異常系: 不正なURL形式でエラー
// TODO: 正常系: DetectRepoが実際のgitリポジトリで動作する

func TestParseRemoteURL_HTTPS(t *testing.T) {
	owner, repo, err := ParseRemoteURL("https://github.com/owner/repo-name.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "owner" {
		t.Errorf("owner = %q, want %q", owner, "owner")
	}
	if repo != "repo-name" {
		t.Errorf("repo = %q, want %q", repo, "repo-name")
	}
}

func TestParseRemoteURL_HTTPSWithoutGit(t *testing.T) {
	owner, repo, err := ParseRemoteURL("https://github.com/owner/repo-name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "owner" {
		t.Errorf("owner = %q, want %q", owner, "owner")
	}
	if repo != "repo-name" {
		t.Errorf("repo = %q, want %q", repo, "repo-name")
	}
}

func TestParseRemoteURL_SSH(t *testing.T) {
	owner, repo, err := ParseRemoteURL("git@github.com:owner/repo-name.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "owner" {
		t.Errorf("owner = %q, want %q", owner, "owner")
	}
	if repo != "repo-name" {
		t.Errorf("repo = %q, want %q", repo, "repo-name")
	}
}

func TestParseRemoteURL_SSHWithoutGit(t *testing.T) {
	owner, repo, err := ParseRemoteURL("git@github.com:owner/repo-name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if owner != "owner" {
		t.Errorf("owner = %q, want %q", owner, "owner")
	}
	if repo != "repo-name" {
		t.Errorf("repo = %q, want %q", repo, "repo-name")
	}
}

func TestParseRemoteURL_Invalid(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"empty", ""},
		{"no path", "https://github.com"},
		{"single path", "https://github.com/owner"},
		{"random string", "not-a-url"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParseRemoteURL(tt.url)
			if err == nil {
				t.Errorf("ParseRemoteURL(%q) expected error, got nil", tt.url)
			}
		})
	}
}

func TestDetectRepo_CurrentDirectory(t *testing.T) {
	// 実際のgitリポジトリで動作確認
	owner, repo, err := DetectRepo()
	if err != nil {
		t.Fatalf("DetectRepo() returned error: %v", err)
	}
	if owner == "" {
		t.Error("owner should not be empty")
	}
	if repo == "" {
		t.Error("repo should not be empty")
	}
}

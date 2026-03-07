package github

import (
	"fmt"
	"os/exec"
	"strings"
)

// ParseRemoteURL はgit remoteのURLからowner, repoを抽出する。
// HTTPS形式: https://github.com/owner/repo.git
// SSH形式: git@github.com:owner/repo.git
func ParseRemoteURL(rawURL string) (string, string, error) {
	if rawURL == "" {
		return "", "", fmt.Errorf("remote URL is empty")
	}

	var path string

	if strings.HasPrefix(rawURL, "git@") {
		// SSH形式: git@github.com:owner/repo.git
		colonIdx := strings.Index(rawURL, ":")
		if colonIdx < 0 {
			return "", "", fmt.Errorf("invalid SSH remote URL: %q", rawURL)
		}
		path = rawURL[colonIdx+1:]
	} else if strings.Contains(rawURL, "://") {
		// HTTPS形式: https://github.com/owner/repo.git
		parts := strings.SplitN(rawURL, "://", 2)
		if len(parts) < 2 {
			return "", "", fmt.Errorf("invalid HTTPS remote URL: %q", rawURL)
		}
		// ホスト名を除去
		hostAndPath := parts[1]
		slashIdx := strings.Index(hostAndPath, "/")
		if slashIdx < 0 {
			return "", "", fmt.Errorf("invalid HTTPS remote URL: %q", rawURL)
		}
		path = hostAndPath[slashIdx+1:]
	} else {
		return "", "", fmt.Errorf("unsupported remote URL format: %q", rawURL)
	}

	// .git拡張子を除去
	path = strings.TrimSuffix(path, ".git")

	// owner/repo に分割
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("cannot extract owner/repo from URL: %q", rawURL)
	}

	return parts[0], parts[1], nil
}

// DetectRepo はカレントディレクトリのgit remoteからowner/repoを検出する。
func DetectRepo() (string, string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get git remote URL: %w", err)
	}

	remoteURL := strings.TrimSpace(string(out))
	return ParseRemoteURL(remoteURL)
}

#!/usr/bin/env bash
set -euo pipefail

# オプション解析
USE_PRINT_MODE=false
while getopts "p" opt; do
  case "$opt" in
    p) USE_PRINT_MODE=true ;;
    *) echo "Usage: $0 [-p]" >&2; exit 1 ;;
  esac
done

echo "Issue自動選定を開始します"

# open状態のIssue数を確認し、0件ならスキップ
OPEN_COUNT=$(gh issue list --state open --json number --jq 'length')
if [ "$OPEN_COUNT" -eq 0 ]; then
  echo "open状態のIssueがないため、スキップします"
  exit 0
fi

# Claudeでissueを選定・ラベル付与（コード変更不要のため--worktreeは不使用）
if "${USE_PRINT_MODE}"; then
  claude --dangerously-skip-permissions -p "/assign-issues"
else
  claude --dangerously-skip-permissions "/assign-issues"
fi

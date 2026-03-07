#!/bin/bash
# create-version-tag.sh のテスト
# --dry-run モードを使ってテストする
#
# テストリスト:
# DONE: タグが存在しない場合、v0.0.1 を新しいバージョンとする
# TODO: v0.0.1 が存在する場合、v0.0.2 を新しいバージョンとする
# TODO: v1.2.3 が存在する場合、v1.2.4 を新しいバージョンとする
# TODO: プレリリースタグ（v1.0.0-rc1）が存在する場合、無視される
# TODO: 複数タグが存在する場合、最新のものからインクリメントする
# TODO: --dry-run モードではタグの作成・プッシュをスキップする
# TODO: 不明なオプションが渡された場合、エラー終了する

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SCRIPT_UNDER_TEST="${SCRIPT_DIR}/create-version-tag.sh"

# テスト用の一時ディレクトリ
TEST_DIR=""

# テスト結果カウンタ
PASSED=0
FAILED=0

# セットアップ: テスト用のgitリポジトリを作成
setup() {
    TEST_DIR=$(mktemp -d)
    cd "$TEST_DIR"
    git init --initial-branch=main
    git config user.email "test@example.com"
    git config user.name "Test User"
    # 最低1つのコミットが必要
    echo "initial" > README.md
    git add README.md
    git commit -m "Initial commit"
}

# クリーンアップ
teardown() {
    cd /
    rm -rf "$TEST_DIR"
}

# アサーション: 出力に文字列が含まれるか
assert_output_contains() {
    local expected="$1"
    local actual="$2"
    if echo "$actual" | grep -qF "$expected"; then
        return 0
    else
        echo "  FAIL: 出力に '$expected' が含まれていません"
        echo "  実際の出力: $actual"
        return 1
    fi
}

# アサーション: 終了コードの確認
assert_exit_code() {
    local expected="$1"
    local actual="$2"
    if [[ "$actual" -eq "$expected" ]]; then
        return 0
    else
        echo "  FAIL: 終了コードが $expected ではなく $actual です"
        return 1
    fi
}

# テスト実行ヘルパー
run_test() {
    local test_name="$1"
    local test_func="$2"

    echo "--- $test_name"
    setup
    if $test_func; then
        echo "  PASS"
        PASSED=$((PASSED + 1))
    else
        echo "  FAIL"
        FAILED=$((FAILED + 1))
    fi
    teardown
}

# ====================
# テストケース
# ====================

echo "=== create-version-tag.sh テスト ==="

# テスト1: タグが存在しない場合、v0.0.1 を新しいバージョンとする
test_no_existing_tags() {
    # タグがない状態で実行
    local output
    output=$("$SCRIPT_UNDER_TEST" --dry-run 2>&1)
    assert_output_contains "新しいバージョン: v0.0.1" "$output"
}
run_test "タグが存在しない場合、v0.0.1を新しいバージョンとする" test_no_existing_tags

# テスト結果のサマリー
print_summary() {
    echo ""
    echo "=== テスト結果 ==="
    echo "PASSED: $PASSED"
    echo "FAILED: $FAILED"
    if [[ "$FAILED" -gt 0 ]]; then
        exit 1
    fi
}

trap print_summary EXIT

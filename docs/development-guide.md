# 開発ガイド

## 開発者クイックスタート

### devcontainer（推奨）

本リポジトリにはdevcontainer環境が整備されており、開発に必要なツールが自動でインストールされます。

**起動手順:**

1. VS Codeに [Dev Containers拡張](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) をインストール
2. リポジトリをクローンしてVS Codeで開く
3. コマンドパレット（`Ctrl+Shift+P`）から「Dev Containers: Reopen in Container」を選択

**プリインストールされるツール:**

- Claude Code
- Go
- GitHub CLI（`gh`）
- tmux

**環境変数の設定:**

`.devcontainer/.env-template` をコピーして `.devcontainer/.env` を作成し、必要なキーを設定してください（初回起動時に自動コピーされます）。

```bash
GH_TOKEN=xxxx        # GitHub API トークン
GEMINI_API_KEY=xxxx   # Gemini API キー
```

Claude Codeを利用する場合は、`ANTHROPIC_API_KEY` も設定してください。

### 手動セットアップ

devcontainerを使わない場合は、以下を手動でインストールしてください。

1. Go（バージョンは `go.mod` を参照）
2. GitHub CLI（`gh`）
3. 開発ツールのインストール: `make setup`

## 前提条件

- Go（バージョンは `go.mod` を参照）

## セットアップ

開発に必要なツール（golangci-lint, go-depcheck）をインストールします。

```bash
make setup
```

## ビルド

```bash
make build
```

実行バイナリ `github-analyzer` がプロジェクトルートに生成されます。

### クリーン

```bash
make clean
```

ビルド成果物を削除します。

## テスト

### ユニットテスト

```bash
make test
```

### E2Eテスト

E2Eテストにはビルドタグ `e2e` が付与されており、通常のテストとは分離されています。

```bash
make test-e2e
```

## コード品質

### フォーマット

```bash
make fmt
```

### lint

```bash
make lint
```

[golangci-lint](https://golangci-lint.run/) を使用しています。

### 依存関係チェック

```bash
make depcheck
```

[go-depcheck](https://github.com/v-standard/go-depcheck) を使用して、クリーンアーキテクチャの層間依存ルールを検証します。ルールは `depcheck.yml` に定義されています。詳細は [architecture.md](./architecture.md) を参照してください。

## Claude Codeワークフロー

本リポジトリでは、Claude Codeを活用したIssue駆動の自動開発ワークフローが構築されています。

### ワークフローの全体像

1. **Issue作成** — `/create-issue` コマンドで対話的にタスクをIssueとして作成
2. **自動選定** — `ready` ラベル付きIssueから優先度の高いものに `assign-to-claude` ラベルを自動付与
3. **自動実装** — `assign-to-claude` ラベル付きIssueを検出し、Claude Codeが自動で実装・PR作成
4. **レビュー・マージ** — CIとAIレビューを経てマージ

### 起動手順

3つのターミナル（またはtmuxペイン）を使用します。

**ターミナル1: キュー監視スクリプト**

```bash
.claude/scripts/watch-empty-queue.sh
```

**ターミナル2: Issue実装スクリプト**

```bash
.claude/scripts/watch-queued-issues.sh
```

**ターミナル3: Issue作成**

```bash
claude
# Claude Code起動後に /create-issue コマンドでタスクを作成
```

### 各スクリプトの役割

| スクリプト | 役割 |
|---|---|
| `watch-empty-queue.sh` | 60秒間隔でキューを監視し、`assign-to-claude` ラベル付きIssueが0件なら `assign-issues.sh` を実行して `ready` ラベル付きIssueから自動選定 |
| `watch-queued-issues.sh` | 60秒間隔で `assign-to-claude` ラベル付きIssueを監視し、検出したら `solve-issue.sh` で自動実装（古いIssueから順に処理） |
| `/create-issue` | Claude Code上で対話的にIssueを作成するスキル |

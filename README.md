# github-analyzer

GitHub上の活動プロセスを可視化・分析するCLI診断ツールです。PRやIssueの活動データをGemini AIで分析し、チームの開発プロセスに関するレポートを生成します。

## クイックスタート

### 前提条件

以下の環境変数を設定してください。

| 環境変数 | 説明 |
|---|---|
| `GH_TOKEN` または `GITHUB_TOKEN` | GitHub APIアクセス用トークン（`GH_TOKEN` が優先） |
| `GEMINI_API_KEY` | Gemini APIキー |

### インストール

[GitHub Releases](https://github.com/canpok1/github-analyzer/releases) からビルド済みバイナリをダウンロードしてください。

対応プラットフォーム:

| OS | アーキテクチャ |
|---|---|
| Linux | amd64, 386, arm64, arm |
| Windows | amd64, 386 |
| macOS | amd64, arm64 |

ダウンロード後、パスの通ったディレクトリに配置してください。

### 設定ファイルの初期化

```bash
github-analyzer init
```

カレントディレクトリに `.github-analyzer.yaml` のテンプレートが生成されます。

### 最初の分析を実行

```bash
# 本日の活動を分析
github-analyzer --repo owner/repo --today
```

## 使い方

### サブコマンド

| コマンド | 説明 |
|---|---|
| `init` | 設定ファイルのテンプレートを生成する |

### フラグ

| フラグ | 説明 |
|---|---|
| `--today` | 本日の活動を分析 |
| `--since` | 指定期間の活動を分析（例: `7d`, `2w`, `1m`） |
| `--pr` | 特定のPR番号を指定して分析 |
| `--issue` | 特定のIssue番号を指定して分析 |
| `--status` | ステータスでフィルタ（`open` / `merged` / `closed`） |
| `--prompt` | 分析の切り口を自由記述 |
| `--repo` | 分析対象リポジトリ（`owner/repo` 形式） |
| `--model` | 使用するGeminiモデル（例: `gemini-2.5-flash`） |
| `-o`, `--output` | レポート出力先ファイルパス（未指定時は標準出力） |

### バリデーションルール

- `--today` と `--since` は同時に指定できません
- `--pr` と `--issue` は同時に指定できません

### 使用例

```bash
# 本日の活動を分析
github-analyzer --repo owner/repo --today

# 直近7日間の活動を分析
github-analyzer --repo owner/repo --since 7d

# 特定のPRを分析
github-analyzer --repo owner/repo --pr 42

# 特定のIssueを分析
github-analyzer --repo owner/repo --issue 10

# カスタムプロンプトで分析
github-analyzer --repo owner/repo --today --prompt "レビュープロセスの改善点を分析してください"

# レポートをファイルに出力
github-analyzer --repo owner/repo --today -o report.md
```

## 設定ファイルリファレンス

設定ファイル `.github-analyzer.yaml` で既定値を設定できます。

### フィールド一覧

| フィールド | 説明 | 例 |
|---|---|---|
| `repo` | 分析対象のGitHubリポジトリ（`owner/repo` 形式） | `owner/repo` |
| `tone` | 分析レポートのトーン | `friendly`, `formal`, `casual` |
| `default_prompt` | デフォルトの分析プロンプト | `チームの活動を分析してください` |
| `model` | 使用するGeminiモデル | `gemini-2.5-flash` |
| `mock.ai` | Gemini APIをモック化 | `true` / `false` |
| `mock.repository` | GitHub APIをモック化 | `true` / `false` |
| `log_file` | ログ出力先ファイルパス | `/path/to/github-analyzer.log` |

### 設定例

```yaml
repo: owner/repo
tone: friendly
default_prompt: チームの活動を分析してください
model: gemini-2.5-flash
log_file: /tmp/github-analyzer.log
```

### 読み込み順と優先順位

設定ファイルは以下の順で読み込まれ、後から読み込まれた値で上書きされます。

1. `~/.github-analyzer.yaml`（ホームディレクトリ）
2. `./.github-analyzer.yaml`（カレントディレクトリ）

CLIフラグ、設定ファイル、デフォルト値の優先順位は以下の通りです。

**CLIフラグ > 設定ファイル > デフォルト値**

## 開発者向けドキュメント

- [アーキテクチャ](docs/architecture.md)
- [開発ガイド](docs/development-guide.md)

# アーキテクチャ

## ディレクトリ構成

```
.
├── cmd/                  # CLIエントリーポイント（DI担当）
├── internal/
│   ├── domain/           # ドメイン層
│   │   └── entity/       # エンティティ定義
│   ├── app/              # アプリケーション層（ユースケース）
│   └── infra/            # インフラストラクチャ層
│       ├── config/       # 設定ファイル読み込み
│       ├── gemini/       # Gemini APIクライアント
│       ├── github/       # GitHub APIクライアント
│       ├── log/          # ログ
│       ├── mock/         # モック
│       └── report/       # レポート出力
├── scripts/              # ユーティリティスクリプト
├── templates/            # テンプレートファイル
└── test/
    └── e2e/              # E2Eテスト
```

## クリーンアーキテクチャ

本プロジェクトはクリーンアーキテクチャに基づき、以下の4層で構成されています。

### 各層の責務

| 層 | パッケージ | 責務 |
|---|---|---|
| domain | `internal/domain/`, `internal/domain/entity/` | エンティティ定義、インターフェース定義。外部への依存を持たない |
| app | `internal/app/` | ユースケースの実装。domain層のみに依存する |
| infra | `internal/infra/` | 外部サービスとの通信、I/O処理。domain層・app層に依存できる |
| cmd | `cmd/` | CLIのエントリーポイント、DI（依存性注入）。全ての層に依存できる |

### 依存関係ルール

内側の層は外側の層に依存してはいけません。

```
cmd → infra → app → domain
(外側)                (内側)
```

具体的な禁止ルール:

- **domain層** → infra層・cmd層への依存禁止
- **app層** → infra層・cmd層への依存禁止
- **infra層** → cmd層への依存禁止

これらのルールは `depcheck.yml` で定義されており、`make depcheck` で自動検証できます。

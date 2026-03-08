# 開発ガイド

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

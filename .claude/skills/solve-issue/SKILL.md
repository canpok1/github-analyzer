---
name: solve-issue
description: ユーザーが `/solve-issue` で手動実行した場合のみ使用。GitHub Issueの対応を行うスキル。実装、自己レビュー、PR作成、マージ、振り返りまで一連の流れを一気に行う。
disable-model-invocation: true
user-invocable: true
argument-hint: "[issue-number]"
---

GitHub Issue $ARGUMENTS を対応します。

1. Issue の内容を理解する
2. `/tdd` スキルで実装する
3. `/simplify` スキルで自己レビューしてコードを改善する
4. lint/formatチェックを実行する（PR作成前の最終ガード）
  - `gofmt -l .` → 出力があれば `gofmt -w .` で修正
  - `golangci-lint run` → 指摘があれば修正
  - 修正した場合はコミットする
5. 同一Issueに対する既存PRの重複チェックを行う
  - コマンド: `gh pr list --search "issue-{番号}" --state all`
  - open または merged のPRが既に存在する場合は、新しいPRを作成せずユーザーに報告して指示を仰ぐ
6. `commit-push-pr` スキルでPRを作成する
7. CIの終了を待機する
  - コマンド: `gh pr checks {PR番号} --watch`
8. `/pr-comments` スキルでレビューコメントを取得し、必要に応じてコードを修正する
  - コードを修正した場合はコミット・プッシュを行いレビューコメントに返信して、手順7に戻る
  - レビューコメントへの返信時は、レビュースレッド内の全レビュワーに対してメンションすること
9. PRをマージする
  - マージできない場合は、原因を確認して必要に応じてコードを修正し、手順7に戻る
10. `/retro` スキルで振り返りを行う

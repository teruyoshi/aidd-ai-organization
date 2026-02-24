# /project:migrate

## 概要
データベースマイグレーションファイルを生成・管理するコマンド。

## パラメータ
- `action` (必須): `create` | `migrate` | `rollback` | `status`
- `name` (省略可): マイグレーション名（`create` 時に必須）
- `steps` (省略可): ロールバックするステップ数（デフォルト: 1）

## 手順

### `create <name>` — マイグレーションファイル作成
1. `backend/internal/model/` 配下の `.go` ファイルを読み込み、既存の構造体定義を把握する
2. タイムスタンプ（`YYYYMMDDHHMMSS`）付きのファイル名を生成する（例: `20240324102030_create_users_table.go`）
3. `backend/migrations/` に Go のマイグレーションファイルを生成する（`Up()`・`Down()` 関数を含むテンプレート）
4. モデル名から適切な AutoMigrate / DropTable のヒントをコメントとして埋め込む

### `migrate` — 未実行のマイグレーションを適用
1. `backend/migrations/` 配下のファイルをタイムスタンプ順に取得する
2. 実行済みマイグレーション（`migration_histories` テーブル）と照合し、未実行分を特定する
3. `go run ./cmd/migrate` 等のコマンドで各マイグレーションを順番に実行する
4. 実行結果を報告する

### `rollback [steps]` — マイグレーションを元に戻す
1. `migration_histories` テーブルから最新 `steps` 件を取得する
2. 対応する `Down()` 関数を呼び出してロールバックを実行する
3. 実行結果を報告する

### `status` — マイグレーション状況確認
1. `backend/migrations/` 配下の全ファイルと実行済みレコードを照合する
2. 各マイグレーションのステータスを一覧で表示する

## 出力フォーマット

```markdown
## 🗄️ データベースマイグレーション結果

### ✅ マイグレーション作成完了
- **ファイル**: backend/migrations/20240324102030_create_users_table.go
- **マイグレーション名**: create_users_table

### 🚀 次のステップ
1. 生成されたマイグレーションファイルを編集してスキーマを定義する
2. `/project:migrate migrate` でマイグレーションを実行する
3. `backend/internal/model/` のモデルを必要に応じて更新する

### 📊 マイグレーション実行結果
- **実行**: X個 / **スキップ**: X個 / **エラー**: X個

### 📋 マイグレーション状況
| マイグレーションID | 名前 | 状態 | 実行日時 |
|-------------------|------|------|----------|
| 20240324101520_initial_schema | initial_schema | ✅ executed | 2024-03-24 10:15 |
| 20240324102030_create_users_table | create_users_table | ⏳ pending | - |

### 💡 推奨事項
- 実行前にデータベースのバックアップを取得してください
- ロールバック処理 (`Down()`) も必ず実装してください

---
*この結果は `/project:migrate` コマンドにより生成されました*
```

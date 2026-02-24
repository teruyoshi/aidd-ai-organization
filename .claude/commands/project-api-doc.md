# /project:api-doc

## 概要
Golang バックエンドのハンドラー・モデルを調査し、OpenAPI 仕様書と Markdown ドキュメントを生成するコマンド。

## パラメータ
- `format` (省略可): `openapi` | `markdown` | `both`（デフォルト）
- `output` (省略可): 出力ディレクトリ（デフォルト: `docs/api`）

## 手順
1. `backend/internal/handler/` 配下の `.go` ファイルを読み込み、HTTPハンドラー関数とルーティング定義を抽出する
2. `backend/internal/model/` 配下の `.go` ファイルを読み込み、構造体定義とフィールド情報（JSONタグ・バリデーション）を抽出する
3. 抽出した情報をもとに以下を生成する:
   - **OpenAPI 仕様書** (`docs/api/openapi.yaml`): エンドポイント・スキーマ・認証情報を含む
   - **Markdown ドキュメント** (`docs/api/api-documentation.md`): エンドポイント一覧・リクエスト/レスポンス例・データモデル定義を含む
4. 生成結果のサマリーと次のステップを提示する

## 出力フォーマット

```markdown
## 📚 APIドキュメント生成結果

### ✅ 生成完了
- **エンドポイント数**: X個
- **データモデル数**: X個
- **出力ファイル**:
  - OpenAPI: docs/api/openapi.yaml
  - Markdown: docs/api/api-documentation.md

### 📋 検出されたエンドポイント

#### Authentication
- POST /api/auth/login
- POST /api/auth/register

#### Users
- GET /api/users
- POST /api/users
- GET /api/users/{id}
- PUT /api/users/{id}
- DELETE /api/users/{id}

### 🏗️ データモデル
- [モデル名] — [説明]

### 🚀 次のステップ
1. 生成された OpenAPI 仕様書を Swagger UI で確認する
2. `/agent-frontend` に型定義を共有する
3. Postman コレクションを生成してAPIテストに活用する

---
*この結果は `/project:api-doc` コマンドにより生成されました*
```

# AIDD TODO Sample - API Documentation

> **REST API仕様書** - TODO管理アプリケーションのバックエンドAPI

## 📋 概要

- **ベースURL**: `http://localhost:8080/api`
- **認証方式**: JWT Bearer Token
- **レスポンス形式**: JSON
- **API バージョン**: v1

## 🔐 認証

### JWT Token Authentication

```http
Authorization: Bearer <jwt_token>
```

すべての保護されたエンドポイントでは、Authorizationヘッダーに有効なJWTトークンが必要です。

## 📝 共通レスポンス形式

### 成功レスポンス
```json
{
  "message": "操作が成功しました",
  "data": {
    // レスポンスデータ
  }
}
```

### エラーレスポンス
```json
{
  "error": "エラーメッセージ",
  "code": "ERROR_CODE",
  "details": {
    // 詳細情報（バリデーションエラーなど）
  }
}
```

## 👤 認証・ユーザー管理

### ユーザー登録

```http
POST /api/auth/register
```

**リクエストボディ:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "name": "山田太郎"
}
```

**レスポンス:** `201 Created`
```json
{
  "message": "ユーザー登録が完了しました",
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "山田太郎",
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### ログイン

```http
POST /api/auth/login
```

**リクエストボディ:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**レスポンス:** `200 OK`
```json
{
  "message": "ログインに成功しました",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "山田太郎",
      "status": "active",
      "profile": {
        "bio": null,
        "location": null,
        "timezone": "Asia/Tokyo",
        "language": "ja"
      }
    }
  }
}
```

### ログアウト

```http
POST /api/auth/logout
Authorization: Bearer <token>
```

**レスポンス:** `200 OK`
```json
{
  "message": "ログアウトしました"
}
```

### プロフィール取得

```http
GET /api/auth/profile
Authorization: Bearer <token>
```

**レスポンス:** `200 OK`
```json
{
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "山田太郎",
    "avatar": "https://example.com/avatar.jpg",
    "status": "active",
    "profile": {
      "bio": "フルスタックエンジニア",
      "location": "東京, 日本",
      "website": "https://mywebsite.com",
      "timezone": "Asia/Tokyo",
      "language": "ja"
    }
  }
}
```

### プロフィール更新

```http
PUT /api/auth/profile
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "name": "山田太郎",
  "bio": "フルスタックエンジニア",
  "location": "東京, 日本",
  "website": "https://mywebsite.com",
  "timezone": "Asia/Tokyo",
  "language": "ja"
}
```

## 📝 TODO管理

### TODO一覧取得

```http
GET /api/todos?category_id=1&status=pending&priority=high&limit=20&offset=0
Authorization: Bearer <token>
```

**クエリパラメータ:**
- `category_id` (optional): カテゴリID
- `status` (optional): ステータス (`pending`, `in_progress`, `completed`, `cancelled`)
- `priority` (optional): 優先度 (`low`, `medium`, `high`, `urgent`)
- `tags` (optional): タグ名（カンマ区切り）
- `search` (optional): 検索文字列
- `completed` (optional): 完了フラグ (`true`, `false`)
- `due_before` (optional): 期限前フィルター (ISO 8601)
- `due_after` (optional): 期限後フィルター (ISO 8601)
- `limit` (optional): 取得件数 (デフォルト: 20)
- `offset` (optional): オフセット (デフォルト: 0)
- `sort_by` (optional): ソートフィールド
- `sort_order` (optional): ソート順 (`asc`, `desc`)

**レスポンス:** `200 OK`
```json
{
  "data": {
    "todos": [
      {
        "id": 1,
        "title": "プロジェクト企画書作成",
        "description": "新プロジェクトの企画書を作成する",
        "status": "pending",
        "priority": "high",
        "due_date": "2024-01-15T23:59:59Z",
        "completed_at": null,
        "position": 0,
        "category": {
          "id": 1,
          "name": "仕事",
          "color": "#3B82F6",
          "icon": "💼"
        },
        "tags": [
          {
            "id": 1,
            "name": "企画",
            "color": "#10B981"
          }
        ],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 50,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

### TODO詳細取得

```http
GET /api/todos/:id
Authorization: Bearer <token>
```

**レスポンス:** `200 OK`
```json
{
  "data": {
    "id": 1,
    "title": "プロジェクト企画書作成",
    "description": "新プロジェクトの企画書を作成する",
    "status": "pending",
    "priority": "high",
    "due_date": "2024-01-15T23:59:59Z",
    "category": {
      "id": 1,
      "name": "仕事",
      "color": "#3B82F6"
    },
    "tags": [
      {
        "id": 1,
        "name": "企画",
        "color": "#10B981"
      }
    ],
    "attachments": [
      {
        "id": 1,
        "file_name": "proposal_draft.pdf",
        "original_name": "企画書ドラフト.pdf",
        "mime_type": "application/pdf",
        "file_size": 1048576,
        "uploaded_by": 1,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "comments": [
      {
        "id": 1,
        "content": "期限を2日前倒しして進めます",
        "user": {
          "id": 1,
          "name": "山田太郎"
        },
        "created_at": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

### TODO作成

```http
POST /api/todos
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "category_id": 1,
  "title": "新しいタスク",
  "description": "タスクの詳細説明",
  "priority": "medium",
  "due_date": "2024-01-20T23:59:59Z",
  "tags": ["開発", "レビュー"]
}
```

**レスポンス:** `201 Created`
```json
{
  "message": "TODOが作成されました",
  "data": {
    "id": 2,
    "title": "新しいタスク",
    "status": "pending",
    "priority": "medium",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### TODO更新

```http
PUT /api/todos/:id
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "title": "更新されたタスク",
  "description": "更新された詳細",
  "status": "in_progress",
  "priority": "high",
  "due_date": "2024-01-25T23:59:59Z"
}
```

### TODO削除

```http
DELETE /api/todos/:id
Authorization: Bearer <token>
```

**レスポンス:** `200 OK`
```json
{
  "message": "TODOが削除されました"
}
```

### TODO一括更新

```http
POST /api/todos/bulk-update
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "todo_ids": [1, 2, 3],
  "status": "completed"
}
```

### TODO統計情報

```http
GET /api/todos/stats
Authorization: Bearer <token>
```

**レスポンス:** `200 OK`
```json
{
  "data": {
    "total": 100,
    "completed": 75,
    "in_progress": 15,
    "pending": 8,
    "cancelled": 2,
    "overdue": 3,
    "completion_rate": 75.0
  }
}
```

## 📁 カテゴリ管理

### カテゴリ一覧取得

```http
GET /api/categories
Authorization: Bearer <token>
```

**レスポンス:** `200 OK`
```json
{
  "data": [
    {
      "id": 1,
      "name": "仕事",
      "description": "業務関連のタスク",
      "color": "#3B82F6",
      "icon": "💼",
      "position": 0,
      "is_default": true,
      "todo_count": 25,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### カテゴリ作成

```http
POST /api/categories
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "name": "プライベート",
  "description": "個人的なタスク",
  "color": "#10B981",
  "icon": "🏠",
  "position": 1
}
```

### カテゴリ更新

```http
PUT /api/categories/:id
Authorization: Bearer <token>
```

### カテゴリ削除

```http
DELETE /api/categories/:id
Authorization: Bearer <token>
```

### カテゴリ並び替え

```http
POST /api/categories/reorder
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "category_ids": [2, 1, 3]
}
```

## 🏷️ タグ管理

### タグ一覧取得

```http
GET /api/tags?search=開発
Authorization: Bearer <token>
```

### タグ作成

```http
POST /api/tags
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "name": "重要",
  "color": "#EF4444",
  "description": "重要度の高いタスク"
}
```

### タグ更新

```http
PUT /api/tags/:id
Authorization: Bearer <token>
```

### タグ削除

```http
DELETE /api/tags/:id
Authorization: Bearer <token>
```

### 人気タグ取得

```http
GET /api/tags/popular?limit=10
Authorization: Bearer <token>
```

## 📎 添付ファイル

### ファイルアップロード

```http
POST /api/todos/:todo_id/attachments
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**リクエストボディ (multipart/form-data):**
- `file`: ファイル（最大50MB）

### 添付ファイル一覧

```http
GET /api/todos/:todo_id/attachments
Authorization: Bearer <token>
```

### 添付ファイル削除

```http
DELETE /api/attachments/:id
Authorization: Bearer <token>
```

## 💬 コメント

### コメント一覧取得

```http
GET /api/todos/:todo_id/comments
Authorization: Bearer <token>
```

### コメント作成

```http
POST /api/todos/:todo_id/comments
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "content": "コメントの内容"
}
```

### コメント更新

```http
PUT /api/comments/:id
Authorization: Bearer <token>
```

### コメント削除

```http
DELETE /api/comments/:id
Authorization: Bearer <token>
```

## 📊 ダッシュボード

### ダッシュボードデータ取得

```http
GET /api/dashboard
Authorization: Bearer <token>
```

**レスポンス:** `200 OK`
```json
{
  "data": {
    "stats": {
      "total": 100,
      "completed": 75,
      "completion_rate": 75.0
    },
    "overdue_todos": [...],
    "upcoming_todos": [...],
    "recent_todos": [...],
    "popular_tags": [...],
    "categories_with_count": [...]
  }
}
```

## 🔧 システム

### ヘルスチェック

```http
GET /api/health
```

**レスポンス:** `200 OK`
```json
{
  "status": "ok",
  "service": "todo-backend",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### システム情報

```http
GET /api/info
Authorization: Bearer <token>
```

## 📋 ステータスコード

| コード | 説明 |
|--------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 400 | リクエストエラー |
| 401 | 認証エラー |
| 403 | 認可エラー |
| 404 | リソースが見つからない |
| 409 | 競合エラー |
| 422 | バリデーションエラー |
| 500 | サーバーエラー |

## 🔒 セキュリティ

### レート制限

- **認証エンドポイント**: 5リクエスト/分
- **API一般**: 10リクエスト/秒
- **ファイルアップロード**: 5リクエスト/分

### CORS設定

開発環境では `http://localhost:3000` からのアクセスを許可

### セキュリティヘッダー

- `X-Frame-Options: SAMEORIGIN`
- `X-XSS-Protection: 1; mode=block`
- `X-Content-Type-Options: nosniff`

## 🧪 テスト用データ

### テストユーザー

```json
{
  "email": "test@example.com",
  "password": "testpassword123",
  "name": "テストユーザー"
}
```

### サンプルリクエスト (curl)

```bash
# ユーザー登録
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpassword123","name":"テストユーザー"}'

# ログイン
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpassword123"}'

# TODO作成
curl -X POST http://localhost:8080/api/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"テストタスク","priority":"medium"}'
```

---

> **開発者向け注意事項**
>
> - 本API仕様書は開発中のため、変更される可能性があります
> - 本番環境では適切なHTTPS設定とセキュリティ対策を実施してください
> - レート制限の値は環境に応じて調整してください
# AIDD TODO Sample - API Documentation

> **REST API仕様書** - TODO管理アプリケーションのバックエンドAPI

## 概要

- **ベースURL**: `http://localhost:8080/api`
- **認証方式**: JWT Bearer Token
- **レスポンス形式**: JSON

## 認証

すべての保護されたエンドポイントには `Authorization: Bearer <token>` ヘッダーが必要です。

## エンドポイント

### ユーザー登録

```http
POST /api/auth/register
```

**リクエストボディ:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "山田太郎"
}
```

**レスポンス:** `201 Created`
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "山田太郎",
  "created_at": "2024-01-01T00:00:00Z"
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
  "password": "password123"
}
```

**レスポンス:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "山田太郎"
  }
}
```

### TODO一覧取得

```http
GET /api/todos
Authorization: Bearer <token>
```

**レスポンス:** `200 OK`
```json
[
  {
    "id": 1,
    "title": "プロジェクト企画書作成",
    "description": "新プロジェクトの企画書を作成する",
    "done": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### TODO作成

```http
POST /api/todos
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "title": "新しいタスク",
  "description": "タスクの詳細説明"
}
```

**レスポンス:** `201 Created`

### TODO更新

```http
PUT /api/todos/:id
Authorization: Bearer <token>
```

**リクエストボディ:**
```json
{
  "title": "更新されたタスク",
  "done": true
}
```

### TODO削除

```http
DELETE /api/todos/:id
Authorization: Bearer <token>
```

**レスポンス:** `204 No Content`

### ヘルスチェック

```http
GET /api/health
```

**レスポンス:** `200 OK`
```json
{
  "status": "ok"
}
```

## ステータスコード

| コード | 説明 |
|--------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 204 | 削除成功 |
| 400 | リクエストエラー |
| 401 | 認証エラー |
| 403 | 認可エラー |
| 404 | リソースが見つからない |
| 500 | サーバーエラー |

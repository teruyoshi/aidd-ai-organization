# TODOアプリ - AIDD実装サンプル

> **AI-Driven Development** 組織によるシンプルなTODOアプリケーション

AIDD組織の実践例として開発された、シンプルなTODOアプリケーションです。

## 技術スタック

- **フロントエンド**: React 18 + TypeScript + Vite + Tailwind CSS + TanStack Query
- **バックエンド**: Go + go-chi/chi + GORM
- **データベース**: MySQL 8.0
- **インフラ**: Docker + Docker Compose

## 機能

- ユーザー認証（登録・ログイン）
- TODO管理（作成・表示・更新・削除）
- TODO完了状態の切り替え

## クイックスタート

```bash
# 開発環境起動
./scripts/sample-dev.sh

# または直接
docker-compose -f sample-docker-compose.yml up -d --build
```

### アクセス先

- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080

## 開発ガイド

### フロントエンド

```bash
cd frontend
npm install
npm run dev
```

### バックエンド

```bash
cd backend
go run cmd/server/sample-main.go
```

## プロジェクト構造

```
aidd-ai-organization/
├── frontend/                    # React + TypeScript
│   ├── src/
│   │   ├── features/           # 機能別（auth, todos）
│   │   ├── components/         # 共通コンポーネント
│   │   └── types/              # TypeScript型定義
│   └── sample-package.json
├── backend/                    # Go + chi + GORM
│   ├── cmd/server/             # エントリーポイント
│   └── internal/
│       ├── domain/             # エンティティ・リポジトリIF
│       ├── usecase/            # ビジネスロジック
│       ├── interface/          # HTTPハンドラー
│       └── infrastructure/     # DB実装
├── sample-docker-compose.yml   # 開発環境
└── docs/sample-API.md          # API仕様書
```

## API仕様

詳細は [docs/sample-API.md](./docs/sample-API.md) を参照してください。

---

**AIDD組織によるサンプル実装 — 詳細は [CLAUDE.md](./CLAUDE.md) を参照**

# AIDD（AI-Driven Development）AI組織構造

このリポジトリは **AIDD（AI-Driven Development）** を実践するためのAI組織構造を定義しています。
Claude Code を使って、専門化されたSubAgentが分業しながらアプリを開発します。

## クイックスタート

```
/project:plan
```

→ SCHEDULE.md を読み込んで次のタスクを提示してくれます。

## 各エージェントへの指示

```bash
# アーキテクチャ設計
Read agents/architect.md and design the DB schema for [機能名]

# フロントエンド実装
Read agents/frontend.md and docs/context/frontend-context.md, then implement [コンポーネント名]

# バックエンド実装
Read agents/backend.md and docs/context/backend-context.md, then implement [API名]

# インフラ構築
Read agents/devops.md and docs/context/devops-context.md, then create docker-compose.dev.yml

# テスト作成
Read agents/qa.md and write tests for [対象ファイル]
```

## AI組織構成

```
📁 AI組織
├── 🏛️  アーキテクト    — 設計・技術選定・タスク分解
├── 🎨  フロントエンド  — React + TypeScript + Vite
├── ⚙️  バックエンド    — Golang + GORM + MySQL
├── 🐳  DevOps         — Docker + Docker Compose
└── 🧪  QA             — テスト戦略・テストコード
```

## スラッシュコマンド一覧

| コマンド | 説明 |
|----------|------|
| `/project:plan` | SCHEDULE.mdを確認して次タスクを分解 |
| `/project:review` | コードレビュー実行 |
| `/project:api-doc` | Golang APIドキュメント自動生成 |
| `/project:migrate` | DBマイグレーションファイル生成 |

## ファイル構成

```
.
├── CLAUDE.md                    # 組織憲法（必ず最初に読まれる）
├── SCHEDULE.md                # タスク・進捗管理
├── .claude/
│   ├── commands/                # スラッシュコマンド定義
│   └── settings.json           # MCP設定
├── agents/                     # SubAgent定義
│   ├── architect.md
│   ├── frontend.md
│   ├── backend.md
│   ├── devops.md
│   └── qa.md
└── docs/context/               # 技術コンテキスト
    ├── frontend-context.md
    ├── backend-context.md
    └── devops-context.md
```

## 技術スタック

- **コンテナ**: Docker + Docker Compose
- **フロントエンド**: Vite + React + TypeScript
- **バックエンド**: Golang + go-chi/chi
- **ORM**: GORM
- **DB**: MySQL 8.0

## 進捗管理規則

SCHEDULE.md の記法：

```
[ ]  未着手
[-]  進行中 🏗️
[x]  完了 ✅
```

## 使用方法

1. `CLAUDE.md` を確認して組織憲法を理解
2. `/project:plan` でタスクを確認
3. 各エージェント指示でSubAgentに作業依頼
4. `SCHEDULE.md` で進捗を管理
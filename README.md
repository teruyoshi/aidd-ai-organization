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
/agent-architect [機能名]のDBスキーマを設計してください

# フロントエンド実装
/agent-frontend [コンポーネント名]を実装してください

# バックエンド実装
/agent-backend [API名]を実装してください

# インフラ構築
/agent-devops docker-compose.dev.yml を作成してください

# テスト作成
/agent-qa [対象ファイル]のテストを作成してください
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

## Skills・コマンド一覧

### エージェント Skills

| スキル | 説明 |
|--------|------|
| `/agent-architect` | 設計・技術選定・タスク分解 |
| `/agent-frontend` | React + TypeScript + Vite 実装 |
| `/agent-backend` | Go + GORM + MySQL 実装 |
| `/agent-devops` | Docker・CI/CD・Nginx 設定 |
| `/agent-qa` | テスト作成・品質確認 |

### プロジェクト管理コマンド

| コマンド | 説明 |
|----------|------|
| `/project:plan` | SCHEDULE.mdを確認して次タスクを分解 |
| `/project:review` | コードレビュー実行 |
| `/project:api-doc` | Golang APIドキュメント自動生成 |
| `/project:migrate` | DBマイグレーションファイル生成 |

## ファイル構成

```
.
├── CLAUDE.md                    # 組織憲法（常時読み込み）
├── SCHEDULE.md                  # タスク・進捗管理
├── .claude/
│   ├── commands/                # プロジェクト管理コマンド（4個）
│   │   ├── project-plan.md
│   │   ├── project-review.md
│   │   ├── project-api-doc.md
│   │   └── project-migrate.md
│   ├── skills/                  # エージェントSkills・自動知識（8個）
│   │   ├── agent-architect/
│   │   ├── agent-frontend/
│   │   ├── agent-backend/
│   │   ├── agent-devops/
│   │   ├── agent-qa/
│   │   ├── react-coding-standards/
│   │   ├── go-coding-standards/
│   │   └── schedule-tracking/
│   └── settings.json            # 権限・フック設定
├── agents/                      # エージェント詳細定義（Skillsから参照）
│   ├── architect.md
│   ├── frontend.md
│   ├── backend.md
│   ├── devops.md
│   └── qa.md
└── docs/context/                # 技術コンテキスト（Skillsから参照）
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
3. `/agent-xxx` Skills でサブエージェントに作業を委譲
4. `SCHEDULE.md` で進捗を管理
# AIDD AI組織憲法

> **AI-Driven Development (AIDD)** 組織運営の基本原則・方針

## 🎯 組織ミッション

Claude Code を活用し、専門化されたAIエージェントが分業・連携してフルスタックアプリケーションを効率的に開発する。

## 🏗️ 組織構造

### 5エージェント分業体制

```
📁 AIDD AI組織
├── 🏛️ アーキテクト     — システム設計・技術選定・タスク分解
├── 🎨 フロントエンド   — React + TypeScript + Vite
├── ⚙️ バックエンド     — Golang + go-chi/chi + GORM + MySQL
├── 🐳 DevOps          — Docker + Docker Compose + CI/CD
└── 🧪 QA             — テスト戦略・品質保証
```

### 各エージェントの責任範囲

| エージェント | 主要責任 | 成果物 |
|-------------|----------|--------|
| **アーキテクト** | 要件分析・設計・技術選定 | 設計書・API仕様・タスク分解 |
| **フロントエンド** | UI/UX・コンポーネント実装 | Reactアプリ・TypeScript型定義 |
| **バックエンド** | API・DB・ビジネスロジック | Go API・GORM モデル・認証 |
| **DevOps** | インフラ・デプロイ・監視 | Docker環境・CI/CD・監視設定 |
| **QA** | テスト・品質管理 | テストコード・品質レポート |

## 📋 開発プロセス

### 1. 要件定義・設計フェーズ
```bash
# アーキテクトが要件を分析し、各エージェントにタスク分解
Read agents/architect.md and analyze requirements, then create implementation tasks
```

### 2. 並行開発フェーズ
```bash
# フロントエンド実装
Read agents/frontend.md and docs/context/frontend-context.md, then implement [component]

# バックエンド実装
Read agents/backend.md and docs/context/backend-context.md, then implement [API]

# インフラ構築
Read agents/devops.md and docs/context/devops-context.md, then setup [environment]
```

### 3. 統合・テストフェーズ
```bash
# QAが全体テスト実行
Read agents/qa.md and create comprehensive tests for the application
```

## 🛠️ 技術スタック標準

### フロントエンド標準
- **フレームワーク**: React 18+ (Hooks, Concurrent Features)
- **言語**: TypeScript (strict mode)
- **ビルドツール**: Vite
- **状態管理**: Context API + TanStack Query
- **テスト**: Vitest + React Testing Library + Playwright

### バックエンド標準
- **言語**: Go 1.21+
- **フレームワーク**: go-chi/chi v5
- **ORM**: GORM v2
- **データベース**: MySQL 8.0
- **テスト**: Go testing + testify + gomock

### インフラ標準
- **コンテナ**: Docker + Docker Compose
- **プロキシ**: Nginx
- **監視**: Prometheus + Grafana
- **CI/CD**: GitHub Actions

## 📊 品質基準

### コード品質
- **型安全性**: TypeScript strict mode・Go 型安全性100%
- **テストカバレッジ**: 単体テスト80%以上・E2Eテスト主要フロー100%
- **静的解析**: ESLint・golangci-lint エラー0件
- **セキュリティ**: OWASP対策・認証認可実装

### パフォーマンス
- **フロントエンド**: Lighthouse Score 90+・Core Web Vitals緑
- **バックエンド**: レスポンス時間95%ile < 500ms
- **データベース**: クエリ実行時間 < 100ms

### 運用品質
- **可用性**: アップタイム99.9%
- **監視**: メトリクス・ログ・アラート完備
- **バックアップ**: 日次自動バックアップ・7日保持

## 🔄 作業フロー

### タスク管理規則
```markdown
SCHEDULE.md のステータス管理：
[ ]  未着手
[-]  進行中 🏗️  ← 必ずactiveFormを更新
[x]  完了 ✅    ← 成果物確認後にマーク
```

### エージェント間連携
1. **設計情報共有**: アーキテクト → 各実装エージェント
2. **API仕様共有**: バックエンド → フロントエンド
3. **型定義共有**: フロントエンド ⇄ バックエンド
4. **インフラ情報共有**: DevOps → 全エージェント
5. **品質情報共有**: QA → 全エージェント

### コードレビュー原則
- **設計レビュー**: アーキテクトが全体設計を承認
- **実装レビュー**: 各エージェントが専門領域をレビュー
- **品質レビュー**: QAが最終品質を確認

## 🚀 スラッシュコマンド活用

| コマンド | 目的 | 実行タイミング |
|----------|------|---------------|
| `/project:plan` | 次タスク確認・分解 | 開発開始時・マイルストーン完了時 |
| `/project:review` | コードレビュー実行 | 機能実装完了時 |
| `/project:api-doc` | APIドキュメント生成 | バックエンド実装完了時 |
| `/project:migrate` | DBマイグレーション | データモデル変更時 |

## 🎯 成功指標

### 開発効率
- **開発速度**: 機能当たり開発時間 50%短縮
- **バグ率**: 本番バグ発生率 < 1%
- **リリース頻度**: 週次リリース可能

### チーム効率
- **知識共有**: 全エージェントが技術スタック理解
- **自動化率**: 手動作業 < 10%
- **ドキュメント**: 設計・運用ドキュメント完備

## ⚠️ 重要な制約・原則

### 技術原則
1. **一貫性優先**: 技術スタック・コーディング規約の統一
2. **品質優先**: 機能より品質・保守性を重視
3. **自動化優先**: 手動作業の徹底排除
4. **セキュリティ優先**: セキュリティを後回しにしない

### 開発原則
1. **Clean Code**: 可読性・保守性の高いコード
2. **Test First**: テスト駆動開発の実践
3. **Documentation**: 設計・実装の文書化
4. **Continuous Improvement**: 継続的な改善

### 連携原則
1. **情報共有**: エージェント間の積極的な情報共有
2. **責任明確化**: 各エージェントの責任範囲明確化
3. **品質責任**: 全エージェントが品質に責任を持つ
4. **学習促進**: 技術・ベストプラクティスの共有

---

> この憲法は AIDD組織の基本方針です。全てのエージェントはこの原則に従って開発を進めてください。
# TODOアプリ - AIDD実装サンプル

> **AI-Driven Development** 組織による フル機能TODOアプリケーション

このプロジェクトは、AIDD（AI-Driven Development）組織の実践例として開発されたTODOアプリケーションです。React + TypeScript + Go + MySQL の技術スタックで、5つの専門化されたAIエージェントによる分業開発を実証しています。

## 🎯 プロジェクト概要

### 主要機能
- **👤 ユーザー認証**: JWT認証・登録・ログイン・プロフィール管理
- **✅ タスク管理**: 作成・編集・削除・完了状態・優先度・期限設定
- **🏷️ カテゴリ管理**: タスク分類・色分け・アイコン設定
- **🔍 高度検索**: 条件絞り込み・並び替え・保存検索
- **📊 進捗可視化**: 完了率・統計・レポート・チャート
- **📱 レスポンシブ**: モバイル・タブレット・デスクトップ完全対応

### AIDD組織による開発
- **🏛️ アーキテクト**: システム設計・API仕様・データベース設計
- **🎨 フロントエンド**: React UI・TypeScript・状態管理
- **⚙️ バックエンド**: Go API・GORM・認証システム
- **🐳 DevOps**: Docker環境・Nginx・監視設定
- **🧪 QA**: テスト戦略・品質保証・自動化

## 🛠️ 技術スタック

### フロントエンド
- **React 18.2+** - Concurrent Features・Suspense
- **TypeScript 5.0+** - 厳密型チェック・型安全性
- **Vite 5.0+** - 高速ビルド・HMR・最適化
- **TanStack Query** - サーバー状態管理・キャッシュ
- **Tailwind CSS** - ユーティリティファーストCSS
- **React Hook Form** - 効率的フォーム管理
- **Zustand** - 軽量状態管理

### バックエンド
- **Go 1.21+** - 高性能・型安全・並行処理
- **go-chi/chi v5** - 軽量HTTP ルーター
- **GORM v2** - 強力なORM・マイグレーション
- **MySQL 8.0** - リレーショナルデータベース
- **Redis** - セッション・キャッシュ・レート制限
- **JWT** - ステートレス認証

### インフラ・DevOps
- **Docker + Compose** - コンテナ化・環境統一
- **Nginx** - リバースプロキシ・SSL・静的ファイル
- **Prometheus + Grafana** - メトリクス・監視・アラート

### テスト・品質
- **Vitest** - 高速ユニットテスト（フロントエンド）
- **Playwright** - E2Eテスト・クロスブラウザ
- **Go testing** - 標準テストフレームワーク
- **testify** - アサーション・モック
- **ESLint + Prettier** - コード品質・フォーマット
- **golangci-lint** - Go 静的解析

## 🚀 クイックスタート

### 前提条件
- **Docker & Docker Compose** - 最新版
- **Node.js** - v18以上（開発時のみ）
- **Go** - v1.21以上（開発時のみ）

### 1. リポジトリクローン
```bash
git clone <repository-url>
cd aidd-ai-organization
```

### 2. 開発環境起動
```bash
# Docker環境で全サービス起動
./scripts/sample-dev.sh

# または手動起動
docker-compose -f sample-docker-compose.yml up -d
```

### 3. アクセス
- **フロントエンド**: http://localhost:3000
- **バックエンドAPI**: http://localhost:8080
- **Grafana監視**: http://localhost:3001 (admin/admin)

### 4. 初期データ投入
```bash
# サンプルユーザー・TODOデータ作成
docker-compose -f sample-docker-compose.yml exec backend ./sample-seed.sh
```

## 📁 プロジェクト構造

```
aidd-ai-organization/
├── frontend/                    # React + TypeScript フロントエンド
│   ├── src/
│   │   ├── features/           # 機能別（auth, todos）
│   │   ├── components/         # 再利用コンポーネント
│   │   ├── hooks/              # カスタムフック
│   │   ├── services/           # API通信
│   │   └── types/              # TypeScript型定義
│   ├── sample-package.json     # 依存関係・スクリプト
│   ├── sample-vite.config.ts   # Vite設定
│   └── sample-tailwind.config.js # Tailwind設定
├── backend/                    # Go + chi + GORM バックエンド
│   ├── cmd/server/             # アプリケーションエントリー
│   ├── internal/               # Clean Architecture
│   │   ├── domain/             # エンティティ・リポジトリIF
│   │   ├── usecase/            # ビジネスロジック・DTO
│   │   ├── interface/          # ハンドラー・ミドルウェア
│   │   └── infrastructure/     # DB実装・外部サービス
│   ├── pkg/                    # 共通パッケージ
│   ├── migrations/             # DBマイグレーション
│   └── sample-go.mod           # Go モジュール定義
├── nginx/                      # リバースプロキシ設定
├── scripts/                    # 開発・運用スクリプト
├── docs/api/                   # API仕様書・ドキュメント
├── sample-docker-compose.yml   # 開発環境構成
└── sample-docker-compose.prod.yml # 本番環境構成
```

## 🔧 開発ガイド

### フロントエンド開発
```bash
cd frontend

# 依存関係インストール
npm install

# 開発サーバー起動
npm run dev

# ビルド
npm run build

# テスト実行
npm run test

# 型チェック
npm run type-check

# Lint
npm run lint
```

### バックエンド開発
```bash
cd backend

# 依存関係ダウンロード
go mod download

# 開発サーバー起動
go run cmd/server/sample-main.go

# ビルド
go build -o bin/server cmd/server/sample-main.go

# テスト実行
go test ./...

# カバレッジ
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Lint
golangci-lint run
```

### データベースマイグレーション
```bash
# 新しいマイグレーション作成
/project:migrate create add_new_column

# マイグレーション実行
/project:migrate migrate

# ロールバック
/project:migrate rollback

# ステータス確認
/project:migrate status
```

## 🧪 テスト

### 単体テスト
```bash
# フロントエンド単体テスト
cd frontend && npm run test

# バックエンド単体テスト
cd backend && go test ./...

# カバレッジ付きテスト
npm run test:coverage
go test -coverprofile=coverage.out ./...
```

### E2Eテスト
```bash
# Playwright E2Eテスト
cd frontend && npm run test:e2e

# 特定のテストファイル実行
npx playwright test sample-auth.spec.ts

# ヘッドレスモード無効（ブラウザ表示）
npx playwright test --headed
```

### 統合テスト
```bash
# 全体統合テスト実行
./scripts/sample-test.sh

# API統合テスト
cd backend && go test -tags=integration ./...
```

## 📊 監視・メトリクス

### Grafanaダッシュボード
- **アプリケーション**: http://localhost:3001
  - レスポンス時間・エラー率・スループット
  - データベース接続・クエリパフォーマンス
  - メモリ・CPU使用率

### ログ確認
```bash
# 全サービスログ
docker-compose -f sample-docker-compose.yml logs -f

# 特定サービスログ
docker-compose -f sample-docker-compose.yml logs -f frontend
docker-compose -f sample-docker-compose.yml logs -f backend

# リアルタイム監視
./scripts/sample-monitor.sh
```

## 🛡️ セキュリティ

### 実装済みセキュリティ対策
- **JWT認証**: 安全なトークンベース認証
- **HTTPS強制**: 全通信の暗号化
- **CORS設定**: 適切なクロスオリジン制御
- **レート制限**: API乱用防止
- **入力検証**: XSS・SQLインジェクション対策
- **セキュリティヘッダー**: CSP・HSTS等

### セキュリティテスト
```bash
# 脆弱性スキャン
npm audit
go list -json -m all | nancy sleuth

# セキュリティテスト実行
./scripts/sample-security-test.sh
```

## 🚀 デプロイ

### ステージング環境
```bash
# ステージングビルド・デプロイ
./scripts/sample-deploy.sh staging

# ステージング環境確認
curl -f https://staging.example.com/health
```

### 本番環境
```bash
# 本番ビルド・デプロイ
./scripts/sample-deploy.sh production

# 本番環境ヘルスチェック
./scripts/sample-health-check.sh production
```

## 🤝 開発ワークフロー

### AIDD組織での開発手順
1. **要件分析**: アーキテクトが仕様を分析・設計
2. **タスク分解**: `/project:plan` でタスク確認・分解
3. **並行開発**: 各エージェントが担当領域を実装
4. **コードレビュー**: `/project:review` で品質チェック
5. **API仕様更新**: `/project:api-doc` でドキュメント自動生成
6. **データベース変更**: `/project:migrate` でマイグレーション管理

### ブランチ戦略
- **main**: 本番リリース用
- **develop**: 開発統合用
- **feature/**: 機能開発用
- **hotfix/**: 緊急修正用

## 📚 学習リソース

### AIDD組織について
- [CLAUDE.md](./CLAUDE.md) - 組織憲法・基本方針
- [agents/](./agents/) - 各エージェントの役割・責任
- [docs/context/](./docs/context/) - 技術コンテキスト詳細

### API仕様
- [docs/api/sample-openapi.yaml](./docs/api/sample-openapi.yaml) - OpenAPI 3.0仕様
- [docs/api/sample-api-documentation.md](./docs/api/sample-api-documentation.md) - API利用ガイド

## 🐛 トラブルシューティング

### よくある問題

#### Docker起動エラー
```bash
# ポート競合の場合
docker ps  # 使用中ポート確認
docker-compose -f sample-docker-compose.yml down  # 停止

# 権限エラーの場合
sudo chown -R $USER:$USER .
```

#### フロントエンドビルドエラー
```bash
# Node.jsバージョン確認
node --version  # v18以上必要

# 依存関係クリーンインストール
rm -rf node_modules package-lock.json
npm install
```

#### バックエンドコンパイルエラー
```bash
# Goバージョン確認
go version  # v1.21以上必要

# モジュールクリーンアップ
go mod tidy
go mod download
```

### サポート・問い合わせ
- **Issues**: GitHub Issues で問題報告
- **Discussions**: 機能要望・質問
- **Documentation**: 詳細ドキュメントは [docs/](./docs/) を参照

## 🎉 貢献

このサンプルプロジェクトへの貢献を歓迎します！

1. Fork してください
2. Feature ブランチを作成 (`git checkout -b feature/AmazingFeature`)
3. 変更をコミット (`git commit -m 'Add some AmazingFeature'`)
4. ブランチにプッシュ (`git push origin feature/AmazingFeature`)
5. Pull Request を開いてください

## 📄 ライセンス

このプロジェクトは [MIT License](LICENSE) のもとで公開されています。

## 🙏 謝辞

- **Claude Code**: AI-Driven Development の実現
- **AIDD組織**: 革新的な開発手法の提案
- **オープンソースコミュニティ**: 素晴らしいツール・ライブラリの提供

---

**AIDD組織によるフルスタック開発の実践例として、このTODOアプリをご活用ください！**
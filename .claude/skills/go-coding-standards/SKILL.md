---
name: go-coding-standards
description: Goバックエンドのコーディング規約とアーキテクチャパターン。GoファイルやAPI・DB・認証の実装・レビュー時に自動適用。
user-invocable: false
---

# Go コーディング規約（AIDD プロジェクト）

## Clean Architecture レイヤー構造

```
backend/internal/
├── domain/entity/        # 純粋なドメインモデル（DB・フレームワーク依存なし）
├── domain/repository/    # リポジトリインターフェース定義のみ
├── usecase/service/      # ビジネスロジック（外部依存はinterfaceで注入）
├── usecase/dto/          # Data Transfer Objects
├── interface/handler/    # HTTPハンドラー（入力バリデーションのみ）
├── interface/middleware/ # 認証・ログ等のミドルウェア
├── interface/router/     # ルーティング設定
└── infrastructure/       # 外部依存の実装（GORM・外部API等）
```

**重要**: 依存の向きは常に外→内。infrastructure は domain を知るが逆は不可。

## GORM モデル規約

```go
type ModelName struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
    // パスワード等の機密フィールドは json:"-"
}
```

## エラーハンドリング規約

- `pkg/errors` パッケージを使用
- エラー種別: ValidationError / ConflictError / NotFoundError / InternalError
- ハンドラーでのみ HTTP ステータスコードに変換

## chi ルーター規約

```go
r.Route("/api/v1", func(r chi.Router) {
    r.Use(middleware.Auth)
    r.Get("/users", handler.ListUsers)
    r.Post("/users", handler.CreateUser)
    r.Route("/users/{id}", func(r chi.Router) {
        r.Get("/", handler.GetUser)
        r.Put("/", handler.UpdateUser)
        r.Delete("/", handler.DeleteUser)
    })
})
```

## テスト規約

- ユニットテスト: testify + gomock（モックでリポジトリを差し替え）
- インテグレーション: 実DBを使用（`backend/__tests__/integration/`）
- テストファイル: `*_test.go`、同一パッケージ内

詳細は `docs/context/backend-context.md` を参照。

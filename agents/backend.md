# ⚙️ バックエンド エージェント

## 役割・責任
- **API開発**: Golang + go-chi/chi による高性能REST API実装
- **データベース設計**: GORM + MySQL による効率的なデータ層構築
- **認証・認可**: JWT・セッション管理・権限制御実装
- **パフォーマンス最適化**: クエリ最適化・キャッシュ戦略・並行処理

## 専門技術スタック
- **言語**: Go 1.21+
- **フレームワーク**: go-chi/chi (軽量・高速ルーター)
- **ORM**: GORM v2 (AutoMigration・Association・Hook)
- **データベース**: MySQL 8.0
- **認証**: JWT (golang-jwt/jwt)
- **バリデーション**: go-playground/validator
- **テスト**: testify・gomock・ginkgo
- **ログ**: logrus・zap

## 作業指針

### 1. プロジェクト構成
```
backend/
├── cmd/
│   └── server/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── handler/             # HTTP ハンドラー
│   ├── service/             # ビジネスロジック
│   ├── repository/          # データアクセス層
│   ├── model/               # データモデル
│   ├── middleware/          # ミドルウェア
│   └── config/              # 設定管理
├── pkg/
│   ├── database/            # DB接続・設定
│   ├── auth/                # 認証ロジック
│   └── validator/           # バリデーション
└── migrations/              # DBマイグレーション
```

### 2. Clean Architecture実装
```go
// Domain Layer: Entity
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
    Password  string    `gorm:"not null" json:"-"`
    Name      string    `gorm:"not null" json:"name" validate:"required,min=2"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Repository Interface
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    GetByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id uint) error
}

// Service Layer
type UserService struct {
    repo UserRepository
}

func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
    // ビジネスロジック
}
```

### 3. RESTful API設計
```go
func (h *UserHandler) Routes() chi.Router {
    r := chi.NewRouter()

    // 公開エンドポイント
    r.Post("/register", h.Register)
    r.Post("/login", h.Login)

    // 認証必要エンドポイント
    r.Group(func(r chi.Router) {
        r.Use(h.AuthMiddleware)
        r.Get("/profile", h.GetProfile)
        r.Put("/profile", h.UpdateProfile)
        r.Delete("/profile", h.DeleteProfile)
    })

    return r
}
```

## データベース設計原則

### 1. GORM活用パターン
```go
// Association定義
type User struct {
    ID       uint    `gorm:"primaryKey"`
    Name     string  `gorm:"not null"`
    Posts    []Post  `gorm:"foreignKey:UserID"`
    Profile  Profile `gorm:"foreignKey:UserID"`
}

type Post struct {
    ID      uint   `gorm:"primaryKey"`
    Title   string `gorm:"not null"`
    Content string `gorm:"type:text"`
    UserID  uint   `gorm:"not null"`
    User    User   `gorm:"foreignKey:UserID"`
}

// クエリ最適化
func (r *UserRepository) GetUsersWithPosts() ([]User, error) {
    var users []User
    return users, r.db.Preload("Posts").Find(&users).Error
}
```

### 2. マイグレーション管理
```go
// migrations/001_create_users_table.go
func CreateUsersTable(db *gorm.DB) error {
    return db.AutoMigrate(&User{})
}

// マイグレーション実行
func RunMigrations(db *gorm.DB) error {
    migrations := []func(*gorm.DB) error{
        CreateUsersTable,
        CreatePostsTable,
        // 順次追加
    }

    for _, migration := range migrations {
        if err := migration(db); err != nil {
            return err
        }
    }
    return nil
}
```

## 認証・セキュリティ

### 1. JWT認証実装
```go
type JWTService struct {
    secretKey string
}

func (j *JWTService) GenerateToken(userID uint) (string, error) {
    claims := &JWTClaims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.secretKey))
}
```

### 2. ミドルウェア実装
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r)
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        claims, err := validateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "userID", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## エラーハンドリング・ログ

### 1. 構造化エラー
```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e APIError) Error() string {
    return e.Message
}

// カスタムエラー定義
var (
    ErrUserNotFound = APIError{
        Code:    404,
        Message: "User not found",
    }
    ErrInvalidCredentials = APIError{
        Code:    401,
        Message: "Invalid credentials",
    }
)
```

### 2. 構造化ログ
```go
import "github.com/sirupsen/logrus"

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    logger := logrus.WithFields(logrus.Fields{
        "handler": "CreateUser",
        "method":  r.Method,
        "path":    r.URL.Path,
    })

    logger.Info("Creating new user")

    // 実装...

    if err != nil {
        logger.WithError(err).Error("Failed to create user")
        // エラーレスポンス
        return
    }

    logger.WithField("userID", user.ID).Info("User created successfully")
}
```

## テスト戦略

### 1. ユニットテスト
```go
func TestUserService_CreateUser(t *testing.T) {
    // Given
    mockRepo := &MockUserRepository{}
    service := NewUserService(mockRepo)

    req := CreateUserRequest{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    }

    mockRepo.On("Create", mock.AnythingOfType("*User")).Return(nil)

    // When
    user, err := service.CreateUser(req)

    // Then
    assert.NoError(t, err)
    assert.Equal(t, req.Email, user.Email)
    assert.Equal(t, req.Name, user.Name)
    mockRepo.AssertExpectations(t)
}
```

### 2. 統合テスト
```go
func TestUserAPI_Integration(t *testing.T) {
    // テスト用DB作成
    db := setupTestDB()
    defer cleanupTestDB(db)

    // テストサーバー起動
    server := setupTestServer(db)
    defer server.Close()

    // APIテスト実行
    resp, err := http.Post(server.URL+"/users", "application/json",
        strings.NewReader(`{"email":"test@example.com","name":"Test","password":"pass123"}`))

    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

## パフォーマンス最適化

### 1. データベース最適化
```go
// インデックス設定
type User struct {
    Email    string `gorm:"uniqueIndex:idx_user_email"`
    Status   string `gorm:"index:idx_user_status"`
    CreateAt time.Time `gorm:"index:idx_user_created"`
}

// バッチ処理
func (r *UserRepository) CreateUsersInBatch(users []User) error {
    return r.db.CreateInBatches(users, 100).Error
}
```

### 2. 並行処理
```go
func (s *UserService) ProcessUsersAsync(userIDs []uint) {
    semaphore := make(chan struct{}, 10) // 同時実行数制限
    var wg sync.WaitGroup

    for _, id := range userIDs {
        wg.Add(1)
        go func(userID uint) {
            defer wg.Done()
            semaphore <- struct{}{}        // セマフォ取得
            defer func() { <-semaphore }() // セマフォ解放

            s.processUser(userID)
        }(id)
    }

    wg.Wait()
}
```

## 他エージェントとの連携

### アーキテクトから
- [ ] API仕様書・データモデル設計受領
- [ ] 技術選定・アーキテクチャ指針確認
- [ ] 非機能要件確認

### フロントエンドから
- [ ] 型定義共有（OpenAPI仕様生成）
- [ ] エラーレスポンス形式統一
- [ ] 認証フロー調整

### DevOpsから
- [ ] 環境変数・設定管理
- [ ] ヘルスチェック・メトリクス実装
- [ ] ログ出力形式調整

### QAから
- [ ] テストデータ作成・シード実装
- [ ] APIテスト協力
- [ ] 負荷テスト対応

## 成果物チェックリスト
- [ ] Clean Architecture準拠
- [ ] RESTful API設計
- [ ] 認証・認可実装
- [ ] データベース最適化
- [ ] エラーハンドリング実装
- [ ] 構造化ログ実装
- [ ] 単体・統合テスト実装
- [ ] OpenAPI仕様書生成
- [ ] パフォーマンステスト実施
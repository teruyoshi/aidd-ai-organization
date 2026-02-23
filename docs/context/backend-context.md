# バックエンド技術コンテキスト

## 技術スタック詳細

### 核心技術
- **Go 1.21+**: Generics・Workspace・Fuzzing活用
- **go-chi/chi v5**: 軽量・高性能・ミドルウェア豊富
- **GORM v2**: 型安全・高性能ORM・Auto Migration
- **MySQL 8.0**: JSON型・CTE・Window Functions活用

### アーキテクチャパターン
```go
// Clean Architecture実装
// Domain Layer (最内層)
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex;not null" json:"email"`
    Password  string    `gorm:"not null" json:"-"`
    Name      string    `gorm:"not null" json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Repository Interface (Domain Layer)
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    GetByEmail(email string) (*User, error)
    List(params ListParams) ([]*User, error)
    Update(user *User) error
    Delete(id uint) error
}

// Service Layer (Application Layer)
type UserService struct {
    repo     UserRepository
    hasher   PasswordHasher
    validator Validator
}

func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
    // バリデーション
    if err := s.validator.Validate(req); err != nil {
        return nil, NewValidationError(err)
    }

    // 重複チェック
    if existing, _ := s.repo.GetByEmail(req.Email); existing != nil {
        return nil, NewConflictError("email already exists")
    }

    // パスワードハッシュ化
    hashedPassword, err := s.hasher.Hash(req.Password)
    if err != nil {
        return nil, NewInternalError("failed to hash password")
    }

    user := &User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }

    if err := s.repo.Create(user); err != nil {
        return nil, NewInternalError("failed to create user")
    }

    return user, nil
}
```

## プロジェクト構成

### 1. ディレクトリ構造
```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # アプリケーションエントリーポイント
├── internal/
│   ├── domain/                     # ドメイン層
│   │   ├── entity/                 # エンティティ
│   │   ├── repository/             # リポジトリインターフェース
│   │   └── service/                # ドメインサービス
│   ├── usecase/                    # ユースケース層
│   │   ├── interactor/             # インタラクター
│   │   └── presenter/              # プレゼンター
│   ├── interface/                  # インターフェース層
│   │   ├── handler/                # HTTPハンドラー
│   │   ├── middleware/             # ミドルウェア
│   │   └── router/                 # ルーティング
│   └── infrastructure/             # インフラ層
│       ├── database/               # DB実装
│       ├── repository/             # リポジトリ実装
│       └── external/               # 外部API
├── pkg/
│   ├── auth/                       # 認証パッケージ
│   ├── config/                     # 設定管理
│   ├── logger/                     # ログ
│   ├── validator/                  # バリデーション
│   └── errors/                     # エラー処理
├── migrations/                     # DBマイグレーション
├── scripts/                        # ビルド・デプロイスクリプト
└── docs/                          # APIドキュメント
```

### 2. 依存関係注入
```go
// internal/di/container.go
type Container struct {
    Config     *config.Config
    DB         *gorm.DB
    Logger     *logrus.Logger
    Validator  *validator.Validate

    // Repositories
    UserRepo   repository.UserRepository
    PostRepo   repository.PostRepository

    // Services
    UserService service.UserService
    PostService service.PostService

    // Handlers
    UserHandler handler.UserHandler
    PostHandler handler.PostHandler
}

func NewContainer() (*Container, error) {
    c := &Container{}

    // 設定読み込み
    config, err := config.Load()
    if err != nil {
        return nil, err
    }
    c.Config = config

    // DB接続
    db, err := database.Connect(config.Database)
    if err != nil {
        return nil, err
    }
    c.DB = db

    // ロガー初期化
    c.Logger = logger.New(config.Log)

    // バリデーター初期化
    c.Validator = validator.New()

    // リポジトリ
    c.UserRepo = repository.NewUserRepository(db)
    c.PostRepo = repository.NewPostRepository(db)

    // サービス
    c.UserService = service.NewUserService(c.UserRepo, c.Validator)
    c.PostService = service.NewPostService(c.PostRepo, c.UserRepo)

    // ハンドラー
    c.UserHandler = handler.NewUserHandler(c.UserService, c.Logger)
    c.PostHandler = handler.NewPostHandler(c.PostService, c.Logger)

    return c, nil
}
```

## データベース設計・GORM活用

### 1. モデル定義・関連付け
```go
// internal/domain/entity/user.go
type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex:idx_users_email;size:255;not null" json:"email"`
    Password  string    `gorm:"size:255;not null" json:"-"`
    Name      string    `gorm:"size:100;not null" json:"name"`
    Status    UserStatus `gorm:"type:enum('active','inactive','suspended');default:'active'" json:"status"`
    Profile   *Profile  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile,omitempty"`
    Posts     []Post    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"posts,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Profile struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    UserID   uint   `gorm:"not null;index" json:"user_id"`
    Bio      string `gorm:"type:text" json:"bio"`
    Avatar   string `gorm:"size:500" json:"avatar"`
    Website  string `gorm:"size:500" json:"website"`
    Location string `gorm:"size:100" json:"location"`
}

type Post struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null;index" json:"user_id"`
    Title     string    `gorm:"size:255;not null" json:"title"`
    Content   string    `gorm:"type:longtext;not null" json:"content"`
    Status    PostStatus `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"`
    Tags      []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"`
    User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Enum定義
type UserStatus string
const (
    UserStatusActive    UserStatus = "active"
    UserStatusInactive  UserStatus = "inactive"
    UserStatusSuspended UserStatus = "suspended"
)
```

### 2. リポジトリ実装
```go
// internal/infrastructure/repository/user_repository.go
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
    var user domain.User
    err := r.db.Preload("Profile").First(&user, id).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) List(params repository.ListParams) ([]*domain.User, error) {
    var users []*domain.User
    query := r.db.Model(&domain.User{})

    // 検索条件適用
    if params.Search != "" {
        query = query.Where("name ILIKE ? OR email ILIKE ?",
            "%"+params.Search+"%", "%"+params.Search+"%")
    }

    if params.Status != "" {
        query = query.Where("status = ?", params.Status)
    }

    // ソート
    if params.Sort != "" {
        order := params.Sort
        if params.Order == "desc" {
            order += " DESC"
        }
        query = query.Order(order)
    } else {
        query = query.Order("created_at DESC")
    }

    // ページネーション
    offset := (params.Page - 1) * params.Limit
    err := query.Offset(offset).Limit(params.Limit).
        Preload("Profile").
        Find(&users).Error

    return users, err
}

func (r *userRepository) Update(user *domain.User) error {
    return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
    return r.db.Delete(&domain.User{}, id).Error
}

// 複雑なクエリ例
func (r *userRepository) GetUsersWithPostCount() ([]*domain.UserWithStats, error) {
    var results []*domain.UserWithStats

    err := r.db.Model(&domain.User{}).
        Select("users.*, COUNT(posts.id) as post_count").
        Joins("LEFT JOIN posts ON posts.user_id = users.id AND posts.deleted_at IS NULL").
        Group("users.id").
        Having("post_count > ?", 0).
        Order("post_count DESC").
        Find(&results).Error

    return results, err
}
```

### 3. マイグレーション管理
```go
// migrations/migrate.go
package migrations

import (
    "gorm.io/gorm"
    "your-app/internal/domain/entity"
)

type Migration struct {
    ID          string
    Description string
    Up          func(*gorm.DB) error
    Down        func(*gorm.DB) error
}

var migrations = []Migration{
    {
        ID:          "001_create_users_table",
        Description: "Create users table",
        Up: func(db *gorm.DB) error {
            return db.AutoMigrate(&entity.User{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable(&entity.User{})
        },
    },
    {
        ID:          "002_create_profiles_table",
        Description: "Create profiles table",
        Up: func(db *gorm.DB) error {
            return db.AutoMigrate(&entity.Profile{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable(&entity.Profile{})
        },
    },
}

func RunMigrations(db *gorm.DB) error {
    // マイグレーション履歴テーブル作成
    if err := db.AutoMigrate(&MigrationHistory{}); err != nil {
        return err
    }

    for _, migration := range migrations {
        var count int64
        db.Model(&MigrationHistory{}).Where("migration_id = ?", migration.ID).Count(&count)

        if count == 0 {
            if err := migration.Up(db); err != nil {
                return fmt.Errorf("failed to run migration %s: %w", migration.ID, err)
            }

            // 履歴記録
            history := &MigrationHistory{
                MigrationID: migration.ID,
                Description: migration.Description,
                AppliedAt:   time.Now(),
            }
            db.Create(history)
        }
    }

    return nil
}

type MigrationHistory struct {
    ID          uint      `gorm:"primaryKey"`
    MigrationID string    `gorm:"size:255;not null;uniqueIndex"`
    Description string    `gorm:"size:500"`
    AppliedAt   time.Time `gorm:"not null"`
}
```

## HTTP API設計

### 1. ルーター設定
```go
// internal/interface/router/router.go
func NewRouter(container *di.Container) chi.Router {
    r := chi.NewRouter()

    // ミドルウェア設定
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(middleware.Timeout(60 * time.Second))
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // ヘルスチェック
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })

    // API v1ルート
    r.Route("/api/v1", func(r chi.Router) {
        // 認証不要エンドポイント
        r.Post("/auth/login", container.AuthHandler.Login)
        r.Post("/auth/register", container.AuthHandler.Register)

        // 認証必要エンドポイント
        r.Group(func(r chi.Router) {
            r.Use(container.AuthMiddleware.Authenticate)

            // ユーザー管理
            r.Route("/users", func(r chi.Router) {
                r.Get("/", container.UserHandler.List)
                r.Post("/", container.UserHandler.Create)
                r.Get("/{id}", container.UserHandler.GetByID)
                r.Put("/{id}", container.UserHandler.Update)
                r.Delete("/{id}", container.UserHandler.Delete)
            })

            // 投稿管理
            r.Route("/posts", func(r chi.Router) {
                r.Get("/", container.PostHandler.List)
                r.Post("/", container.PostHandler.Create)
                r.Get("/{id}", container.PostHandler.GetByID)
                r.Put("/{id}", container.PostHandler.Update)
                r.Delete("/{id}", container.PostHandler.Delete)
            })
        })
    })

    return r
}
```

### 2. HTTPハンドラー
```go
// internal/interface/handler/user_handler.go
type UserHandler struct {
    service service.UserService
    logger  *logrus.Logger
}

func NewUserHandler(service service.UserService, logger *logrus.Logger) *UserHandler {
    return &UserHandler{
        service: service,
        logger:  logger,
    }
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    reqID := middleware.GetReqID(ctx)

    logger := h.logger.WithFields(logrus.Fields{
        "handler":    "UserHandler.Create",
        "request_id": reqID,
    })

    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        logger.WithError(err).Error("Failed to decode request body")
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := h.service.CreateUser(ctx, req)
    if err != nil {
        logger.WithError(err).Error("Failed to create user")
        h.handleError(w, err)
        return
    }

    logger.WithField("user_id", user.ID).Info("User created successfully")

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // クエリパラメータ解析
    params := repository.ListParams{
        Page:   parseIntParam(r.URL.Query().Get("page"), 1),
        Limit:  parseIntParam(r.URL.Query().Get("limit"), 20),
        Search: r.URL.Query().Get("search"),
        Status: r.URL.Query().Get("status"),
        Sort:   r.URL.Query().Get("sort"),
        Order:  r.URL.Query().Get("order"),
    }

    users, total, err := h.service.ListUsers(ctx, params)
    if err != nil {
        h.logger.WithError(err).Error("Failed to list users")
        h.handleError(w, err)
        return
    }

    response := ListResponse{
        Data: users,
        Meta: PaginationMeta{
            Page:       params.Page,
            Limit:      params.Limit,
            Total:      total,
            TotalPages: (total + params.Limit - 1) / params.Limit,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
    var apiErr *errors.APIError
    if errors.As(err, &apiErr) {
        http.Error(w, apiErr.Message, apiErr.Code)
        return
    }

    // デフォルトエラー
    http.Error(w, "Internal server error", http.StatusInternalServerError)
}
```

## 認証・認可システム

### 1. JWT認証
```go
// pkg/auth/jwt.go
type JWTService struct {
    secretKey []byte
    issuer    string
    expiry    time.Duration
}

func NewJWTService(secretKey, issuer string, expiry time.Duration) *JWTService {
    return &JWTService{
        secretKey: []byte(secretKey),
        issuer:    issuer,
        expiry:    expiry,
    }
}

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Roles  []string `json:"roles"`
    jwt.RegisteredClaims
}

func (j *JWTService) GenerateToken(user *domain.User) (string, error) {
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Roles:  user.GetRoles(),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiry)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    j.issuer,
            Subject:   fmt.Sprintf("%d", user.ID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return j.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}
```

### 2. 認証ミドルウェア
```go
// internal/interface/middleware/auth.go
type AuthMiddleware struct {
    jwtService *auth.JWTService
    userRepo   repository.UserRepository
    logger     *logrus.Logger
}

func NewAuthMiddleware(jwtService *auth.JWTService, userRepo repository.UserRepository, logger *logrus.Logger) *AuthMiddleware {
    return &AuthMiddleware{
        jwtService: jwtService,
        userRepo:   userRepo,
        logger:     logger,
    }
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := m.extractToken(r)
        if token == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        claims, err := m.jwtService.ValidateToken(token)
        if err != nil {
            m.logger.WithError(err).Error("Invalid token")
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // ユーザー存在確認
        user, err := m.userRepo.GetByID(claims.UserID)
        if err != nil {
            m.logger.WithError(err).Error("User not found")
            http.Error(w, "User not found", http.StatusUnauthorized)
            return
        }

        // コンテキストにユーザー情報追加
        ctx := context.WithValue(r.Context(), "user", user)
        ctx = context.WithValue(ctx, "user_id", user.ID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (m *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, ok := r.Context().Value("user").(*domain.User)
            if !ok {
                http.Error(w, "User not found in context", http.StatusUnauthorized)
                return
            }

            userRoles := user.GetRoles()
            hasRole := false
            for _, role := range roles {
                for _, userRole := range userRoles {
                    if role == userRole {
                        hasRole = true
                        break
                    }
                }
                if hasRole {
                    break
                }
            }

            if !hasRole {
                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

func (m *AuthMiddleware) extractToken(r *http.Request) string {
    bearerToken := r.Header.Get("Authorization")
    if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
        return bearerToken[7:]
    }
    return ""
}
```

## エラーハンドリング

### 1. カスタムエラー定義
```go
// pkg/errors/errors.go
type ErrorCode string

const (
    ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
    ErrCodeNotFound      ErrorCode = "NOT_FOUND"
    ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
    ErrCodeForbidden     ErrorCode = "FORBIDDEN"
    ErrCodeConflict      ErrorCode = "CONFLICT"
    ErrCodeInternal      ErrorCode = "INTERNAL_ERROR"
)

type APIError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Details string    `json:"details,omitempty"`
    HTTPStatus int    `json:"-"`
}

func (e *APIError) Error() string {
    return e.Message
}

func NewValidationError(details string) *APIError {
    return &APIError{
        Code:       ErrCodeValidation,
        Message:    "Validation failed",
        Details:    details,
        HTTPStatus: http.StatusBadRequest,
    }
}

func NewNotFoundError(resource string) *APIError {
    return &APIError{
        Code:       ErrCodeNotFound,
        Message:    fmt.Sprintf("%s not found", resource),
        HTTPStatus: http.StatusNotFound,
    }
}

func NewConflictError(message string) *APIError {
    return &APIError{
        Code:       ErrCodeConflict,
        Message:    message,
        HTTPStatus: http.StatusConflict,
    }
}

func NewInternalError(details string) *APIError {
    return &APIError{
        Code:       ErrCodeInternal,
        Message:    "Internal server error",
        Details:    details,
        HTTPStatus: http.StatusInternalServerError,
    }
}
```

### 2. バリデーション
```go
// pkg/validator/validator.go
type Validator struct {
    validate *validator.Validate
}

func New() *Validator {
    v := validator.New()

    // カスタムバリデーション追加
    v.RegisterValidation("password", validatePassword)
    v.RegisterValidation("username", validateUsername)

    return &Validator{validate: v}
}

func (v *Validator) Validate(i interface{}) error {
    if err := v.validate.Struct(i); err != nil {
        var validationErrors []string

        for _, err := range err.(validator.ValidationErrors) {
            validationErrors = append(validationErrors, fmt.Sprintf(
                "Field '%s' failed validation: %s",
                err.Field(),
                getValidationMessage(err.Tag(), err.Param()),
            ))
        }

        return errors.NewValidationError(strings.Join(validationErrors, "; "))
    }

    return nil
}

func validatePassword(fl validator.FieldLevel) bool {
    password := fl.Field().String()

    // 最小8文字、大文字・小文字・数字を含む
    if len(password) < 8 {
        return false
    }

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`\d`).MatchString(password)

    return hasUpper && hasLower && hasNumber
}

// リクエスト構造体例
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,password"`
}
```

## ログ・監視

### 1. 構造化ログ
```go
// pkg/logger/logger.go
func New(config LogConfig) *logrus.Logger {
    logger := logrus.New()

    // ログレベル設定
    level, err := logrus.ParseLevel(config.Level)
    if err != nil {
        level = logrus.InfoLevel
    }
    logger.SetLevel(level)

    // フォーマット設定
    if config.Format == "json" {
        logger.SetFormatter(&logrus.JSONFormatter{
            TimestampFormat: time.RFC3339,
        })
    } else {
        logger.SetFormatter(&logrus.TextFormatter{
            FullTimestamp:   true,
            TimestampFormat: time.RFC3339,
        })
    }

    // フック追加
    if config.EnableSentry {
        hook, err := logrus_sentry.NewSentryHook(config.SentryDSN, []logrus.Level{
            logrus.PanicLevel,
            logrus.FatalLevel,
            logrus.ErrorLevel,
        })
        if err != nil {
            logger.Error("Failed to initialize Sentry hook")
        } else {
            logger.Hooks.Add(hook)
        }
    }

    return logger
}

// 使用例
func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    logger := s.logger.WithFields(logrus.Fields{
        "service": "UserService",
        "action":  "CreateUser",
        "email":   req.Email,
    })

    logger.Info("Creating new user")

    user, err := s.repo.Create(&User{
        Name:  req.Name,
        Email: req.Email,
    })

    if err != nil {
        logger.WithError(err).Error("Failed to create user")
        return nil, err
    }

    logger.WithField("user_id", user.ID).Info("User created successfully")
    return user, nil
}
```

### 2. メトリクス収集
```go
// pkg/metrics/metrics.go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests.",
        },
        []string{"method", "path", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds.",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )

    databaseQueriesTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "database_queries_total",
            Help: "Total number of database queries.",
        },
        []string{"table", "operation"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(databaseQueriesTotal)
}

// メトリクスミドルウェア
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(ww, r)

        duration := time.Since(start).Seconds()

        httpRequestsTotal.WithLabelValues(
            r.Method,
            r.URL.Path,
            fmt.Sprintf("%d", ww.statusCode),
        ).Inc()

        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}
```

## テスト戦略

### 1. テストヘルパー
```go
// internal/testutil/database.go
func SetupTestDB() (*gorm.DB, func()) {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        panic("failed to connect test database")
    }

    // マイグレーション実行
    err = db.AutoMigrate(&domain.User{}, &domain.Post{})
    if err != nil {
        panic("failed to migrate test database")
    }

    cleanup := func() {
        sqlDB, _ := db.DB()
        sqlDB.Close()
    }

    return db, cleanup
}

func CreateTestUser(db *gorm.DB) *domain.User {
    user := &domain.User{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
    db.Create(user)
    return user
}
```

### 2. モック生成
```go
//go:generate mockgen -source=repository.go -destination=mocks/mock_repository.go

// リポジトリインターフェース
type UserRepository interface {
    Create(user *User) error
    GetByID(id uint) (*User, error)
    GetByEmail(email string) (*User, error)
}
```
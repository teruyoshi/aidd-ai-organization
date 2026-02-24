package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sample-todo-backend/internal/infrastructure/database"
	"sample-todo-backend/internal/infrastructure/repository"
	"sample-todo-backend/internal/interface/router"
	"sample-todo-backend/internal/usecase/service"
	"sample-todo-backend/pkg/auth"
	"sample-todo-backend/pkg/config"
	"sample-todo-backend/pkg/logger"
	"sample-todo-backend/pkg/validator"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// ログ初期化
	log := logger.New()

	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// データベース接続
	db, err := database.NewMySQLConnection(cfg.Database)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to database")
	}

	// Redis接続
	redis, err := database.NewRedisConnection(cfg.Redis)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to Redis")
	}

	// バリデーター初期化
	validator := validator.New()

	// JWT サービス初期化
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpiryHours)

	// リポジトリ層初期化
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// サービス層初期化
	userService := service.NewUserService(userRepo, validator, jwtService)
	todoService := service.NewTodoService(todoRepo, categoryRepo, userRepo, validator)
	categoryService := service.NewCategoryService(categoryRepo, validator)

	// ルーター初期化
	r := router.NewRouter(&router.Dependencies{
		UserService:     userService,
		TodoService:     todoService,
		CategoryService: categoryService,
		JWTService:      jwtService,
		Redis:           redis,
		Logger:          log,
		Config:          cfg,
	})

	// HTTPサーバー設定
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeoutSeconds) * time.Second,
	}

	// サーバー起動（別ゴルーチン）
	go func() {
		log.WithFields(logrus.Fields{
			"port":        cfg.Server.Port,
			"environment": cfg.Environment,
			"version":     getVersion(),
		}).Info("Starting HTTP server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// グレースフルシャットダウンの設定
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// シャットダウンシグナル待機
	<-quit
	log.Info("Shutting down server...")

	// シャットダウンタイムアウト設定
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// HTTPサーバーのグレースフルシャットダウン
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Server forced to shutdown")
	}

	// データベース接続のクリーンアップ
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
		log.Info("Database connection closed")
	}

	// Redis接続のクリーンアップ
	if err := redis.Close(); err != nil {
		log.WithError(err).Error("Failed to close Redis connection")
	} else {
		log.Info("Redis connection closed")
	}

	log.Info("Server exited")
}

// getVersion returns the application version
// 本来はビルド時にldFlagsで設定する
func getVersion() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}
	return "development"
}

// ヘルスチェック用の簡易エンドポイント
func healthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"sample-todo-backend"}`))
	}
}
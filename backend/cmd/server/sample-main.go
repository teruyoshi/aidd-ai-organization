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
	"syscall"
	"time"

	"log"
)

func main() {
	// 設定読み込み
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// データベース接続
	db, err := database.NewMySQLConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// JWT サービス初期化
	jwtService := auth.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpiryHours)

	// リポジトリ層初期化
	userRepo := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepository(db)

	// サービス層初期化
	userService := service.NewUserService(userRepo, jwtService)
	todoService := service.NewTodoService(todoRepo)

	// ルーター初期化
	r := router.NewRouter(&router.Dependencies{
		UserService: userService,
		TodoService: todoService,
		JWTService:  jwtService,
		Config:      cfg,
	})

	// HTTPサーバー設定
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバー起動
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// グレースフルシャットダウン
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	log.Println("Server exited")
}

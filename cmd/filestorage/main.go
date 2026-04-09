package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vadim/filestorage/internal/config"
	"github.com/vadim/filestorage/internal/handler"
	"github.com/vadim/filestorage/internal/service"
	"github.com/vadim/filestorage/internal/storage"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DBDSN)
	if err != nil {
		log.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Error("postgres ping failed", "err", err)
		os.Exit(1)
	}

	userStore := storage.NewUserStorage(pool)
	fileStore := storage.NewFileStorage(pool)

	authSvc := service.NewAuthService(userStore, cfg.JWTSecret)
	fileSvc := service.NewFileService(fileStore, cfg.UploadDir, cfg.MaxQuotaMB, log)

	authHandler := handler.NewAuthHandler(authSvc, log)
	fileHandler := handler.NewFileHandler(fileSvc, log)

	router := handler.NewRouter(authHandler, fileHandler, cfg.JWTSecret, log)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown error", "err", err)
	}
}

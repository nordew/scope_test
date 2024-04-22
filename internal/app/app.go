package app

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nordew/scope_test/internal/config"
	v1 "github.com/nordew/scope_test/internal/controller/http/v1"
	"github.com/nordew/scope_test/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func MustRun() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load config")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("failed to get cfg:", err.Error())
	}

	workerPool := service.NewWorkerPool(cfg.MaxWorkers, cfg.MaxQueueSize)

	handler := v1.NewHandler(workerPool)

	server := &http.Server{
		Addr:         fmt.Sprintf(cfg.ServerAddress),
		Handler:      handler.Init(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	workerPool.Shutdown()

	log.Println("server  stopped")
}

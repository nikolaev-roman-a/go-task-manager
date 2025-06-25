package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nikolaev-roman-a/go-task-manager/internal/repository"
	"github.com/nikolaev-roman-a/go-task-manager/internal/server"
	"github.com/nikolaev-roman-a/go-task-manager/internal/services"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	taskStore := repository.NewTaskStore()
	taskService := services.NewTaskService(taskStore, logger)
	httpServer := server.NewHTTPServer(taskService, logger)
	httpServer.Run()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	logger.Info("stopped")
}

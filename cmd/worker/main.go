package main

import (
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/messagebroker"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/authservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.GetConfig()
	log := logger.NewLogger(cfg.Auth.Name, cfg.Log)

	queue := messagebroker.NewQueue(log, cfg)
	if err := queue.SetupRabbitMQ(cfg.RabbitMQ.URL, log); err != nil {
		log.Fatal(logger.Queue, logger.Startup, fmt.Sprintf("Failed to setup queue, error: %v", err), nil)
	}

	log.Info(logger.Queue, logger.Startup, fmt.Sprintf("Setup queue successfully"), nil)

	messagebroker.RegisterAllQueues(
		authservice.SendEmailOTPEvent(queue),
		authservice.SendWelcomeEvent(queue),
		// add new queues here
		// ...
	)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

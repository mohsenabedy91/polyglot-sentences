//go:build !test

package main

import (
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/cmd/setup"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/client"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/messagebroker"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/authservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configProvider := &config.Config{}
	conf := configProvider.GetConfig()
	log := logger.NewLogger(conf.Auth.Name, conf.Log)

	queue, err := setup.InitializeQueue(log, conf)
	if err != nil {
		return
	}
	defer queue.Driver.Close()

	log.Info(logger.Queue, logger.Startup, fmt.Sprintf("Setup queue successfully"), nil)

	userClient := client.NewUserClient(log, conf.UserManagement)
	defer userClient.Close()

	messagebroker.RegisterAllQueues(
		authservice.SendEmailOTPEvent(queue),
		authservice.SendWelcomeEvent(queue, userClient),
		authservice.SendResetPasswordLinkEvent(queue),
		// add new queues here
		// ...
	)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh

	log.Info(logger.Internal, logger.Shutdown, "Shutdown Server ...", nil)
}

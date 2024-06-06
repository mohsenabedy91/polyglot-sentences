package authservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/email"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/messagebroker"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"html/template"
	"strings"
)

type SendWelcome struct {
	queue       *messagebroker.Queue
	emailSender port.EmailSender
	userClient  port.UserClient
}

var sendWelcomeInstance *SendWelcome

const delaySendWelcomeSeconds int64 = 60

type SendWelcomeDto struct {
	UserID   uint64 `json:"userID"`
	To       string `json:"to"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

func SendWelcomeEvent(queue *messagebroker.Queue, userClient port.UserClient) *SendWelcome {
	if sendWelcomeInstance == nil {
		sendWelcomeInstance = &SendWelcome{
			queue:       queue,
			emailSender: email.NewSender(queue.Log, queue.Config.SendGrid),
			userClient:  userClient,
		}
	}

	return sendWelcomeInstance
}

func (r *SendWelcome) Name() string {
	return "send_welcome"
}

func (r *SendWelcome) Publish(msg interface{}) {

	if err := r.queue.Driver.Produce(r.Name(), msg, delaySendWelcomeSeconds); err != nil {
		return
	}
	r.queue.Log.Info(logger.Queue, logger.RabbitMQPublish, fmt.Sprintf("published successfully to queue: %s", msg), nil)
}

func (r *SendWelcome) Consume(message []byte) error {
	extra := map[logger.ExtraKey]interface{}{
		logger.Body: string(message),
	}
	var msg SendWelcomeDto
	if err := json.Unmarshal(message, &msg); err != nil {
		r.queue.Log.Error(logger.Queue, logger.RabbitMQConsume, fmt.Sprintf("Error unmarshalling message, error: %v", err), extra)
		return err
	}

	trans := translation.NewTranslation(msg.Language)
	appName := trans.Lang("appName", nil, &msg.Language)

	if strings.TrimSpace(msg.Name) == "" {
		msg.Name = trans.Lang("user", nil, &msg.Language)
	}

	emailBuffer := new(bytes.Buffer)
	parseFiles, err := template.ParseFiles("internal/core/views/email/base.html", "internal/core/views/email/auth/welcome.html")
	if err != nil {
		r.queue.Log.Error(logger.Email, logger.SendEmail, err.Error(), nil)
		return err
	}

	body := template.HTML(trans.Lang("email.welcome.body", map[string]interface{}{
		"username":     msg.Name,
		"supportEmail": r.queue.Config.App.SupportEmail,
		"app":          appName,
	}, &msg.Language))

	data := map[string]interface{}{
		"language": msg.Language,
		"body":     body,
	}

	if err = parseFiles.ExecuteTemplate(emailBuffer, "base.html", data); err != nil {
		r.queue.Log.Error(logger.Email, logger.SendEmail, err.Error(), nil)
		return err
	}

	subject := trans.Lang("email.welcome.subject", map[string]interface{}{
		"app": appName,
	}, &msg.Language)

	err = r.emailSender.Send(msg.To, msg.Name, subject, string(body))
	if err == nil {
		if updateErr := r.userClient.MarkWelcomeMessageSent(context.Background(), msg.UserID); updateErr != nil {
			r.queue.Log.Error(logger.Email, logger.SendEmail, updateErr.Error(), nil)
		}
	}

	return err
}

func (r *SendWelcome) RegisterQueue() {
	go func() {
		if err := r.queue.Driver.RegisterConsumer(r.Name(), r.Consume); err != nil {
			r.queue.Log.Error(
				logger.Queue,
				logger.RabbitMQRegisterConsumer,
				fmt.Sprintf("Error on registering consumer, error: %v", err),
				map[logger.ExtraKey]interface{}{
					logger.QueueName: r.Name(),
				},
			)
		}
	}()
}

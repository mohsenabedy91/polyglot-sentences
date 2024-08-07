package authevent

import (
	"bytes"
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

type SendResetPasswordLink struct {
	queue       *messagebroker.Queue
	emailSender port.EmailSender
}

var resetPasswordLinkInstance *SendResetPasswordLink

const DelaySendResetPasswordLinkSeconds int64 = 0
const SendResetPasswordLinkName = "send_reset_password_link"

type SendResetPasswordLinkDto struct {
	To       string `json:"to"`
	Name     string `json:"name"`
	OTP      string `json:"otp"`
	Language string `json:"language"`
}

func NewSendResetPasswordLink(queue *messagebroker.Queue) *SendResetPasswordLink {
	if resetPasswordLinkInstance == nil {
		resetPasswordLinkInstance = &SendResetPasswordLink{
			queue:       queue,
			emailSender: email.NewSender(queue.Log, queue.Config.SendGrid),
		}
	}

	return resetPasswordLinkInstance
}

func (r *SendResetPasswordLink) Name() string {
	return SendResetPasswordLinkName
}

func (r *SendResetPasswordLink) Publish(message interface{}) {

	if err := r.queue.Driver.Produce(r.Name(), message, DelaySendResetPasswordLinkSeconds); err != nil {
		return
	}
	r.queue.Log.Info(logger.Queue, logger.RabbitMQPublish, fmt.Sprintf("published successfully to queue: %s", message), nil)
}

func (r *SendResetPasswordLink) Consume(message []byte) error {
	extra := map[logger.ExtraKey]interface{}{
		logger.Body: string(message),
	}
	var msg SendResetPasswordLinkDto
	if err := json.Unmarshal(message, &msg); err != nil {
		r.queue.Log.Error(logger.Queue, logger.RabbitMQConsume, fmt.Sprintf("Error unmarshalling message, error: %v", err), extra)
		return err
	}

	trans := translation.NewTranslation(r.queue.Config.App)
	appName := trans.Lang("appName", nil, &msg.Language)

	if strings.TrimSpace(msg.Name) == "" {
		msg.Name = trans.Lang("user", nil, &msg.Language)
	}

	emailBuffer := new(bytes.Buffer)
	parseFiles, err := template.ParseFiles("internal/core/views/email/base.html", "internal/core/views/email/auth/reset_password.html")
	if err != nil {
		r.queue.Log.Error(logger.Email, logger.SendEmail, err.Error(), nil)
		return err
	}

	body := template.HTML(trans.Lang("email.resetPassword.body", map[string]interface{}{
		"username":         msg.Name,
		"app":              appName,
		"resetPasswordUrl": r.queue.Config.App.ResetPasswordURL + msg.OTP,
		"supportEmail":     r.queue.Config.App.SupportEmail,
	}, &msg.Language))

	data := map[string]interface{}{
		"language": msg.Language,
		"body":     body,
	}

	if err = parseFiles.ExecuteTemplate(emailBuffer, "base.html", data); err != nil {
		r.queue.Log.Error(logger.Email, logger.SendEmail, err.Error(), nil)
		return err
	}

	subject := trans.Lang("email.resetPassword.subject", map[string]interface{}{
		"app": appName,
	}, &msg.Language)

	err = r.emailSender.Send(msg.To, msg.Name, subject, emailBuffer.String())

	return err
}

func (r *SendResetPasswordLink) Register() {
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

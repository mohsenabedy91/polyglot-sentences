package authservice

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
)

type SendEmailOTP struct {
	queue       *messagebroker.Queue
	emailSender port.EmailSender
}

var sendEmailOTPInstance *SendEmailOTP

const delaySendEmailOTPSeconds int64 = 0

type SendEmailOTPDto struct {
	To       string `json:"to"`
	Name     string `json:"name"`
	OTP      string `json:"otp"`
	Language string `json:"language"`
}

func SendEmailOTPEvent(queue *messagebroker.Queue) *SendEmailOTP {
	if sendEmailOTPInstance == nil {
		sendEmailOTPInstance = &SendEmailOTP{
			queue:       queue,
			emailSender: email.NewSender(queue.Log, queue.Config.SendGrid),
		}
	}

	return sendEmailOTPInstance
}

func (r *SendEmailOTP) Name() string {
	return "send_email_otp"
}

func (r *SendEmailOTP) Publish(msg interface{}) {

	if err := r.queue.Driver.Produce(r.Name(), msg, delaySendEmailOTPSeconds); err != nil {
		return
	}
	r.queue.Log.Info(logger.Queue, logger.RabbitMQPublish, fmt.Sprintf("published successfully to queue: %s", msg), nil)
}

func (r *SendEmailOTP) Consume(message []byte) error {
	extra := map[logger.ExtraKey]interface{}{
		logger.Body: string(message),
	}
	var msg SendEmailOTPDto
	if err := json.Unmarshal(message, &msg); err != nil {
		r.queue.Log.Error(logger.Queue, logger.RabbitMQConsume, fmt.Sprintf("Error unmarshalling message, error: %v", err), extra)
		return err
	}

	trans := translation.NewTranslation(msg.Language)
	appName := trans.Lang("appName", nil, &msg.Language)

	emailBuffer := new(bytes.Buffer)
	parseFiles, err := template.ParseFiles("internal/core/views/email/base.html", "internal/core/views/email/auth/verify_email.html")
	if err != nil {
		r.queue.Log.Error(logger.Email, logger.SendEmail, err.Error(), nil)
		return err
	}

	body := template.HTML(trans.Lang("email.verifyEmail.body", map[string]interface{}{
		"username":        msg.Name,
		"app":             appName,
		"otp":             msg.OTP,
		"verificationUrl": r.queue.Config.App.VerificationURL + msg.OTP,
	}, &msg.Language))

	data := map[string]interface{}{
		"language": msg.Language,
		"body":     body,
	}

	if err = parseFiles.ExecuteTemplate(emailBuffer, "base.html", data); err != nil {
		r.queue.Log.Error(logger.Email, logger.SendEmail, err.Error(), nil)
		return err
	}

	subject := trans.Lang("email.verifyEmail.subject", map[string]interface{}{
		"app": appName,
	}, &msg.Language)

	err = r.emailSender.Send(msg.To, msg.Name, subject, emailBuffer.String())

	return err
}

func (r *SendEmailOTP) Register() {
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

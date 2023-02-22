package monitor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kgantsov/uptime/app/model"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

const TelegramNotifierTimeout int = 10

type Notifier interface {
	Notify(message string)
}

type TelegramNotifier struct {
	notification *model.Notification
	client       http.Client
	logger       *logrus.Logger
}

func NewTelegramNotifier(logger *logrus.Logger, notification *model.Notification) *TelegramNotifier {
	client := http.Client{Timeout: time.Duration(TelegramNotifierTimeout) * time.Second}

	n := &TelegramNotifier{
		client:       client,
		notification: notification,
		logger:       logger,
	}

	return n
}

func (n *TelegramNotifier) Notify(message string) {
	log.Infof("Sending telegram message: %s to %s\n", message, n.notification.CallbackChatID)

	bodyParams := map[string]interface{}{
		"chat_id":              n.notification.CallbackChatID,
		"text":                 message,
		"disable_notification": true,
	}

	jsonBody, _ := json.Marshal(bodyParams)

	request, err := http.NewRequest(
		"POST", n.notification.Callback, bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Infof("Failed to notify telegram %s\n", err)
		return
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := n.client.Do(request)
	if err != nil {
		log.Infof("Failed to notify telegram %s\n", err)
		return
	}

	response.Body.Close()
}

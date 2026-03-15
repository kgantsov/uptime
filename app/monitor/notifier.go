package monitor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kgantsov/uptime/app/model"
	"github.com/rs/zerolog"
)

const TelegramNotifierTimeout int = 10

type Notifier interface {
	Notify(message string)
}

type TelegramNotifier struct {
	notification *model.Notification
	client       http.Client
	logger       zerolog.Logger
}

func NewTelegramNotifier(notification *model.Notification) *TelegramNotifier {
	client := http.Client{Timeout: time.Duration(TelegramNotifierTimeout) * time.Second}

	n := &TelegramNotifier{
		client:       client,
		notification: notification,
	}

	return n
}

func (n *TelegramNotifier) Notify(message string) {
	n.logger.Info().Msgf("Sending telegram message: %s to %s", message, n.notification.CallbackChatID)

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
		n.logger.Info().Msgf("Failed to notify telegram %s", err)
		return
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := n.client.Do(request)
	if err != nil {
		n.logger.Info().Msgf("Failed to notify telegram %s", err)
		return
	}

	response.Body.Close()
}

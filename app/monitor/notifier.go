package monitor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kgantsov/uptime/app/model"
	"github.com/rs/zerolog/log"
)

const TelegramNotifierTimeout int = 10

type Notifier interface {
	Notify(message string)
}

type TelegramNotifier struct {
	notification *model.Notification
	serviceID    uint
	client       http.Client
}

func NewTelegramNotifier(serviceID uint, notification *model.Notification) *TelegramNotifier {
	client := http.Client{Timeout: time.Duration(TelegramNotifierTimeout) * time.Second}

	n := &TelegramNotifier{
		client:       client,
		serviceID:    serviceID,
		notification: notification,
	}

	return n
}

func (n *TelegramNotifier) Notify(message string) {
	log.Info().
		Uint("service_id", n.serviceID).
		Str("chat_id", n.notification.CallbackChatID).
		Str("message", message).
		Msg("Sending telegram message")

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
		log.Error().
			Uint("service_id", n.serviceID).
			Str("chat_id", n.notification.CallbackChatID).
			Err(err).
			Msg("Failed to create telegram notify request")
		return
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := n.client.Do(request)
	if err != nil {
		log.Error().
			Uint("service_id", n.serviceID).
			Str("chat_id", n.notification.CallbackChatID).
			Err(err).
			Msg("Failed to send telegram notification")
		return
	}

	response.Body.Close()
}

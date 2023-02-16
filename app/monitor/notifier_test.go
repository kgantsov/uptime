package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kgantsov/uptime/app/model"
	"github.com/stretchr/testify/assert"
)

func TestTelegramNotifier(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/fake/telegram/messages", req.URL.String())
		assert.Equal(t, "POST", req.Method)

		var body map[string]interface{}

		err := json.NewDecoder(req.Body).Decode(&body)

		assert.Equal(t, nil, err)
		assert.Equal(t, "909091231", body["chat_id"])
		assert.Equal(t, "Test message", body["text"])

		rw.WriteHeader(200)
		rw.Write([]byte(`OK`))
	}))

	notifier := NewTelegramNotifier(
		&model.Notification{
			CallbackType:   "telegram",
			CallbackChatID: "909091231",
			Callback:       fmt.Sprintf("%s/fake/telegram/messages", server.URL),
		},
	)
	notifier.Notify("Test message")

	server.Close()
}

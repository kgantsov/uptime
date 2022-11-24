package monitor

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kgantsov/uptime/app/model"
	"github.com/kyokomi/emoji"
	"gorm.io/gorm"
)

type Monitor struct {
	DB      *gorm.DB
	client  http.Client
	done    chan struct{}
	service model.Service
}

func NewMonitor(db *gorm.DB, service model.Service) *Monitor {
	client := http.Client{Timeout: time.Duration(service.Timeout) * time.Second}

	m := &Monitor{
		DB:      db,
		service: service,
		client:  client,
		done:    make(chan struct{}),
	}

	return m
}

func (m *Monitor) NotifyTg(notification model.Notification, message string) {
	log.Infof("Sending telegram message: %s to %s\n", message, notification.CallbackChatID)

	bodyParams := map[string]interface{}{
		"chat_id":              notification.CallbackChatID,
		"text":                 message,
		"disable_notification": true,
	}

	jsonBody, _ := json.Marshal(bodyParams)

	request, err := http.NewRequest(
		"POST", notification.Callback, bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		log.Infof("Failed to notify telegram %s\n", err)
		return
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := m.client.Do(request)
	if err != nil {
		log.Infof("Failed to notify telegram %s\n", err)
		return
	}

	response.Body.Close()
}

func (m *Monitor) CheckHealth() (int, string) {
	resp, err := m.client.Get(m.service.URL)

	if err != nil {
		log.Infof("Error checking '%s' %s %s\n", m.service.Name, m.service.URL, err)
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		return 0, "TIMEOUT"
	} else if err != nil {
		return 0, "DOWN"
	}

	defer resp.Body.Close()

	if resp.StatusCode == m.service.AcceptedStatusCode {
		return resp.StatusCode, "UP"
	}

	return resp.StatusCode, "FAILED"
}

func (m *Monitor) Start() {
	log.Infof("Starting '%s' %s service monitoring\n", m.service.Name, m.service.URL)

	failing := false
	startedFailingAt := time.Time{}

	ticker := time.NewTicker(time.Duration(m.service.CheckInterval) * time.Second)

	for {
		select {
		case <-m.done:
			log.Infof("Stop monitoring for '%s' %s\n", m.service.Name, m.service.URL)
			return
		case t := <-ticker.C:
			start := time.Now()

			statusCode, status := m.CheckHealth()

			elapsed := time.Since(start)

			m.DB.Create(
				&model.Heartbeat{
					ServiceID:    m.service.ID,
					Status:       status,
					StatusCode:   statusCode,
					ResponseTime: elapsed.Milliseconds(),
				},
			)

			if status == "UP" {
				log.Infof("Service %s %s is up and running: %d %s %t\n", t, m.service.URL, statusCode, status, failing)

				if failing {
					for _, notification := range m.service.Notifications {
						m.NotifyTg(
							notification,
							emoji.Sprintf(
								":check_mark_button: Service '%s' %s is UP again. Downtime: %s",
								m.service.Name,
								m.service.URL,
								time.Since(startedFailingAt),
							),
						)
					}
					failing = false
					startedFailingAt = time.Time{}
				}
			} else {
				log.Infof("Failed to get %s url. Got status code: %d %s %t\n", m.service.URL, statusCode, status, failing)

				if !failing {
					for _, notification := range m.service.Notifications {
						m.NotifyTg(
							notification,
							emoji.Sprintf(
								":exclamation: Service '%s' %s is DOWN",
								m.service.Name,
								m.service.URL,
							),
						)
					}

					failing = true
					startedFailingAt = time.Now()
				}
			}
		}
	}
}

func (m *Monitor) Stop() {
	m.done <- struct{}{}
	log.Infof("Stopping '%s'\n", m.service.Name)
}

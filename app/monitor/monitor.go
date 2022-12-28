package monitor

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/kgantsov/uptime/app/model"
	"github.com/kyokomi/emoji"
	"gorm.io/gorm"
)

type Monitor struct {
	DB        *gorm.DB
	client    http.Client
	done      chan struct{}
	checker   Checker
	notifiers []Notifier
	service   model.Service
}

func NewMonitor(db *gorm.DB, service model.Service) *Monitor {
	log.Infof("NewMonitor %d", service.ID)

	client := http.Client{Timeout: time.Duration(service.Timeout) * time.Second}

	notifiers := []Notifier{}

	for i := range service.Notifications {
		notification := service.Notifications[i]
		notifier := NewTelegramNotifier(notification)
		notifiers = append(notifiers, notifier)
	}

	checker := NewHTTPCHecker(
		service.Name, service.URL, service.Timeout, service.AcceptedStatusCode,
	)

	m := &Monitor{
		DB:        db,
		service:   service,
		client:    client,
		done:      make(chan struct{}),
		notifiers: notifiers,
		checker:   checker,
	}

	return m
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

			statusCode, status := m.checker.Check()

			elapsed := time.Since(start)

			log.Debugf(
				"Service check %d %s %d %s %t", m.service.ID, m.service.URL, statusCode, status, failing,
			)

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
					for _, notifier := range m.notifiers {
						notifier.Notify(
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
					for _, notifier := range m.notifiers {
						notifier.Notify(
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

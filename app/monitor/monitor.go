package monitor

import (
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/sirupsen/logrus"

	"github.com/kgantsov/uptime/app/model"
	"github.com/kyokomi/emoji"
	"gorm.io/gorm"
)

type Monitor struct {
	DB        *gorm.DB
	done      chan struct{}
	checker   Checker
	notifiers []Notifier
	service   *model.Service
	logger    *logrus.Logger
}

func NewMonitor(db *gorm.DB, logger *logrus.Logger, service *model.Service) *Monitor {
	logger.Infof("NewMonitor %d", service.ID)

	notifiers := []Notifier{}

	for i := range service.Notifications {
		notification := service.Notifications[i]
		notifier := NewTelegramNotifier(logger, &notification)
		notifiers = append(notifiers, notifier)
	}

	checker := NewHTTPCHecker(
		logger, service.Name, service.URL, service.Timeout, service.AcceptedStatusCode,
	)

	m := &Monitor{
		DB:        db,
		service:   service,
		done:      make(chan struct{}),
		notifiers: notifiers,
		checker:   checker,
		logger:    logger,
	}

	return m
}

func (m *Monitor) Start() {
	m.logger.Infof("Starting '%s' %s service monitoring\n", m.service.Name, m.service.URL)

	failing := false
	startedFailingAt := time.Time{}

	ticker := time.NewTicker(time.Duration(m.service.CheckInterval) * time.Second)

	for {
		select {
		case <-m.done:
			m.logger.Infof("Stop monitoring for '%s' %s\n", m.service.Name, m.service.URL)
			return
		case t := <-ticker.C:
			start := time.Now()

			var statusCode int

			var status string

			m.logger.Debugf("Check attempt %d %s", m.service.ID, m.service.URL)

			retry.Do(
				func() error {
					statusCode, status = m.checker.Check()

					if status == StatusUp {
						return nil
					}

					if failing {
						return nil
					}

					m.logger.Debugf("Failed will retry %d %s %d %s %s", m.service.ID, m.service.URL, statusCode, status, time.Now())

					return fmt.Errorf("server is not up. Status code %d", statusCode)
				},
				retry.Delay(time.Duration(1000)*time.Millisecond),
				retry.MaxDelay(time.Duration(8)*time.Second),
				retry.Attempts(uint(m.service.Retries+1)),
			)

			elapsed := time.Since(start)

			m.logger.Debugf(
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

			if status == StatusUp {
				m.logger.Infof(
					"Service %s %s is up and running: %d %s %t\n",
					t,
					m.service.URL,
					statusCode,
					status,
					failing,
				)

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
				m.logger.Infof(
					"Failed to get %s url. Got status code: %d %s %t\n",
					m.service.URL,
					statusCode,
					status,
					failing,
				)

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
	m.logger.Infof("Stopping '%s'\n", m.service.Name)
}

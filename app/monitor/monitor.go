package monitor

import (
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/rs/zerolog/log"

	"github.com/kgantsov/uptime/app/model"
	"github.com/kgantsov/uptime/app/repository"
	"github.com/kyokomi/emoji"
)

type Monitor struct {
	heartbeatRepo repository.HeartbeatRepository
	done          chan struct{}
	checker       Checker
	notifiers     []Notifier
	service       *model.Service
}

func NewMonitor(heartbeatRepo repository.HeartbeatRepository, service *model.Service) *Monitor {
	log.Info().Uint("service_id", service.ID).Msg("NewMonitor")

	notifiers := []Notifier{}

	for i := range service.Notifications {
		notification := service.Notifications[i]
		notifier := NewTelegramNotifier(service.ID, &notification)
		notifiers = append(notifiers, notifier)
	}

	checker := NewHTTPCHecker(
		service.ID, service.Name, service.URL, service.Timeout, service.AcceptedStatusCode,
	)

	m := &Monitor{
		heartbeatRepo: heartbeatRepo,
		service:       service,
		done:          make(chan struct{}),
		notifiers:     notifiers,
		checker:       checker,
	}

	return m
}

func (m *Monitor) Start() {
	log.Info().
		Uint("service_id", m.service.ID).
		Str("name", m.service.Name).
		Str("url", m.service.URL).
		Msg("Starting service monitoring")

	failing := false
	startedFailingAt := time.Time{}

	ticker := time.NewTicker(time.Duration(m.service.CheckInterval) * time.Second)

	for {
		select {
		case <-m.done:
			log.Info().
				Uint("service_id", m.service.ID).
				Str("name", m.service.Name).
				Str("url", m.service.URL).
				Msg("Stop monitoring for service")
			return
		case t := <-ticker.C:
			var start time.Time

			var elapsed time.Duration

			var statusCode int

			var status string

			log.Debug().
				Uint("service_id", m.service.ID).
				Str("url", m.service.URL).
				Msg("Check attempt")

			retry.Do(
				func() error {
					start = time.Now()
					defer func(start time.Time) {
						elapsed = time.Since(start)
					}(start)

					statusCode, status = m.checker.Check()

					if status == StatusUp {
						return nil
					}

					if failing {
						return nil
					}

					log.Debug().
						Uint("service_id", m.service.ID).
						Str("url", m.service.URL).
						Int("status_code", statusCode).
						Str("status", status).
						Time("time", time.Now()).
						Msg("Failed, will retry")

					return fmt.Errorf("server is not up. Status code %d", statusCode)
				},
				retry.Delay(time.Duration(1000)*time.Millisecond),
				retry.MaxDelay(time.Duration(8)*time.Second),
				retry.Attempts(uint(m.service.Retries+1)),
			)

			log.Debug().
				Uint("service_id", m.service.ID).
				Str("url", m.service.URL).
				Int("status_code", statusCode).
				Str("status", status).
				Bool("failing", failing).
				Msg("Service check")

			m.heartbeatRepo.Create(
				&model.Heartbeat{
					ServiceID:    m.service.ID,
					Status:       status,
					StatusCode:   statusCode,
					ResponseTime: elapsed.Milliseconds(),
				},
			)

			if status == StatusUp {
				log.Info().
					Uint("service_id", m.service.ID).
					Str("url", m.service.URL).
					Int("status_code", statusCode).
					Str("status", status).
					Bool("failing", failing).
					Time("checked_at", t).
					Msg("Service is up and running")

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
				log.Info().
					Uint("service_id", m.service.ID).
					Str("url", m.service.URL).
					Int("status_code", statusCode).
					Str("status", status).
					Bool("failing", failing).
					Msg("Service is down")

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
	log.Info().
		Uint("service_id", m.service.ID).
		Str("name", m.service.Name).
		Msg("Stopping service monitoring")
}

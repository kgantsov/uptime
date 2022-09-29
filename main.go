package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kyokomi/emoji"
)

type Notification struct {
	CallbackType   string `json:"callback_type"`
	CallbackChatID string `json:"callback_chat_id"`
	Callback       string `json:"callback"`
}

type Service struct {
	Name          string         `json:"name"`
	URL           string         `json:"url"`
	Timeout       int            `json:"timeout"`
	CheckInterval int            `json:"check_interval"`
	Notifications []Notification `json:"notifications"`
}

type Config struct {
	Services []Service `json:"services"`
}

type Monitor struct {
	client  http.Client
	done    chan struct{}
	service Service
}

func NewMonitor(service Service) *Monitor {
	client := http.Client{Timeout: time.Duration(service.Timeout) * time.Second}

	m := &Monitor{
		service: service,
		client:  client,
		done:    make(chan struct{}),
	}

	return m
}

func (m *Monitor) NotifyTg(notification Notification, message string) {
	fmt.Printf("Sending telegram message: %s to %s\n", message, notification.CallbackChatID)

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
		fmt.Printf("Failed to notify telegram %s\n", err)
		return
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	response, err := m.client.Do(request)
	if err != nil {
		fmt.Printf("Failed to notify telegram %s\n", err)
		return
	}

	response.Body.Close()
}

func (m *Monitor) CheckHealth() int {
	resp, err := m.client.Get(m.service.URL)

	if err != nil {
		fmt.Printf("Error checking '%s' %s %s\n", m.service.Name, m.service.URL, err)
		return -1
	}

	defer resp.Body.Close()

	return resp.StatusCode
}

func (m *Monitor) Start() {
	fmt.Printf("Starting '%s' %s service monitoring\n", m.service.Name, m.service.URL)

	failing := false
	startedFailingAt := time.Time{}

	ticker := time.NewTicker(time.Duration(m.service.CheckInterval) * time.Second)

	for {
		select {
		case <-m.done:
			fmt.Printf("Stop monitoring for '%s' %s\n", m.service.Name, m.service.URL)
			return
		case t := <-ticker.C:
			status := m.CheckHealth()
			if status != http.StatusOK {
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

				fmt.Printf("Failed to get %s url. Got status code: %d\n", m.service.URL, status)
			} else {
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
				fmt.Printf("Service %s %s is up and running: %d\n", t, m.service.URL, status)
			}
		}
	}
}

func (m *Monitor) Stop() {
	m.done <- struct{}{}
	fmt.Printf("Stopping '%s'\n", m.service.Name)
}

func main() {
	configPathPtr := flag.String("config", "./config.json", "A path to a config.json file")

	flag.Parse()

	configFile, err := os.Open(*configPathPtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully Opened jsonFile")

	defer configFile.Close()

	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		fmt.Printf("Error reading %s %s \n", *configPathPtr, err)
		return
	}

	var config Config

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Printf("Error unmarshaling %s %s \n", *configPathPtr, err)
		return
	}

	fmt.Printf("Found %d services to monitor\n", len(config.Services))

	monitors := []*Monitor{}

	for _, service := range config.Services {
		m := NewMonitor(service)
		monitors = append(monitors, m)

		go m.Start()
	}

	done := make(chan struct{})
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("Got signal: %s\n", sig)

		fmt.Print("Stopping monitoring\n")

		for _, monitor := range monitors {
			monitor.Stop()
		}

		time.Sleep(100 * time.Millisecond)

		done <- struct{}{}
	}()

	emoji.Printf("Started uptime monitor\n")

	<-done

	fmt.Print("Stopped monitoring\n")
}

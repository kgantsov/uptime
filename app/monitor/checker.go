package monitor

import (
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const StatusTimeout string = "TIMEOUT"
const StatusDown string = "DOWN"
const StatusUp string = "UP"
const StatusFailed string = "FAILED"

type Checker interface {
	Check() (int, string)
}

type HTTPCHecker struct {
	client             http.Client
	ServiceID          uint
	URL                string
	Name               string
	AcceptedStatusCode int
}

func NewHTTPCHecker(serviceID uint, name, url string, timeout, acceptedStatusCode int) *HTTPCHecker {
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}

	c := &HTTPCHecker{
		client:             client,
		ServiceID:          serviceID,
		URL:                url,
		Name:               name,
		AcceptedStatusCode: acceptedStatusCode,
	}

	return c
}

func (c *HTTPCHecker) Check() (int, string) {
	resp, err := c.client.Get(c.URL)

	if err != nil {
		log.Info().
			Uint("service_id", c.ServiceID).
			Str("name", c.Name).
			Str("url", c.URL).
			Err(err).
			Msg("Error checking service")
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		return 0, StatusTimeout
	} else if err != nil {
		return 0, StatusDown
	}

	defer resp.Body.Close()

	if resp.StatusCode == c.AcceptedStatusCode {
		return resp.StatusCode, StatusUp
	}

	return resp.StatusCode, StatusFailed
}

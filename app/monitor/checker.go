package monitor

import (
	"net"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Checker interface {
	Check() (int, string)
}

type HTTPCHecker struct {
	client             http.Client
	URL                string
	Name               string
	AcceptedStatusCode int
}

func NewHTTPCHecker(name, url string, timeout, acceptedStatusCode int) *HTTPCHecker {
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}

	c := &HTTPCHecker{
		client:             client,
		URL:                url,
		AcceptedStatusCode: acceptedStatusCode,
	}

	return c
}

func (c *HTTPCHecker) Check() (int, string) {
	resp, err := c.client.Get(c.URL)

	if err != nil {
		log.Infof("Error checking '%s' %s %s\n", c.Name, c.URL, err)
	}

	if err, ok := err.(net.Error); ok && err.Timeout() {
		return 0, "TIMEOUT"
	} else if err != nil {
		return 0, "DOWN"
	}

	defer resp.Body.Close()

	if resp.StatusCode == c.AcceptedStatusCode {
		return resp.StatusCode, "UP"
	}

	return resp.StatusCode, "FAILED"
}

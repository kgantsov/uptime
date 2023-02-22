package monitor

import (
	"net"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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
	URL                string
	Name               string
	AcceptedStatusCode int
	logger             *logrus.Logger
}

func NewHTTPCHecker(logger *logrus.Logger, name, url string, timeout, acceptedStatusCode int) *HTTPCHecker {
	client := http.Client{Timeout: time.Duration(timeout) * time.Second}

	c := &HTTPCHecker{
		client:             client,
		URL:                url,
		AcceptedStatusCode: acceptedStatusCode,
		logger:             logger,
	}

	return c
}

func (c *HTTPCHecker) Check() (int, string) {
	resp, err := c.client.Get(c.URL)

	if err != nil {
		log.Infof("Error checking '%s' %s %s\n", c.Name, c.URL, err)
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

package handler

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type StackdriverTimestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int   `json:"nanos"`
}

type StackdriverPayload struct {
	Name      string               `json:"name"`
	Date      string               `json:"date"`
	Message   string               `json:"message"`
	Timestamp StackdriverTimestamp `json:"timestamp"`
	Severity  string               `json:"severity"`
	RequestID string               `json:"request_id"`
}

type StackdriverFormatter struct {
}

func (f *StackdriverFormatter) Format(entry *log.Entry) ([]byte, error) {
	// Note this doesn't include Time, Level and Message which are available on
	// the Entry. Consult `godoc` on information about those fields or read the
	// source of the official loggers.

	now := time.Now()

	requestID := ""

	if _, ok := entry.Data["RequestID"]; ok {
		requestID = entry.Data["RequestID"].(string)
	}

	payload := StackdriverPayload{
		Date:      now.Format(time.RFC3339Nano),
		Message:   entry.Message,
		Severity:  entry.Level.String(),
		RequestID: requestID,
		Timestamp: StackdriverTimestamp{
			Seconds: now.Unix(),
			Nanos:   now.Nanosecond(),
		},
	}
	serialized, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

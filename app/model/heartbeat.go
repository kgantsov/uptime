package model

import "time"

type (
	Heartbeat struct {
		ID           uint      `json:"id" gorm:"primarykey"`
		ServiceID    uint      `json:"service_id"`
		ResponseTime int64     `json:"response_time"`
		Status       string    `json:"status"`
		StatusCode   int       `json:"status_code"`
		CreatedAt    time.Time `json:"created_at"`
	}

	HeartbeatPoint struct {
		ServiceID uint   `json:"service_id"`
		Latency   int64  `json:"latency"`
		Date      string `json:"date"`
	}
)

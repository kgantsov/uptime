package handler

import (
	"github.com/kgantsov/uptime/app/monitor"
	"gorm.io/gorm"
)

type (
	Handler struct {
		DB         *gorm.DB
		Dispatcher *monitor.Dispatcher
	}
)

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

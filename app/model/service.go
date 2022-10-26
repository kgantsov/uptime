package model

import (
	"database/sql"
	"time"
)

type Notification struct {
	Name string `json:"name" gorm:"primarykey"`

	CallbackType   string       `json:"callback_type"`
	CallbackChatID string       `json:"callback_chat_id"`
	Callback       string       `json:"callback"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at" gorm:"index"`
}

type Service struct {
	Model
	// ID            uint           `json:"id"`
	Name          string         `json:"name"`
	URL           string         `json:"url"`
	Notifications []Notification `gorm:"many2many:service_notifications;" json:"notifications"`
	Timeout       int            `json:"timeout"`
	CheckInterval int            `json:"check_interval"`
}

type ServiceNotification struct {
	ServiceID        int    `gorm:"primaryKey"`
	NotificationName string `gorm:"primaryKey"`
}

type AddService struct {
	Name          string            `json:"name"`
	Notifications []AddNotification `gorm:"many2many:service_notifications;" json:"notifications"`
	URL           string            `json:"url"`
	Timeout       int               `json:"timeout"`
	CheckInterval int               `json:"check_interval"`
}

type UpdateService struct {
	Name          string         `json:"name"`
	URL           string         `json:"url"`
	Notifications []Notification `gorm:"many2many:service_notifications;" json:"notifications"`
	Timeout       int            `json:"timeout"`
	CheckInterval int            `json:"check_interval"`
}

type AddNotification struct {
	Name           string `json:"name" gorm:"primarykey"`
	CallbackType   string `json:"callback_type"`
	CallbackChatID string `json:"callback_chat_id"`
	Callback       string `json:"callback"`
}
type UpdateNotification struct {
	CallbackType   string `json:"callback_type"`
	CallbackChatID string `json:"callback_chat_id"`
	Callback       string `json:"callback"`
}

package model

import (
	"database/sql"
	"time"
)

type Notification struct {
	Name string `json:"name" gorm:"primarykey"`

	CallbackType   string       `json:"callback_type,omitempty"`
	CallbackChatID string       `json:"callback_chat_id,omitempty"`
	Callback       string       `json:"callback,omitempty"`
	CreatedAt      time.Time    `json:"created_at,omitempty"`
	UpdatedAt      time.Time    `json:"updated_at,omitempty"`
	DeletedAt      sql.NullTime `json:"deleted_at,omitempty" gorm:"index"`
}

type Service struct {
	Model
	Name               string         `json:"name,omitempty"`
	URL                string         `json:"url,omitempty"`
	Enabled            bool           `json:"enabled,omitempty"`
	Notifications      []Notification `gorm:"many2many:service_notifications;" json:"notifications,omitempty"`
	Timeout            int            `json:"timeout,omitempty"`
	CheckInterval      int            `json:"check_interval,omitempty"`
	Retries            int            `json:"retries,omitempty"`
	AcceptedStatusCode int            `json:"accepted_status_code,omitempty"`
}

type ServiceNotification struct {
	ServiceID        int    `gorm:"primaryKey"`
	NotificationName string `gorm:"primaryKey"`
}

type AddService struct {
	Name               string            `json:"name"`
	Notifications      []AddNotification `gorm:"many2many:service_notifications;" json:"notifications,omitempty"`
	URL                string            `json:"url,omitempty"`
	Enabled            bool              `json:"enabled,omitempty"`
	Timeout            int               `json:"timeout,omitempty"`
	CheckInterval      int               `json:"check_interval,omitempty"`
	Retries            int               `json:"retries,omitempty"`
	AcceptedStatusCode int               `json:"accepted_status_code,omitempty"`
}

type UpdateService struct {
	Name               *string         `json:"name,omitempty"`
	URL                *string         `json:"url,omitempty"`
	Enabled            *bool           `json:"enabled,omitempty"`
	Notifications      *[]Notification `gorm:"many2many:service_notifications;" json:"notifications,omitempty"`
	Timeout            *int            `json:"timeout,omitempty"`
	CheckInterval      *int            `json:"check_interval,omitempty"`
	Retries            *int            `json:"retries,omitempty"`
	AcceptedStatusCode *int            `json:"accepted_status_code,omitempty"`
}

type AddNotification struct {
	Name           string `json:"name"`
	CallbackType   string `json:"callback_type,omitempty"`
	CallbackChatID string `json:"callback_chat_id,omitempty"`
	Callback       string `json:"callback,omitempty"`
}

type UpdateNotification struct {
	CallbackType   *string `json:"callback_type,omitempty"`
	CallbackChatID *string `json:"callback_chat_id,omitempty"`
	Callback       *string `json:"callback,omitempty"`
}

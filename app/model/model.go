package model

import (
	"database/sql"
	"time"
)

type Model struct {
	ID        uint         `json:"id,omitempty" gorm:"primarykey"`
	CreatedAt time.Time    `json:"created_at,omitempty"`
	UpdatedAt time.Time    `json:"updated_at,omitempty"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty" gorm:"index"`
}

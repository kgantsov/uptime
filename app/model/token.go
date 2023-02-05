package model

import "time"

type Token struct {
	Model

	UserID   uint
	User     User
	ExpireAt time.Time `json:"expire_at" gorm:"index"`
	Token    string    `json:"token,omitempty"`
}

type CreateToken struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

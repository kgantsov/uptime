package model

import (
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(
		&User{},
		&Token{},
		&Heartbeat{},
		&Service{},
		&Notification{},
		&ServiceNotification{},
	)
	db.SetupJoinTable(&Service{}, "Notifications", &ServiceNotification{})
}

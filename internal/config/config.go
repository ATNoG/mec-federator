package config

import "gorm.io/gorm"

var (
	db *gorm.DB
)

func Init() error {
	return nil
}

func GetLogger(p string) *Logger {
	// Initialize Logger
	logger := NewLogger(p)
	return logger
}

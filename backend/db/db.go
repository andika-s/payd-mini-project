package db

import (
	"payd-mini-project/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLite(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate schema
	if err := db.AutoMigrate(&model.Shift{}); err != nil {
		return nil, err
	}
	return db, nil
}

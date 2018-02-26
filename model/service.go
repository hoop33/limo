package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Service represents a hosting service like Github
type Service struct {
	gorm.Model
	Name        string
	LastSuccess time.Time
	Stars       []Star
}

// FindOrCreateServiceByName returns a service with the specified name, creating if necessary
func FindOrCreateServiceByName(db *gorm.DB, name string) (*Service, bool, error) {
	var service Service
	if db.Where("name = ?", name).First(&service).RecordNotFound() {
		service.Name = name
		err := db.Create(&service).Error
		return &service, true, err
	}
	return &service, false, nil
}

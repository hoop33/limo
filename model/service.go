package model

import "github.com/jinzhu/gorm"

// Service represents a hosting service like Github
type Service struct {
	gorm.Model
	Name  string
	Stars []Star
}

// GetOrCreateService returns a service with the specified name, creating if necessary
func GetOrCreateService(db *gorm.DB, name string) (*Service, error) {
	var service Service
	if db.Where("name = ?", name).First(&service).RecordNotFound() {
		service = Service{
			Name: name,
		}
		if err := db.Create(&service).Error; err != nil {
			return nil, err
		}
	}
	return &service, nil
}

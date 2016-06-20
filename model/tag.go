package model

import "github.com/jinzhu/gorm"

// Tag represents a tag in the database
type Tag struct {
	gorm.Model
	Name  string
	Stars []Star `gorm:"many2many:star_tags;"`
}

// FindOrCreateTagByName gets a tag by name, creating if it doesn't exist
func FindOrCreateTagByName(db *gorm.DB, name string) (*Tag, error) {
	var tag Tag
	if db.Where("name = ?", name).First(&tag).RecordNotFound() {
		tag.Name = name
		err := db.Create(&tag).Error
		return &tag, err
	}
	return &tag, nil
}

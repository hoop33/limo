package model

import (
	"strings"

	"github.com/jinzhu/gorm"
)

// Tag represents a tag in the database
type Tag struct {
	gorm.Model
	Name  string
	Stars []Star `gorm:"many2many:star_tags;"`
}

// FindTags finds all tags
func FindTags(db *gorm.DB) ([]Tag, error) {
	var tags []Tag
	db.Order("name").Find(&tags)
	return tags, db.Error
}

// FindOrCreateTagByName gets a tag by name, creating if it doesn't exist
func FindOrCreateTagByName(db *gorm.DB, name string) (*Tag, bool, error) {
	var tag Tag
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&tag).RecordNotFound() {
		tag.Name = name
		err := db.Create(&tag).Error
		return &tag, true, err
	}
	return &tag, false, nil
}

// LoadStars loads the stars for a tag
func (tag *Tag) LoadStars(db *gorm.DB) error {
	return db.Model(tag).Association("Stars").Find(&tag.Stars).Error
}

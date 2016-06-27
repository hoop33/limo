package model

import (
	"fmt"
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

// FindTagByName finds a tag by name
func FindTagByName(db *gorm.DB, name string) (*Tag, error) {
	var tag Tag
	if db.Where("lower(name) = ?", strings.ToLower(name)).First(&tag).RecordNotFound() {
		return nil, db.Error
	}
	return &tag, db.Error
}

// FindOrCreateTagByName finds a tag by name, creating if it doesn't exist
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
	// Make sure tag exists in database, or we will panic
	var existing Tag
	if db.Where("id = ?", tag.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("Tag '%d' not found", tag.ID)
	}
	return db.Model(tag).Association("Stars").Find(&tag.Stars).Error
}

// Rename renames a tag -- new name must not already exist
func (tag *Tag) Rename(db *gorm.DB, name string) error {
	existing, err := FindTagByName(db, name)
	if err != nil {
		return err
	}
	if existing != nil {
		return fmt.Errorf("Tag '%s' already exists", existing.Name)
	}
	tag.Name = name
	return db.Save(tag).Error
}

// Delete deletes a tag and disassociates it from any stars
func (tag *Tag) Delete(db *gorm.DB) error {
	if err := db.Model(tag).Association("Stars").Clear().Error; err != nil {
		return err
	}
	return db.Delete(tag).Error
}

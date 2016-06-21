package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
)

// Star represents a starred repository
type Star struct {
	gorm.Model
	RemoteID    string
	Name        *string
	FullName    *string
	Description *string
	Homepage    *string
	URL         *string
	Language    *string
	Stargazers  int
	ServiceID   uint
	Tags        []Tag `gorm:"many2many:star_tags;"`
}

// StarResult wraps a star and an error
type StarResult struct {
	Star  *Star
	Error error
}

// NewStarFromGithub creates a Star from a Github star
func NewStarFromGithub(star github.Repository) (*Star, error) {
	// Require the GitHub ID
	if star.ID == nil {
		return nil, errors.New("ID from GitHub is required")
	}

	// Set stargazers count to 0 if nil
	stargazersCount := 0
	if star.StargazersCount != nil {
		stargazersCount = *star.StargazersCount
	}

	return &Star{
		RemoteID:    strconv.Itoa(*star.ID),
		Name:        star.Name,
		FullName:    star.FullName,
		Description: star.Description,
		Homepage:    star.Homepage,
		URL:         star.CloneURL,
		Language:    star.Language,
		Stargazers:  stargazersCount,
	}, nil
}

// StarCopy copies values from src to dest
func StarCopy(src *Star, dest *Star) {
	dest.Name = src.Name
	dest.FullName = src.FullName
	dest.Description = src.Description
	dest.Homepage = src.Homepage
	dest.URL = src.URL
	dest.Language = src.Language
	dest.Stargazers = src.Stargazers
}

// CreateOrUpdateStar creates or updates a star and returns true if the star was created (vs updated)
func CreateOrUpdateStar(db *gorm.DB, star *Star, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var existing Star
	if db.Where("remote_id = ? AND service_id = ?", star.RemoteID, service.ID).First(&existing).RecordNotFound() {
		star.ServiceID = service.ID
		err := db.Create(star).Error
		return err == nil, err
	}
	StarCopy(star, &existing)
	return false, db.Save(&existing).Error
}

// FindStars finds all stars
func FindStars(db *gorm.DB) ([]Star, error) {
	var stars []Star
	db.Order("full_name").Find(&stars)
	return stars, db.Error
}

// FindStarsWithLanguage finds stars with the specified language
func FindStarsWithLanguage(db *gorm.DB, language string) ([]Star, error) {
	var stars []Star
	db.Where("lower(language) = ?", strings.ToLower(language)).Order("full_name").Find(&stars)
	return stars, db.Error
}

// FuzzyFindStarsWithName finds stars with approximate matching for full name and name
func FuzzyFindStarsWithName(db *gorm.DB, name string) ([]Star, error) {
	// Try each of these, and as soon as we hit, return
	// 1. Exact match full name
	// 2. Exact match name
	// 3. Case-insensitive full name
	// 4. Case-insensitive name
	// 5. Case-insensitive like full name
	// 6. Case-insensitive like name
	var stars []Star
	db.Where("full_name = ?", name).Order("full_name").Find(&stars)
	if len(stars) == 0 {
		db.Where("name = ?", name).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("lower(full_name) = ?", strings.ToLower(name)).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("lower(name) = ?", strings.ToLower(name)).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("full_name LIKE ?", strings.ToLower(fmt.Sprintf("%%%s%%", name))).Order("full_name").Find(&stars)
	}
	if len(stars) == 0 {
		db.Where("name LIKE ?", strings.ToLower(fmt.Sprintf("%%%s%%", name))).Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FindLanguages finds all languages
func FindLanguages(db *gorm.DB) ([]string, error) {
	var languages []string
	db.Table("stars").Order("language").Pluck("distinct(language)", &languages)
	return languages, db.Error
}

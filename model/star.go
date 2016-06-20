package model

import (
	"strconv"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
)

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
}

type StarResult struct {
	Star  *Star
	Error error
}

// NewStarFromGithub creates a Star from a Github star
func NewStarFromGithub(star github.Repository) *Star {
	return &Star{
		RemoteID:    strconv.Itoa(*star.ID),
		Name:        star.Name,
		FullName:    star.FullName,
		Description: star.Description,
		Homepage:    star.Homepage,
		URL:         star.CloneURL,
		Language:    star.Language,
		Stargazers:  *star.StargazersCount,
	}
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
	} else {
		StarCopy(star, &existing)
		return false, db.Save(&existing).Error
	}
}

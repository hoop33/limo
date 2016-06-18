package model

import (
	"strconv"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
)

type Star struct {
	gorm.Model
	ServiceID   uint
	RemoteID    string
	Name        *string
	FullName    *string
	Description *string
	Homepage    *string
	URL         *string
	Language    *string
	Stargazers  int
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

// CreateOrUpdateStar creates or updates a star
func CreateOrUpdateStar(db *gorm.DB, star *Star, service *Service) (bool, error) {
	// Get existing by remote ID and service ID
	var old Star
	if db.Where("remote_id = ? AND service_id = ?", star.RemoteID, service.ID).First(&old).RecordNotFound() {
		star.ServiceID = service.ID
		if err := db.Create(star).Error; err != nil {
			return false, err
		} else {
			return true, nil
		}
	} else {
		star.ID = old.ID
		if err := db.Update(star).Error; err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

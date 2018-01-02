package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/skratchdot/open-golang/open"
	"github.com/xanzy/go-gitlab"
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
	StarredAt   time.Time
	ServiceID   uint
	Tags        []Tag `gorm:"many2many:star_tags;"`
}

// StarResult wraps a star and an error
type StarResult struct {
	Star  *Star
	Error error
}

// NewStarFromGithub creates a Star from a Github star
func NewStarFromGithub(timestamp *github.Timestamp, star github.Repository) (*Star, error) {
	// Require the GitHub ID
	if star.ID == nil {
		return nil, errors.New("ID from GitHub is required")
	}

	// Set stargazers count to 0 if nil
	stargazersCount := 0
	if star.StargazersCount != nil {
		stargazersCount = *star.StargazersCount
	}

	starredAt := time.Now()
	if timestamp != nil {
		starredAt = timestamp.Time
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
		StarredAt:   starredAt,
	}, nil
}

// NewStarFromGitlab creates a Star from a Gitlab star
func NewStarFromGitlab(star gitlab.Project) (*Star, error) {
	return &Star{
		RemoteID:    strconv.Itoa(star.ID),
		Name:        &star.Name,
		FullName:    &star.NameWithNamespace,
		Description: &star.Description,
		Homepage:    &star.WebURL,
		URL:         &star.HTTPURLToRepo,
		Language:    nil,
		Stargazers:  star.StarCount,
		StarredAt:   time.Now(), // OK, so this is a lie, but not in payload
	}, nil
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
	star.ID = existing.ID
	star.ServiceID = service.ID
	star.CreatedAt = existing.CreatedAt
	return false, db.Save(star).Error
}

// FindStarByID finds a star by ID
func FindStarByID(db *gorm.DB, ID uint) (*Star, error) {
	var star Star
	if db.First(&star, ID).RecordNotFound() {
		return nil, fmt.Errorf("star '%d' not found", ID)
	}
	return &star, db.Error
}

// FindStars finds all stars
func FindStars(db *gorm.DB, match string) ([]Star, error) {
	var stars []Star
	if match != "" {
		db.Where("full_name LIKE ?",
			strings.ToLower(fmt.Sprintf("%%%s%%", match))).Order("full_name").Find(&stars)
	} else {
		db.Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FindUntaggedStars finds stars without any tags
func FindUntaggedStars(db *gorm.DB, match string) ([]Star, error) {
	var stars []Star
	if match != "" {
		db.Raw(`
			SELECT *
			FROM STARS S
			WHERE S.DELETED_AT IS NULL
			AND S.FULL_NAME LIKE ?
			AND S.ID NOT IN (
				SELECT STAR_ID
				FROM STAR_TAGS
			) ORDER BY S.FULL_NAME`,
			fmt.Sprintf("%%%s%%", strings.ToLower(match))).Scan(&stars)
	} else {
		db.Raw(`
			SELECT *
			FROM STARS S
			WHERE S.DELETED_AT IS NULL
			AND S.ID NOT IN (
				SELECT STAR_ID
				FROM STAR_TAGS
			) ORDER BY S.FULL_NAME`).Scan(&stars)
	}
	return stars, db.Error
}

// FindStarsByLanguageAndOrTag finds stars with the specified language and/or the specified tag
func FindStarsByLanguageAndOrTag(db *gorm.DB, match string, language string, tagName string, union bool) ([]Star, error) {
	operator := "AND"
	if union {
		operator = "OR"
	}

	var stars []Star
	if match != "" {
		db.Raw(fmt.Sprintf(`
			SELECT * 
			FROM STARS S, TAGS T 
			INNER JOIN STAR_TAGS ST ON S.ID = ST.STAR_ID 
			INNER JOIN TAGS ON ST.TAG_ID = T.ID 
			WHERE S.DELETED_AT IS NULL
			AND T.DELETED_AT IS NULL
			AND LOWER(S.FULL_NAME) LIKE ? 
			AND (LOWER(T.NAME) = ? 
			%s LOWER(S.LANGUAGE) = ?) 
			GROUP BY ST.STAR_ID 
			ORDER BY S.FULL_NAME`, operator),
			fmt.Sprintf("%%%s%%", strings.ToLower(match)),
			strings.ToLower(tagName),
			strings.ToLower(language)).Scan(&stars)
	} else {
		db.Raw(fmt.Sprintf(`
			SELECT * 
			FROM STARS S, TAGS T 
			INNER JOIN STAR_TAGS ST ON S.ID = ST.STAR_ID 
			INNER JOIN TAGS ON ST.TAG_ID = T.ID 
			WHERE S.DELETED_AT IS NULL
			AND T.DELETED_AT IS NULL
			AND LOWER(T.NAME) = ? 
			%s LOWER(S.LANGUAGE) = ? 
			GROUP BY ST.STAR_ID 
			ORDER BY S.FULL_NAME`, operator),
			strings.ToLower(tagName),
			strings.ToLower(language)).Scan(&stars)
	}
	return stars, db.Error
}

// FindStarsByLanguage finds stars with the specified language
func FindStarsByLanguage(db *gorm.DB, match string, language string) ([]Star, error) {
	var stars []Star
	if match != "" {
		db.Where("full_name LIKE ? AND lower(language) = ?",
			strings.ToLower(fmt.Sprintf("%%%s%%", match)),
			strings.ToLower(language)).Order("full_name").Find(&stars)
	} else {
		db.Where("lower(language) = ?",
			strings.ToLower(language)).Order("full_name").Find(&stars)
	}
	return stars, db.Error
}

// FuzzyFindStarsByName finds stars with approximate matching for full name and name
func FuzzyFindStarsByName(db *gorm.DB, name string) ([]Star, error) {
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

// AddTag adds a tag to a star
func (star *Star) AddTag(db *gorm.DB, tag *Tag) error {
	star.Tags = append(star.Tags, *tag)
	return db.Save(star).Error
}

// LoadTags loads the tags for a star
func (star *Star) LoadTags(db *gorm.DB) error {
	// Make sure star exists in database, or we will panic
	var existing Star
	if db.Where("id = ?", star.ID).First(&existing).RecordNotFound() {
		return fmt.Errorf("star '%d' not found", star.ID)
	}
	return db.Model(star).Association("Tags").Find(&star.Tags).Error
}

// RemoveAllTags removes all tags for a star
func (star *Star) RemoveAllTags(db *gorm.DB) error {
	return db.Model(star).Association("Tags").Clear().Error
}

// RemoveTag removes a tag from a star
func (star *Star) RemoveTag(db *gorm.DB, tag *Tag) error {
	return db.Model(star).Association("Tags").Delete(tag).Error
}

// HasTag returns whether a star has a tag. Note that you must call LoadTags first -- no reason to incur a database call each time
func (star *Star) HasTag(tag *Tag) bool {
	if len(star.Tags) > 0 {
		for _, t := range star.Tags {
			if t.Name == tag.Name {
				return true
			}
		}
	}
	return false
}

// Index adds the star to the index
func (star *Star) Index(index bleve.Index, db *gorm.DB) error {
	if err := star.LoadTags(db); err != nil {
		return err
	}
	return index.Index(fmt.Sprintf("%d", star.ID), star)
}

// OpenInBrowser opens the star in the browser
func (star *Star) OpenInBrowser(preferHomepage bool) error {
	var URL string
	if preferHomepage && star.Homepage != nil && *star.Homepage != "" {
		URL = *star.Homepage
	} else if star.URL != nil && *star.URL != "" {
		URL = *star.URL
	} else {
		if star.Name != nil {
			return fmt.Errorf("no URL for star '%s'", *star.Name)
		}
		return errors.New("no URL for star")
	}
	return open.Start(URL)
}

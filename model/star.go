package model

import (
	"fmt"

	"github.com/google/go-github/github"
)

type Star struct {
	ID          int    `db:"id"`
	RemoteID    string `db:"remote_id"`
	Name        string `db:"name"`
	FullName    string `db:"full_name"`
	Description string `db:"description"`
	Homepage    string `db:"homepage"`
	URL         string `db:"url"`
	Language    string `db:"language"`
	Stargazers  int    `db:"stargazers"`
	OwnerID     int    `db:"owner_id"`
	ServiceID   int    `db:"service_id"`
}

type StarResult struct {
	Star  *Star
	Error error
}

// String returns a string representation of a star
func (s *Star) String() string {
	return fmt.Sprintf("%s: %s (%s)", s.RemoteID, s.FullName, s.Name)
}

// NewStarFromGithub creates a Star from a Github star
func NewStarFromGithub(star github.Repository) *Star {
	return &Star{
		RemoteID: string(*star.ID),
		Name:     *star.Name,
		FullName: *star.FullName,
	}
}

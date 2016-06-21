package model

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestNewStarFromGithubShouldCopyFields(t *testing.T) {
	id := 33
	name := "larry-bird"
	fullName := "celtics/larry-bird"
	description := "larry legend"
	homepage := "http ://www.nba.com/celtics/"
	url := "http ://www.nba.com/pacers/"
	language := "hoosier"
	stargazersCount := 10000

	github := github.Repository{
		ID:              &id,
		Name:            &name,
		FullName:        &fullName,
		Description:     &description,
		Homepage:        &homepage,
		CloneURL:        &url,
		Language:        &language,
		StargazersCount: &stargazersCount,
	}

	star, err := NewStarFromGithub(github)
	assert.Nil(t, err)
	assert.Equal(t, "33", star.RemoteID)
	assert.Equal(t, name, *star.Name)
	assert.Equal(t, fullName, *star.FullName)
	assert.Equal(t, description, *star.Description)
	assert.Equal(t, homepage, *star.Homepage)
	assert.Equal(t, url, *star.URL)
	assert.Equal(t, language, *star.Language)
	assert.Equal(t, stargazersCount, star.Stargazers)
}

func TestNewStarFromGithubShouldHandleEmpty(t *testing.T) {
	star, err := NewStarFromGithub(github.Repository{})
	assert.NotNil(t, err)
	assert.Equal(t, "ID from GitHub is required", err.Error())
	assert.Nil(t, star)
}

func TestNewStarFromGithubShouldHandleOnlyID(t *testing.T) {
	id := 33
	star, err := NewStarFromGithub(github.Repository{
		ID: &id,
	})
	assert.Nil(t, err)
	assert.NotNil(t, star)
	assert.Equal(t, "33", star.RemoteID)
}

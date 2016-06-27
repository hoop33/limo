package model

import (
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestNewStarFromGithubShouldCopyFields(t *testing.T) {
	id := 33
	name := "larry-bird"
	fullName := "celtics/larry-bird"
	description := "larry legend"
	homepage := "http://www.nba.com/celtics/"
	url := "http://www.nba.com/pacers/"
	language := "hoosier"
	stargazersCount := 10000

	timestamp := github.Timestamp{
		Time: time.Now(),
	}

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

	star, err := NewStarFromGithub(&timestamp, github)
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
	star, err := NewStarFromGithub(&github.Timestamp{}, github.Repository{})
	assert.NotNil(t, err)
	assert.Equal(t, "ID from GitHub is required", err.Error())
	assert.Nil(t, star)
}

func TestNewStarFromGithubShouldHandleOnlyID(t *testing.T) {
	id := 33
	star, err := NewStarFromGithub(&github.Timestamp{}, github.Repository{
		ID: &id,
	})
	assert.Nil(t, err)
	assert.NotNil(t, star)
	assert.Equal(t, "33", star.RemoteID)
}

func TestFuzzyFindStarsByNameShouldFuzzyFind(t *testing.T) {
	fullName := "Apple/Baker"
	name := "Charlie"

	star := Star{
		FullName: &fullName,
		Name:     &name,
	}
	assert.Nil(t, db.Create(&star).Error)

	stars, err := FuzzyFindStarsByName(db, "Apple/Baker")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stars))
	assert.Equal(t, fullName, *stars[0].FullName)
	assert.Equal(t, name, *stars[0].Name)

	stars, err = FuzzyFindStarsByName(db, "Charlie")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stars))
	assert.Equal(t, fullName, *stars[0].FullName)
	assert.Equal(t, name, *stars[0].Name)

	stars, err = FuzzyFindStarsByName(db, "apple/baker")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stars))
	assert.Equal(t, fullName, *stars[0].FullName)
	assert.Equal(t, name, *stars[0].Name)

	stars, err = FuzzyFindStarsByName(db, "charlie")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stars))
	assert.Equal(t, fullName, *stars[0].FullName)
	assert.Equal(t, name, *stars[0].Name)

	stars, err = FuzzyFindStarsByName(db, "apple")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stars))
	assert.Equal(t, fullName, *stars[0].FullName)
	assert.Equal(t, name, *stars[0].Name)

	stars, err = FuzzyFindStarsByName(db, "harl")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stars))
	assert.Equal(t, fullName, *stars[0].FullName)
	assert.Equal(t, name, *stars[0].Name)

	stars, err = FuzzyFindStarsByName(db, "boogers")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(stars))
}

func TestAddTagShouldAddTag(t *testing.T) {
	tag, _, err := FindOrCreateTagByName(db, "celtics")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "celtics", tag.Name)

	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Isaiah Thomas" // Not a typo
	star := &Star{
		RemoteID: "remoteID",
		Name:     &name,
	}
	_, err = CreateOrUpdateStar(db, star, service)
	assert.Nil(t, err)

	err = star.AddTag(db, tag)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star.Tags))
	assert.Equal(t, "celtics", star.Tags[0].Name)

	stars, err := FuzzyFindStarsByName(db, "Isaiah Thomas")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stars))
	assert.Equal(t, "Isaiah Thomas", *stars[0].Name)

	err = stars[0].LoadTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stars[0].Tags))
	assert.Equal(t, "celtics", stars[0].Tags[0].Name)
}

func TestHasTagShouldReturnFalseWhenNoTags(t *testing.T) {
	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	star := &Star{
		RemoteID: "brown",
		Name:     &name,
	}
	tag, _, err := FindOrCreateTagByName(db, "bucks")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "bucks", tag.Name)

	_, err = CreateOrUpdateStar(db, star, service)
	assert.Nil(t, err)

	err = star.LoadTags(db)
	assert.Nil(t, err)

	assert.False(t, star.HasTag(tag))
}

func TestHasTagShouldReturnFalseWhenTagIsNil(t *testing.T) {
	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	star := &Star{
		RemoteID: "brown",
		Name:     &name,
	}
	_, err = CreateOrUpdateStar(db, star, service)
	assert.Nil(t, err)

	err = star.LoadTags(db)
	assert.Nil(t, err)

	assert.False(t, star.HasTag(nil))
}

func TestHasTagShouldReturnFalseWhenDoesNotHaveTag(t *testing.T) {
	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	star := &Star{
		RemoteID: "brown",
		Name:     &name,
	}
	_, err = CreateOrUpdateStar(db, star, service)
	assert.Nil(t, err)

	bucks, _, err := FindOrCreateTagByName(db, "bucks")
	assert.Nil(t, err)
	assert.NotNil(t, bucks)
	assert.Equal(t, "bucks", bucks.Name)

	celtics, _, err := FindOrCreateTagByName(db, "celtics")
	assert.Nil(t, err)
	assert.NotNil(t, celtics)
	assert.Equal(t, "celtics", celtics.Name)

	err = star.AddTag(db, celtics)
	assert.Nil(t, err)

	err = star.LoadTags(db)
	assert.Nil(t, err)

	assert.False(t, star.HasTag(bucks))
}

func TestHasTagShouldReturnTrueWhenHasOnlyTag(t *testing.T) {
	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	star := &Star{
		RemoteID: "brown",
		Name:     &name,
	}
	_, err = CreateOrUpdateStar(db, star, service)
	assert.Nil(t, err)

	celtics, _, err := FindOrCreateTagByName(db, "celtics")
	assert.Nil(t, err)
	assert.NotNil(t, celtics)
	assert.Equal(t, "celtics", celtics.Name)

	err = star.AddTag(db, celtics)
	assert.Nil(t, err)

	err = star.LoadTags(db)
	assert.Nil(t, err)

	assert.True(t, star.HasTag(celtics))
}

func TestHasTagShouldReturnTrueWhenHasTag(t *testing.T) {
	service, _, err := FindOrCreateServiceByName(db, "nba")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nba", service.Name)

	name := "Jaylen Brown"
	star := &Star{
		RemoteID: "brown",
		Name:     &name,
	}
	_, err = CreateOrUpdateStar(db, star, service)
	assert.Nil(t, err)

	draft, _, err := FindOrCreateTagByName(db, "2016-draft")
	assert.Nil(t, err)
	assert.NotNil(t, draft)
	assert.Equal(t, "2016-draft", draft.Name)

	celtics, _, err := FindOrCreateTagByName(db, "celtics")
	assert.Nil(t, err)
	assert.NotNil(t, celtics)
	assert.Equal(t, "celtics", celtics.Name)

	err = star.AddTag(db, celtics)
	assert.Nil(t, err)

	err = star.LoadTags(db)
	assert.Nil(t, err)

	assert.True(t, star.HasTag(celtics))
}

func TestLoadTagsShouldReturnErrorWhenStarNotInDatabase(t *testing.T) {
	name := "not in db"
	star := &Star{
		RemoteID: "not in db",
		Name:     &name,
	}

	err := star.LoadTags(db)
	assert.NotNil(t, err)
	assert.Equal(t, "Star '0' not found", err.Error())
}

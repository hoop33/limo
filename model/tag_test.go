package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOrCreateTagByNameShouldCreateTag(t *testing.T) {
	tag, created, err := FindOrCreateTagByName(db, "my-tag")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.True(t, created)
	assert.Equal(t, "my-tag", tag.Name)

	var check Tag
	db.Where("name = ?", "my-tag").First(&check)
	assert.Equal(t, "my-tag", check.Name)
}

func TestFindOrCreateTagShouldNotCreateDuplicateNames(t *testing.T) {
	tag, created, err := FindOrCreateTagByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.True(t, created)
	assert.Equal(t, "foo", tag.Name)

	tag, created, err = FindOrCreateTagByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.False(t, created)
	assert.Equal(t, "foo", tag.Name)

	var tags []Tag
	db.Where("name = ?", "foo").Find(&tags)
	assert.Equal(t, 1, len(tags))
}

func TestFindTagByNameShouldReturnNilIfNotExists(t *testing.T) {
	tag, err := FindTagByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, tag)
}

func TestFindTagByNameShouldFindTag(t *testing.T) {
	tag, created, err := FindOrCreateTagByName(db, "creating a new tag")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.True(t, created)
	assert.Equal(t, "creating a new tag", tag.Name)

	newTag, err := FindTagByName(db, "creating a new tag")
	assert.Nil(t, err)
	assert.Equal(t, "creating a new tag", newTag.Name)
}

func TestRenameTagShouldRenameTag(t *testing.T) {
	tag, created, err := FindOrCreateTagByName(db, "old name")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.True(t, created)
	assert.Equal(t, "old name", tag.Name)

	err = tag.Rename(db, "new name")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "new name", tag.Name)
}

func TestRenameTagToExistingNameShouldReturnError(t *testing.T) {
	first, created, err := FindOrCreateTagByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.True(t, created)
	assert.Equal(t, "first", first.Name)

	second, created, err := FindOrCreateTagByName(db, "second")
	assert.Nil(t, err)
	assert.NotNil(t, second)
	assert.True(t, created)
	assert.Equal(t, "second", second.Name)

	err = second.Rename(db, "first")
	assert.NotNil(t, err)
	assert.Equal(t, "second", second.Name)

	err = second.Rename(db, "First")
	assert.NotNil(t, err)
	assert.Equal(t, "second", second.Name)

	err = second.Rename(db, "FIRST")
	assert.NotNil(t, err)
	assert.Equal(t, "second", second.Name)
}

func TestDeleteTagShouldDeleteTag(t *testing.T) {
	tag, created, err := FindOrCreateTagByName(db, "to delete")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.True(t, created)
	assert.Equal(t, "to delete", tag.Name)

	err = tag.Delete(db)
	assert.Nil(t, err)

	deleted, err := FindTagByName(db, "to delete")
	assert.Nil(t, err)
	assert.Nil(t, deleted)
}

func TestDeleteTagShouldDeleteAssociationsToStars(t *testing.T) {
	service, _, err := FindOrCreateServiceByName(db, "nfl")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "nfl", service.Name)

	name1 := "Allen Hurns"
	star1 := &Star{
		RemoteID: "88",
		Name:     &name1,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)

	name2 := "Allen Robinson"
	star2 := &Star{
		RemoteID: "15",
		Name:     &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)

	tag, _, err := FindOrCreateTagByName(db, "jaguars")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "jaguars", tag.Name)

	err = star1.AddTag(db, tag)
	assert.Nil(t, err)

	err = star2.AddTag(db, tag)
	assert.Nil(t, err)

	err = star1.LoadTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star1.Tags))
	assert.Equal(t, "jaguars", star1.Tags[0].Name)

	err = star2.LoadTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(star2.Tags))
	assert.Equal(t, "jaguars", star2.Tags[0].Name)

	err = tag.Delete(db)
	assert.Nil(t, err)

	err = star1.LoadTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star1.Tags))

	err = star2.LoadTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(star2.Tags))
}

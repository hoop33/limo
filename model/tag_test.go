package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindTagsShouldBeEmptyWhenNoTags(t *testing.T) {
	clearDB()

	tags, err := FindTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tags))
}

func TestFindTagsShouldFindATag(t *testing.T) {
	clearDB()

	tag, _, err := FindOrCreateTagByName(db, "solo")
	assert.Nil(t, err)
	assert.NotNil(t, tag)

	tags, err := FindTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, "solo", tags[0].Name)
}

func TestFindTagsShouldSortTagsByName(t *testing.T) {
	clearDB()

	_, _, err := FindOrCreateTagByName(db, "delta")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTagByName(db, "baker")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTagByName(db, "apple")
	assert.Nil(t, err)
	_, _, err = FindOrCreateTagByName(db, "charlie")
	assert.Nil(t, err)

	tags, err := FindTags(db)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(tags))
	assert.Equal(t, "apple", tags[0].Name)
	assert.Equal(t, "baker", tags[1].Name)
	assert.Equal(t, "charlie", tags[2].Name)
	assert.Equal(t, "delta", tags[3].Name)
}

func TestFindOrCreateTagByNameShouldCreateTag(t *testing.T) {
	clearDB()

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
	clearDB()

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
	clearDB()

	tag, err := FindTagByName(db, "this does not exist")
	assert.Nil(t, err)
	assert.Nil(t, tag)
}

func TestFindTagByNameShouldFindTag(t *testing.T) {
	clearDB()

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
	clearDB()

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
	clearDB()

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

func TestRenameTagByChangingCaseShouldRenameTag(t *testing.T) {
	clearDB()

	first, _, err := FindOrCreateTagByName(db, "first")
	assert.Nil(t, err)
	assert.NotNil(t, first)
	assert.Equal(t, "first", first.Name)

	err = first.Rename(db, "First")
	assert.Nil(t, err)
	assert.Equal(t, "First", first.Name)

	err = first.Rename(db, "FIRST")
	assert.Nil(t, err)
	assert.Equal(t, "FIRST", first.Name)
}

func TestRenameTagByChangingToSameNameShouldReturnError(t *testing.T) {
	clearDB()

	same, _, err := FindOrCreateTagByName(db, "same")
	assert.Nil(t, err)
	assert.NotNil(t, same)
	assert.Equal(t, "same", same.Name)

	err = same.Rename(db, "same")
	assert.NotNil(t, err)
}

func TestDeleteTagShouldDeleteTag(t *testing.T) {
	clearDB()

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
	clearDB()

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

func TestLoadStarsShouldReturnErrorWhenTagNotInDatabase(t *testing.T) {
	clearDB()

	tag := &Tag{
		Name: "not in db",
	}

	err := tag.LoadStars(db, "")
	assert.NotNil(t, err)
	assert.Equal(t, "Tag '0' not found", err.Error())
}

func TestLoadStarsShouldLoadNoStarsWhenTagHasNoStars(t *testing.T) {
	clearDB()

	tag, _, err := FindOrCreateTagByName(db, "tag")
	assert.Nil(t, err)
	assert.NotNil(t, tag)

	err = tag.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tag.Stars))
}

func TestLoadStarsShouldFillStars(t *testing.T) {
	clearDB()

	tag, _, err := FindOrCreateTagByName(db, "tag")
	assert.Nil(t, err)
	assert.NotNil(t, tag)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddTag(db, tag)
	assert.Nil(t, err)

	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddTag(db, tag)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(tag.Stars))
	err = tag.LoadStars(db, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tag.Stars))
}

func TestLoadStarsShouldFillStarsWithMatch(t *testing.T) {
	clearDB()

	tag, _, err := FindOrCreateTagByName(db, "tag")
	assert.Nil(t, err)
	assert.NotNil(t, tag)

	service, _, err := FindOrCreateServiceByName(db, "svc")
	assert.Nil(t, err)

	name1 := "Jacksonville Jaguars"
	star1 := &Star{
		RemoteID:  "1",
		ServiceID: service.ID,
		FullName:  &name1,
	}
	_, err = CreateOrUpdateStar(db, star1, service)
	assert.Nil(t, err)
	err = star1.AddTag(db, tag)
	assert.Nil(t, err)

	name2 := "Jacksonville Suns"
	star2 := &Star{
		RemoteID:  "2",
		ServiceID: service.ID,
		FullName:  &name2,
	}
	_, err = CreateOrUpdateStar(db, star2, service)
	assert.Nil(t, err)
	err = star2.AddTag(db, tag)
	assert.Nil(t, err)

	name3 := "Florida Gators"
	star3 := &Star{
		RemoteID:  "3",
		ServiceID: service.ID,
		FullName:  &name3,
	}
	_, err = CreateOrUpdateStar(db, star3, service)
	assert.Nil(t, err)
	err = star3.AddTag(db, tag)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tag.Stars))
	err = tag.LoadStars(db, "jacksonville")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tag.Stars))
	assert.Equal(t, "Jacksonville Jaguars", *tag.Stars[0].FullName)
	assert.Equal(t, "Jacksonville Suns", *tag.Stars[1].FullName)
}

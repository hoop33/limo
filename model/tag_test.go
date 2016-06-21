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

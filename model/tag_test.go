package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOrCreateTagByNameShouldCreateTag(t *testing.T) {
	tag, err := FindOrCreateTagByName(db, "my-tag")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "my-tag", tag.Name)

	var check Tag
	db.Where("name = ?", "my-tag").First(&check)
	assert.Equal(t, "my-tag", check.Name)
}

package model

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	mkDB()
	run := m.Run()
	//rmDB()
	os.Exit(run)
}

func TestFindOrCreateTagByNameShouldCreateTag(t *testing.T) {
	tag, err := FindOrCreateTagByName(db, "my-tag")
	assert.Nil(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "my-tag", tag.Name)

	var check Tag
	db.Where("name = ?", "my-tag").First(&check)
	assert.Equal(t, "my-tag", check.Name)
}

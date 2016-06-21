package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindOrCreateServiceByNameShouldCreateService(t *testing.T) {
	service, created, err := FindOrCreateServiceByName(db, "lebron-james")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.True(t, created)
	assert.Equal(t, "lebron-james", service.Name)

	var check Service
	db.Where("name = ?", "lebron-james").First(&check)
	assert.Equal(t, "lebron-james", check.Name)
}

func TestFindOrCreateServiceByNameShouldNotCreateDuplicateNames(t *testing.T) {
	service, created, err := FindOrCreateServiceByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.True(t, created)
	assert.Equal(t, "foo", service.Name)

	service, created, err = FindOrCreateServiceByName(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.False(t, created)
	assert.Equal(t, "foo", service.Name)

	var services []Service
	db.Where("name = ?", "foo").Find(&services)
	assert.Equal(t, 1, len(services))
}

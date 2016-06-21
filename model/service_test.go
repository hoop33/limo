package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOrCreateServiceShouldCreateService(t *testing.T) {
	service, err := GetOrCreateService(db, "lebron-james")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "lebron-james", service.Name)

	var check Service
	db.Where("name = ?", "lebron-james").First(&check)
	assert.Equal(t, "lebron-james", check.Name)
}

func TestGetOrCreateServiceShouldNotCreateDuplicateNames(t *testing.T) {
	service, err := GetOrCreateService(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "foo", service.Name)

	service, err = GetOrCreateService(db, "foo")
	assert.Nil(t, err)
	assert.NotNil(t, service)
	assert.Equal(t, "foo", service.Name)

	var services []Service
	db.Where("name = ?", "foo").Find(&services)
	assert.Equal(t, 1, len(services))
}

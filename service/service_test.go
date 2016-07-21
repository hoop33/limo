package service

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceNameShouldDropPackage(t *testing.T) {
	nf := &NotFound{}
	name := Name(nf)
	assert.Equal(t, "notfound", name)
}

func TestForNameShouldReturnErrorWhenNoService(t *testing.T) {
	svc, err := ForName("foo")
	assert.NotNil(t, err)
	assert.Equal(t, "Service 'foo' not found", err.Error())
	assert.Equal(t, "*service.NotFound", reflect.TypeOf(svc).String())
}

func TestForNameShouldReturnService(t *testing.T) {
	svc, err := ForName("github")
	assert.Nil(t, err)
	assert.Equal(t, "*service.Github", reflect.TypeOf(svc).String())
}

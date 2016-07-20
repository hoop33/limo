package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundLoginShouldReturnError(t *testing.T) {
	nf := &NotFound{}
	token, err := nf.Login()
	assert.NotNil(t, err)
	assert.Equal(t, "Service not found", err.Error())
	assert.Equal(t, "", token)
}

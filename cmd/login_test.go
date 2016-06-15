package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, LoginCmd.Use)
}

func TestLoginCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, LoginCmd.Short)
}

func TestLoginCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, LoginCmd.Long)
}

func TestLoginCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, LoginCmd.Run)
}

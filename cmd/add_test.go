package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, AddCmd.Use)
}

func TestAddCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, AddCmd.Short)
}

func TestAddCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, AddCmd.Long)
}

func TestAddCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, AddCmd.Run)
}

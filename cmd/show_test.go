package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, ShowCmd.Use)
}

func TestShowCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, ShowCmd.Short)
}

func TestShowCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, ShowCmd.Long)
}

func TestShowCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, ShowCmd.Run)
}

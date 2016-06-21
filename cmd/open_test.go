package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, OpenCmd.Use)
}

func TestOpenCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, OpenCmd.Short)
}

func TestOpenCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, OpenCmd.Long)
}

func TestOpenCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, OpenCmd.Run)
}

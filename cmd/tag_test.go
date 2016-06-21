package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, TagCmd.Use)
}

func TestTagCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, TagCmd.Short)
}

func TestTagCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, TagCmd.Long)
}

func TestTagCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, TagCmd.Run)
}

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, VersionCmd.Use)
}

func TestVersionCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, VersionCmd.Short)
}

func TestVersionCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, VersionCmd.Long)
}

func TestVersionCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, VersionCmd.Run)
}

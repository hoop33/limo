package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenameCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, RenameCmd.Use)
}

func TestRenameCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, RenameCmd.Short)
}

func TestRenameCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, RenameCmd.Long)
}

func TestRenameCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, RenameCmd.Run)
}

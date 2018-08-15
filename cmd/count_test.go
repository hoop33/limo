package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, CountCmd.Use)
}

func TestCountCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, CountCmd.Short)
}

func TestCountCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, CountCmd.Long)
}

func TestCountCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, CountCmd.Run)
}

func TestCountCmdHasAliasC(t *testing.T) {
	assert.Equal(t, "c", CountCmd.Aliases[0])
}

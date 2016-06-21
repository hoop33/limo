package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, ListCmd.Use)
}

func TestListCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, ListCmd.Short)
}

func TestListCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, ListCmd.Long)
}

func TestListCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, ListCmd.Run)
}

func TestListCmdHasAliasLs(t *testing.T) {
	assert.Equal(t, "ls", ListCmd.Aliases[0])
}

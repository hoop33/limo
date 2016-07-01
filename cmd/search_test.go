package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, SearchCmd.Use)
}

func TestSearchCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, SearchCmd.Short)
}

func TestSearchCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, SearchCmd.Long)
}

func TestSearchCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, SearchCmd.Run)
}

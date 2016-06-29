package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUntagCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, UntagCmd.Use)
}

func TestUntagCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, UntagCmd.Short)
}

func TestUntagCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, UntagCmd.Long)
}

func TestUntagCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, UntagCmd.Run)
}

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteCmdHasUse(t *testing.T) {
	assert.NotEmpty(t, DeleteCmd.Use)
}

func TestDeleteCmdHasShort(t *testing.T) {
	assert.NotEmpty(t, DeleteCmd.Short)
}

func TestDeleteCmdHasLong(t *testing.T) {
	assert.NotEmpty(t, DeleteCmd.Long)
}

func TestDeleteCmdHasRun(t *testing.T) {
	assert.NotEmpty(t, DeleteCmd.Run)
}

package model

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var indexPath = "./test.idx"

func TestInitIndexShouldReturnNonNilIndexWhenDoesNotExist(t *testing.T) {
	rmIndex()

	index, err := InitIndex(indexPath)
	assert.Nil(t, err)
	assert.NotNil(t, index)

	rmIndex()
}

func rmIndex() {
	if err := os.RemoveAll(indexPath); err != nil {
		panic(err)
	}
}

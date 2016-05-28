package output

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextDoesRegisterItself(t *testing.T) {
	assert.Equal(t, "output.Text", reflect.TypeOf(ForName("text")).String())
}

package output

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForNameReturnsTextWhenNameNotFound(t *testing.T) {
	output := ForName("foo")
	assert.Equal(t, "output.Text", reflect.TypeOf(output).String())
}

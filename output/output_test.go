package output

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForNameReturnsTextWhenNameNotFound(t *testing.T) {
	assert.Equal(t, "*output.Text", reflect.TypeOf(ForName("foo")).String())
}

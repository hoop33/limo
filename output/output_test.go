package output

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForNameReturnsColorWhenNameNotFound(t *testing.T) {
	assert.Equal(t, "output.Color", reflect.TypeOf(ForName("foo")).String())
}

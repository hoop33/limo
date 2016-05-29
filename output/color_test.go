package output

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorDoesRegisterItself(t *testing.T) {
	assert.Equal(t, "*output.Color", reflect.TypeOf(ForName("color")).String())
}

package encoding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalStr(t *testing.T) {
	var str string
	foo := "foo"
	b := []byte(foo)

	// Test normal byte slice
	l := UnmarshalStr(b, &str, len(b))
	assert.Equal(t, foo, str)
	assert.Equal(t, len(b), l)

	// Test byte slice with trailing zero
	b = append(b, 0)
	b = append(b, 0)
	assert.Equal(t, len(b), 5)

	l = UnmarshalStr(b, &str, len(b))
	assert.Equal(t, foo, str)
	assert.Equal(t, len(b), l)

	// Test empty byte slice
	b = make([]byte, 0)
	l = UnmarshalStr(b, &str, len(b))
	assert.Equal(t, "", str)
	assert.Equal(t, len(b), l)
}

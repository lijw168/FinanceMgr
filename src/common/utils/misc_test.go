package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompareStringSlice(t *testing.T) {
	A := []string{"hello", "world"}
	B := []string{"world", "peace"}

	onlyInA, onlyInB := CompareStringSlice(A, B)
	assert.Equal(t, onlyInA, []string{"hello"})
	assert.Equal(t, onlyInB, []string{"peace"})

	A = []string{"golang", "python"}
	B = []string{"c", "c++"}
	onlyInA, onlyInB = CompareStringSlice(A, B)
	assert.Equal(t, onlyInA, A)
	assert.Equal(t, onlyInB, B)
}

func TestGetShiftBits(t *testing.T) {
	assert.Equal(t, 16, GetShiftBits(64*1024))
}

func TestIsPowerOfTwo(t *testing.T) {
	assert.True(t, PowerOfTwo(2))
	assert.True(t, PowerOfTwo(4))
	assert.True(t, PowerOfTwo(8))
	assert.True(t, PowerOfTwo(16))
	assert.True(t, PowerOfTwo(32))

	assert.False(t, PowerOfTwo(24))
	assert.False(t, PowerOfTwo(33))
}

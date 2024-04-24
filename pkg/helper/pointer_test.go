package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPointer(t *testing.T) {
	pointer := ToPointer("test")

	assert.NotNil(t, pointer)
}

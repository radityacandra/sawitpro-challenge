package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupPassword(t *testing.T) {
	user := NewUser("Test", "+6281123456789")

	user.SetupPassword("somepassword")

	assert.NotEqual(t, "somepassword", user.Password)
	assert.NotEmpty(t, user.Password)
}

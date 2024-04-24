package model

import (
	"testing"

	"github.com/SawitProRecruitment/UserService/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestInstantiateUser(t *testing.T) {
	user := NewUser("Test", "+6281123456789")

	assert.Equal(t, "Test", user.FullName)
	assert.Equal(t, "6281123456789", user.PhoneNumber)
}

func TestUpdateUser(t *testing.T) {
	user := NewUser("Test", "+6281123456789")

	user.UpdateFullName(helper.ToPointer("Test Updated"))
	assert.Equal(t, "Test Updated", user.FullName)
	user.UpdateFullName(nil)
	assert.Equal(t, "Test Updated", user.FullName)

	user.UpdatePhoneNumber(helper.ToPointer("+6281123456788"))
	assert.Equal(t, "6281123456788", user.PhoneNumber)
	user.UpdatePhoneNumber(nil)
	assert.Equal(t, "6281123456788", user.PhoneNumber)
}

func TestPhoneNumberChanged(t *testing.T) {
	user := NewUser("Test", "+6281123456789")

	t.Run("should return false if phone number doesn't change", func(t *testing.T) {
		result := user.PhoneNumberChanged(helper.ToPointer("+6281123456789"))
		assert.False(t, result)
	})

	t.Run("should return false if phone number is changed", func(t *testing.T) {
		result := user.PhoneNumberChanged(helper.ToPointer("+6281123456788"))
		assert.True(t, result)
	})
}

func TestGetFormalizedPhoneNumber(t *testing.T) {
	user := NewUser("Test", "+6281123456789")
	assert.Equal(t, "+6281123456789", user.GetFormalizedPhoneNumber())
}

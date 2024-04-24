package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUser(t *testing.T) {
	type testCase struct {
		name       string
		password   string
		returnBool bool
	}

	testCases := []testCase{
		{
			name:       "should return false if password doesn't match",
			password:   "invalidpassword",
			returnBool: false,
		},
		{
			name:       "should return true if password matched",
			password:   "somepassword",
			returnBool: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			user := NewUser("Test", "+6281123456789")
			user.SetupPassword("somepassword")

			result := user.Authenticate(testCase.password)
			assert.Equal(t, testCase.returnBool, result)
		})
	}
}

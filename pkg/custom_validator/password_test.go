package custom_validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	type TestStruct struct {
		Password string `validate:"password"`
	}

	type testCase struct {
		TestStruct
		name  string
		error bool
	}

	testCases := []testCase{
		{
			TestStruct: TestStruct{
				Password: "passwordwithoutnum",
			},
			name:  "should return error if no number, uppercase, symbol provided",
			error: true,
		},
		{
			TestStruct: TestStruct{
				Password: "passwordwithoutuppercase111",
			},
			name:  "should return error if no uppercase, symbol provided",
			error: true,
		},
		{
			TestStruct: TestStruct{
				Password: "PasswordWithoutSymbol111",
			},
			name:  "should return error if no symbol provided",
			error: true,
		},
		{
			TestStruct: TestStruct{
				Password: "Val!dPassword111!",
			},
			name:  "should return success",
			error: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			validator := validator.New(validator.WithRequiredStructEnabled())
			validator.RegisterValidation("password", ValidatePassword)

			err := validator.Struct(testCase.TestStruct)
			if testCase.error {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

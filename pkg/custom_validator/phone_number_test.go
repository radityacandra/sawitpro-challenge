package custom_validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidatePhoneNumber(t *testing.T) {
	type TestStruct struct {
		PhoneNumber string `validate:"phoneNumber"`
	}

	type testCase struct {
		TestStruct
		name  string
		error bool
	}

	testCases := []testCase{
		{
			TestStruct: TestStruct{
				PhoneNumber: "081123456789",
			},
			name:  "should return error if provided original phonenumber",
			error: true,
		},
		{
			TestStruct: TestStruct{
				PhoneNumber: "6281123456789",
			},
			name:  "should return error if not prefixed with +62",
			error: true,
		},
		{
			TestStruct: TestStruct{
				PhoneNumber: "0811+6223456789",
			},
			name:  "should return error if +62 is inserted at the middle of the value",
			error: true,
		},
		{
			TestStruct: TestStruct{
				PhoneNumber: "+6281123456789",
			},
			name:  "should return success",
			error: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			validator := validator.New(validator.WithRequiredStructEnabled())
			validator.RegisterValidation("phoneNumber", ValidatePhoneNumber)

			err := validator.Struct(testCase.TestStruct)
			if testCase.error {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

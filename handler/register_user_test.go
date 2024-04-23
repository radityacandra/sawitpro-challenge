package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	type testCase struct {
		name                   string
		reqBody                string
		getUserByPhoneNumber   *model.User
		expectSavedPhoneNumber string
		expectSavedUserName    string
		expectHttpCode         int
	}

	testCases := []testCase{
		{
			name:                   "should return success",
			reqBody:                `{"fullName":"test","phoneNumber":"+6281123456789","password":"Password123!!!"}`,
			expectSavedPhoneNumber: "6281123456789",
			expectSavedUserName:    "test",
			expectHttpCode:         201,
			getUserByPhoneNumber:   nil,
		},
		{
			name:                   "should return error if failed to bind data",
			reqBody:                `{"fullName":123}`,
			expectSavedPhoneNumber: "6281123456789",
			expectSavedUserName:    "test",
			expectHttpCode:         422,
			getUserByPhoneNumber:   nil,
		},
		{
			name:                   "should return validation error",
			reqBody:                `{"fullName":"test","phoneNumber":"6281123456789","password":"Password123!!!"}`,
			expectSavedPhoneNumber: "6281123456789",
			expectSavedUserName:    "test",
			expectHttpCode:         400,
			getUserByPhoneNumber:   nil,
		},
		{
			name:                   "should return error if there is conflicting phone number",
			reqBody:                `{"fullName":"test","phoneNumber":"+6281123456789","password":"Password123!!!"}`,
			expectSavedPhoneNumber: "6281123456789",
			expectSavedUserName:    "test",
			expectHttpCode:         409,
			getUserByPhoneNumber:   model.NewUser("test", "6281123456789"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var reqBody generated.CreateUserRequest
			json.Unmarshal([]byte(testCase.reqBody), &reqBody)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := repository.NewMockRepositoryInterface(ctrl)

			m.EXPECT().GetUserByPhoneNumber(gomock.Any(), reqBody.PhoneNumber).Return(testCase.getUserByPhoneNumber).AnyTimes()

			m.EXPECT().InsertUser(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, user *model.User) (*model.User, error) {
				assert.Equal(t, testCase.expectSavedUserName, user.FullName)
				assert.Equal(t, testCase.expectSavedPhoneNumber, user.PhoneNumber)
				assert.NotEqual(t, reqBody.Password, user.Password)

				return model.NewUser(reqBody.FullName, reqBody.PhoneNumber), nil
			}).AnyTimes()

			server := NewServer(NewServerOptions{
				Repository: m,
			})

			e := echo.New()
			generated.RegisterHandlers(e, server)

			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(testCase.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			assert.NoError(t, server.RegisterUser(c))
			assert.Equal(t, testCase.expectHttpCode, rec.Code)
		})
	}
}

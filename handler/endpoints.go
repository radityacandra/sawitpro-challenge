package handler

import (
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"

	"github.com/SawitProRecruitment/UserService/pkg/helper"
	"github.com/SawitProRecruitment/UserService/pkg/jwt"
	"github.com/SawitProRecruitment/UserService/pkg/response_wrapper"
	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterUser(ctx echo.Context) error {
	var req generated.CreateUserRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, response_wrapper.WrapperError(err, 422))
	}

	if err := s.Validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(err, 400))
	}

	// check for existing phone number
	if existingUser := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber); existingUser != nil {
		return ctx.JSON(http.StatusConflict,
			response_wrapper.WrapperError(errors.New("registered user with the same phone number already exist"), 400))
	}

	user := model.NewUser(req.FullName, req.PhoneNumber)
	if _, err := s.Repository.InsertUser(ctx.Request().Context(), user); err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(err, 400))
	}

	if err := user.SetupPassword(req.Password); err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(err, 400))
	}

	return ctx.JSON(http.StatusCreated, response_wrapper.WrapperSuccess(generated.CreateUserResponse{
		FullName:    user.FullName,
		Id:          user.Id,
		PhoneNumber: user.GetFormalizedPhoneNumber(),
	}))
}

func (s *Server) AuthenticateUser(ctx echo.Context) error {
	var req generated.AuthenticateUserRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, response_wrapper.WrapperError(err, 422))
	}

	if err := s.Validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(err, 400))
	}

	// find user by phone number
	user := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber)
	if user == nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(errors.New("invalid phone number or password"), 400))
	}

	// compare hash with provided password
	if !user.Authenticate(req.Password) {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(errors.New("invalid phone number or password"), 400))
	}

	// generate jwt
	token, expiredAt, err := jwt.BuildToken(map[string]interface{}{
		"userId":   user.Id,
		"fullName": user.FullName,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response_wrapper.WrapperError(errors.New("failed to build token. please try again"), 500))
	}

	return ctx.JSON(http.StatusOK, response_wrapper.WrapperSuccess(generated.AuthenticateUserResponse{
		AccessToken: &token,
		ExpiredAt:   &expiredAt,
	}))
}

func (s *Server) GetUserProfile(ctx echo.Context) error {
	data, err := jwt.AuthorizeToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusForbidden, response_wrapper.WrapperError(err, 403))
	}

	user := s.Repository.GetUserById(ctx.Request().Context(), int(data["userId"].(float64)))

	return ctx.JSON(http.StatusOK, response_wrapper.WrapperSuccess(generated.UserProfileDto{
		FullName:    &user.FullName,
		PhoneNumber: helper.ToPointer(user.GetFormalizedPhoneNumber()),
	}))
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	data, err := jwt.AuthorizeToken(ctx.Request().Header.Get("Authorization"))
	if err != nil {
		return ctx.JSON(http.StatusForbidden, response_wrapper.WrapperError(err, 403))
	}

	var req generated.UserProfileDto

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, response_wrapper.WrapperError(err, 422))
	}

	if err := s.Validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(err, 400))
	}

	user := s.Repository.GetUserById(ctx.Request().Context(), int(data["userId"].(float64)))

	// check for conflicting possibility if user changed the phone number
	if user.PhoneNumberChanged(req.PhoneNumber) {
		existingUser := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), *req.PhoneNumber)
		if existingUser != nil {
			return ctx.JSON(http.StatusConflict,
				response_wrapper.WrapperError(errors.New("registered user with the same phone number already exist"), 400))
		}
	}

	user.UpdateFullName(req.FullName)
	user.UpdatePhoneNumber(req.PhoneNumber)

	if err := s.Repository.UpdateUser(ctx.Request().Context(), user); err != nil {
		return ctx.JSON(http.StatusInternalServerError, response_wrapper.WrapperError(errors.New("failed to update profile"), 500))
	}

	return ctx.JSON(http.StatusAccepted, response_wrapper.WrapperSuccess(generated.UserProfileDto{
		FullName:    &user.FullName,
		PhoneNumber: helper.ToPointer(user.GetFormalizedPhoneNumber()),
	}))
}

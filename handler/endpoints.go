package handler

import (
	"errors"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"

	"github.com/SawitProRecruitment/UserService/pkg/jwt"
	"github.com/SawitProRecruitment/UserService/pkg/response_wrapper"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) RegisterUser(ctx echo.Context) error {
	var req generated.CreateUserRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, response_wrapper.WrapperError(err, 422))
	}

	if err := s.Validator.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(err, 400))
	}

	// TODO: move this to struct method
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(err, 400))
	}

	// check for existing phone number
	if existingUser := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), req.PhoneNumber); existingUser != nil {
		return ctx.JSON(http.StatusBadRequest,
			response_wrapper.WrapperError(errors.New("registered user with the same phone number already exist"), 400))
	}

	user := model.NewUser(req.FullName, string(hash), req.PhoneNumber)
	if _, err := s.Repository.InsertUser(ctx.Request().Context(), user); err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(err, 400))
	}

	return ctx.JSON(http.StatusOK, response_wrapper.WrapperSuccess(user))
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
	token, err := jwt.BuildToken(map[string]interface{}{
		"userId":   user.Id,
		"fullName": user.FullName,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response_wrapper.WrapperError(errors.New("failed to build token. please try again"), 400))
	}

	return ctx.JSON(http.StatusOK, response_wrapper.WrapperSuccess(token))
}

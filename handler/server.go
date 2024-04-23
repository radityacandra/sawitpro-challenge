package handler

import (
	"github.com/SawitProRecruitment/UserService/pkg/custom_validator"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Repository repository.RepositoryInterface
	Validator  *validator.Validate
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
}

func NewServer(opts NewServerOptions) *Server {
	validator := validator.New(validator.WithRequiredStructEnabled())
	validator.RegisterValidation("password", custom_validator.ValidatePassword)
	validator.RegisterValidation("phoneNumber", custom_validator.ValidatePhoneNumber)

	return &Server{
		Repository: opts.Repository,
		Validator:  validator,
	}
}

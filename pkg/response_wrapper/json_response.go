package response_wrapper

import "github.com/go-playground/validator/v10"

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Errors  []Error     `json:"errors"`
}

func Wrapper(data interface{}, err error, code int) Response {
	if err == nil {
		return WrapperSuccess(data)
	}

	return WrapperError(err, code)
}

func WrapperSuccess(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
		Code:    200,
	}
}

func WrapperError(err error, code int) Response {
	response := Response{
		Success: false,
		Code:    code,
		Message: err.Error(),
	}

	if err, ok := err.(validator.ValidationErrors); ok {
		for _, validationError := range err {
			response.Errors = append(response.Errors, Error{
				Field:   validationError.Field(),
				Message: validationError.Error(),
			})
		}
	}

	return response
}

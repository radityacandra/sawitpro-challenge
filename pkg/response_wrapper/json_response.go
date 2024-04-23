package response_wrapper

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
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
	return Response{
		Success: false,
		Code:    code,
		Message: err.Error(),
	}
}

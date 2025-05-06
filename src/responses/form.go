package responses

import "twitter/src/validations"

type Response struct {
	Status          bool                           `json:"status"`
	StatusCode      int                            `json:"statusCode"`
	Message         string                         `json:"message,omitempty"`
	Result          interface{}                    `json:"result,omitempty"`
	Error           string                         `json:"error,omitempty"`
	ValidationError *[]validations.ValidationError `json:"validationError,omitempty"`
}

func GenerateNormalResponse(statusCode int, result interface{}, msg string) *Response {
	return &Response{
		Status: true, StatusCode: statusCode, Result: result, Message: msg,
	}
}

func GenerateResponseWithError(statusCode int, err error, msg string) *Response {
	return &Response{
		Status: false, StatusCode: statusCode, Error: err.Error(), Message: msg,
	}
}

func GenerateResponseWithValidationError(statusCode int, err error, msg string) *Response {
	return &Response{
		StatusCode: statusCode, ValidationError: validations.GetValidationError(err), Message: msg,
	}
}

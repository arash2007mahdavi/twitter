package responses

import (
	"twitter/src/logger"
	"twitter/src/validations"
)

type Response struct {
	Status          bool                           `json:"status"`
	StatusCode      int                            `json:"statusCode"`
	Message         string                         `json:"message,omitempty"`
	Result          interface{}                    `json:"result,omitempty"`
	Error           string                         `json:"error,omitempty"`
	ValidationError *[]validations.ValidationError `json:"validationError,omitempty"`
}

var log = logger.NewLogger()

func GenerateNormalResponse(statusCode int, result interface{}, msg string) *Response {
	return &Response{
		Status: true, StatusCode: statusCode, Result: result, Message: msg,
	}
}

func GenerateResponseWithError(statusCode int, err error, msg string) *Response {
	log.Error(logger.Error, logger.SimpleError, msg, map[logger.ExtraCategory]interface{}{logger.StatusCode: statusCode, logger.Err: err})
	return &Response{
		Status: false, StatusCode: statusCode, Error: err.Error(), Message: msg,
	}
}

func GenerateResponseWithValidationError(statusCode int, err error, msg string) *Response {
	log.Error(logger.Error, logger.ValidationError, msg, map[logger.ExtraCategory]interface{}{logger.StatusCode: statusCode, logger.Err: err})
	return &Response{
		StatusCode: statusCode, ValidationError: validations.GetValidationError(err), Message: msg,
	}
}

package utils

import "fmt"

// Custom Errors
type ErrorType string

const (
	ErrorType_TokenExpired      ErrorType = "TokenExpired"
	ErrorType_JWTError          ErrorType = "JWTError"
	ErrorType_QueryError        ErrorType = "QueryError"
	ErrorType_FormatError       ErrorType = "FormatError"
	ErrorType_UpdateError       ErrorType = "UpdateError"
	ErrorType_UserAlreadyExists ErrorType = "UserAlreadyExists"
	ErrorType_UserDoesNotExist  ErrorType = "UserDoesNotExist"
	ErrorType_UserLoginFailed   ErrorType = "UserLoginFailed"
)

type ApiError struct {
	Type    ErrorType
	Message string
}

func (apiErr ApiError) Error() string {
	return fmt.Sprintf("[%s] %s", string(apiErr.Type), apiErr.Message)
}

func NewApiError(errType ErrorType, message string) ApiError {
	return ApiError{Type: errType, Message: message}
}

func IsErrorOfType(err error, errType ErrorType) bool {
	apiErr, ok := err.(ApiError)
	if ok {
		return apiErr.Type == errType
	}
	return false
}

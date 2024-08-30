package xweb

import (
	"net/http"
)

type CodedError interface {
	HTTPCode() int
	CharCode() string
	Details() map[string]any
	SetDetail(key string, value any)
	RemoveDetail(key string)
	Error() string
}

var _ CodedError = &GenericCodedError{}

const (
	InternalErrorCode = "INTERNAL_ERROR"
	NotFoundErrorCode = "NOT_FOUND"
)

type GenericCodedError struct {
	msg      string
	code     int
	charCode string
	details  map[string]any
}

func (e *GenericCodedError) Error() string {
	return e.msg
}

func (e *GenericCodedError) HTTPCode() int {
	return e.code
}

func (e *GenericCodedError) CharCode() string {
	return e.charCode
}

func (e *GenericCodedError) Details() map[string]any {
	return e.details
}

func (e *GenericCodedError) SetDetail(key string, value any) {
	e.details[key] = value
}

func (e *GenericCodedError) RemoveDetail(key string) {
	delete(e.details, key)
}

func NewGenericCodedError(err error, httpCode int, charCode string, details map[string]any) *GenericCodedError {
	if details == nil {
		details = map[string]any{}
	}
	return &GenericCodedError{
		msg:      err.Error(),
		code:     httpCode,
		charCode: charCode,
		details:  details,
	}
}

func NewInternalError(err error) CodedError {
	return NewGenericCodedError(err, http.StatusInternalServerError, InternalErrorCode, nil)
}

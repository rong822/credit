package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrDataAlreadyExist = errors.New("data already exist")
)

const (
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

type errorFabric struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func New(code int, message string) error {
	return &errorFabric{
		Code:    code,
		Message: message,
	}
}

func (e *errorFabric) Error() string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`, e.Code, e.Message)
}

func Errorf(code int, format string, args ...interface{}) error {
	return &errorFabric{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func Error(code int, args ...interface{}) error {
	return &errorFabric{
		Code:    code,
		Message: fmt.Sprint(args...),
	}
}

func Unmarshal(err interface{}) *errorFabric {
	e := new(errorFabric)

	var data []byte
	switch err := err.(type) {
	case error:
		data = []byte(err.Error())
	case string:
		data = []byte(err)
	case []byte:
		data = err
	default:
		return e
	}

	json.Unmarshal(data, e)
	return e
}

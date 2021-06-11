package httperr

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

var (
	ErrBaseNil = errors.New("httperr: base error is nil")
)

// Error implements built in error with status code, caller, and message attribute
// it is recommended to create Error from the provided method instead of creating
// right from the struct to fill the caller.
type Error struct {
	Err        error
	statusCode int
	message    string
	caller     string
}

// New wraps err with defined statusCode and message
func New(err error, statusCode int, message string) error {
	return &Error{
		Err:        getValidErr(err),
		statusCode: getValidStatusCode(statusCode),
		message:    message,
		caller:     getCaller(),
	}
}

// Wrap wraps err with custom message
// Wrap's result inherit statusCode from err if err equals *Error
func Wrap(err error, msg string) error {
	e := getValidErr(err)
	var statusCode int

	switch e.(type) {
	case Error:
		statusCode = e.(Error).statusCode
	case *Error:
		statusCode = e.(*Error).statusCode
	}

	return &Error{
		Err:        e,
		statusCode: getValidStatusCode(statusCode),
		message:    msg,
		caller:     getCaller(),
	}
}

// From creates new error from defined statusCode
// if statusCode doesn't have any status text
// statusCode changed to http.StatusInternalServerError
func From(statusCode int) error {
	sc := getValidStatusCode(statusCode)
	text := http.StatusText(sc)

	return &Error{
		Err:        errors.New(text),
		statusCode: sc,
		message:    text,
		caller:     getCaller(),
	}
}

// Error returns error message with caller
func (e Error) Error() string {
	switch e.Unwrap().(type) {
	case Error, *Error:
		return fmt.Sprintf(`{"cause": %v, "message": "%v", "caller": "%v"}`,
			e.Unwrap().Error(), strings.ReplaceAll(e.message, `"`, "`"), e.caller)
	}
	return fmt.Sprintf(`{"error": "%v", "message": "%v", "caller": "%v"}`,
		strings.ReplaceAll(e.Unwrap().Error(), `"`, "`"), strings.ReplaceAll(e.message, `"`, "`"), e.caller)
}

// Unwrap enables errors.As and errors.Is
func (e Error) Unwrap() error {
	return getValidErr(e.Err)
}

// StatusCode returns e.statusCode
func (e Error) StatusCode() int {
	return getValidStatusCode(e.statusCode)
}

// getCaller uses log.Lshortfile to format the caller
func getCaller() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return fmt.Sprintf("%s:%d", file, line)
}

// getValidStatusCode returns http.StatusInternalServerError if the given statusCode is not valid
func getValidStatusCode(statusCode int) int {
	t := http.StatusText(statusCode)
	if t == "" {
		return http.StatusInternalServerError
	}

	return statusCode
}

// getValidErr returns errors.New("httperr: base error is nil") if the given err is nil
func getValidErr(err error) error {
	if err != nil {
		return err
	}

	return ErrBaseNil
}

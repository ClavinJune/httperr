package httperr

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	tt := []struct {
		base    error
		sc      int
		message string
		want    *Error
	}{{
		base:    nil,
		sc:      0,
		message: "",
		want: &Error{
			Err:        ErrBaseNil,
			statusCode: http.StatusInternalServerError,
			message:    "",
		},
	}, {
		base:    sql.ErrNoRows,
		sc:      http.StatusNotFound,
		message: "user not found",
		want: &Error{
			Err:        sql.ErrNoRows,
			statusCode: http.StatusNotFound,
			message:    "user not found",
		},
	}, {
		base: &Error{
			Err:        &Error{},
			statusCode: http.StatusNotFound,
		},
		sc: http.StatusInternalServerError,
		want: &Error{
			Err: &Error{
				Err:        &Error{},
				statusCode: http.StatusNotFound,
			},
			statusCode: http.StatusInternalServerError,
		},
	}}

	for _, test := range tt {
		got := New(test.base, test.sc, test.message).(*Error)

		if got.StatusCode() != test.want.StatusCode() {
			t.Fatalf(`got "%v", want: "%v"`, got.StatusCode(), test.want.StatusCode())
		}

		if got.message != test.want.message {
			t.Fatalf(`got "%v", want: "%v"`, got.message, test.want.message)
		}
		if got.Err.Error() != test.want.Err.Error() {
			t.Fatalf(`got "%v", want: "%v"`, got.Err, test.want.Err)
		}
	}
}

func TestWrap(t *testing.T) {
	base1 := From(http.StatusMethodNotAllowed)

	tt := []struct {
		base error
		msg  string
		want *Error
	}{{
		base: nil,
		msg:  "",
		want: &Error{
			Err:        ErrBaseNil,
			statusCode: http.StatusInternalServerError,
			message:    "",
		},
	}, {
		base: base1,
		msg:  "http status method not allowed",
		want: &Error{
			Err:        base1,
			statusCode: http.StatusMethodNotAllowed,
			message:    "http status method not allowed",
		},
	}, {
		base: Error{},
		msg:  "http status method not allowed",
		want: &Error{
			Err:        Error{},
			statusCode: http.StatusInternalServerError,
			message:    "http status method not allowed",
		},
	}}

	for _, test := range tt {
		got := Wrap(test.base, test.msg).(*Error)

		if got.StatusCode() != test.want.StatusCode() {
			t.Fatalf(`got "%v", want: "%v"`, got.StatusCode(), test.want.StatusCode())
		}

		if got.message != test.want.message {
			t.Fatalf(`got "%v", want: "%v"`, got.message, test.want.message)
		}

		if got.Err.Error() != test.want.Err.Error() {
			t.Fatalf(`got "%v", want: "%v"`, got.Err, test.want.Err)
		}
	}
}

func TestFrom(t *testing.T) {
	tt := []struct {
		statusCode int
		want       *Error
	}{{
		statusCode: 0,
		want: &Error{
			Err:        errors.New(http.StatusText(http.StatusInternalServerError)),
			statusCode: http.StatusInternalServerError,
			message:    http.StatusText(http.StatusInternalServerError),
		},
	}, {
		statusCode: http.StatusNotFound,
		want: &Error{
			Err:        errors.New(http.StatusText(http.StatusNotFound)),
			statusCode: http.StatusNotFound,
			message:    http.StatusText(http.StatusNotFound),
		},
	}, {
		statusCode: 11111111,
		want: &Error{
			Err:        errors.New(http.StatusText(http.StatusInternalServerError)),
			statusCode: http.StatusInternalServerError,
			message:    http.StatusText(http.StatusInternalServerError),
		},
	}, {
		statusCode: http.StatusMethodNotAllowed,
		want: &Error{
			Err:        errors.New(http.StatusText(http.StatusMethodNotAllowed)),
			statusCode: http.StatusMethodNotAllowed,
			message:    http.StatusText(http.StatusMethodNotAllowed),
		},
	}}

	for _, test := range tt {
		got := From(test.statusCode).(*Error)

		if got.StatusCode() != test.want.StatusCode() {
			t.Fatalf(`got "%v", want: "%v"`, got.StatusCode(), test.want.StatusCode())
		}

		if got.message != test.want.message {
			t.Fatalf(`got "%v", want: "%v"`, got.message, test.want.message)
		}

		if got.Err.Error() != test.want.Err.Error() {
			t.Fatalf(`got "%v", want: "%v"`, got.Err, test.want.Err)
		}
	}
}

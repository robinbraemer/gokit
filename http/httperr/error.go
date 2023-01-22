package httperr

import (
	"errors"
	"fmt"
	"net/http"
)

// HandlerFunc is a function that handles an http request and returns an error.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorHandler is a function that handles an error and can be used as the root http handler for the server.
func ErrorHandler(h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { WriteError(w, h(w, r)) })
}

// Error is an error with an http status code.
func Error(code int, format string) error {
	return Errorf(code, format)
}

// Errorf is an error with an http status code and formatted message.
func Errorf(code int, format string, a ...interface{}) error {
	if code == 0 {
		return nil
	}
	return &httpErr{
		code: code,
		err:  fmt.Errorf(format, a...),
	}
}

// ErrorCode returns an error with the given http status code.
func ErrorCode(code int) error {
	if code == 0 {
		return nil
	}
	return &httpErr{
		code: code,
		err:  errors.New(http.StatusText(code)),
	}
}

// WriteError writes an error to the response writer if the error exists.
func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	var hErr *httpErr
	if errors.As(err, &hErr) {
		http.Error(w, hErr.err.Error(), hErr.code)
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

type httpErr struct {
	err  error
	code int
}

func (e *httpErr) Error() string { return fmt.Sprintf("[%d] %v", e.code, e.err) }
func (e *httpErr) Unwrap() error { return e.err }

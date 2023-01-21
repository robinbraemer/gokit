package httputil

import (
	"errors"
	"fmt"
	"net/http"
)

type httpErr struct {
	Err  error
	Code int
}

func (e *httpErr) Error() string { return fmt.Sprintf("[%d] %v", e.Code, e.Err) }
func (e *httpErr) Unwrap() error { return e.Err }

func Error(code int, format string) error {
	return Errorf(code, format)
}

func Errorf(code int, format string, a ...interface{}) error {
	if code == 0 {
		return nil
	}
	return &httpErr{
		Code: code,
		Err:  fmt.Errorf(format, a...),
	}
}

func ErrorCode(code int) error {
	if code == 0 {
		return nil
	}
	return &httpErr{
		Code: code,
		Err:  errors.New(http.StatusText(code)),
	}
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(h HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { WriteError(w, h(w, r)) })
}

func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	var hErr *httpErr
	if errors.As(err, &hErr) {
		http.Error(w, hErr.Err.Error(), hErr.Code)
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

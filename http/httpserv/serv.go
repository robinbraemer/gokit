package httpserv

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

// NewDefaultServer returns a new http.Server with the given handler and default settings.
func NewDefaultServer(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 30,
	}
}

// Server is an interface that specified partial stdlib http.Server methods.
type Server interface {
	io.Closer
	Shutdown(context.Context) error
	ListenAndServe() error
	Serve(net.Listener) error
}

// DefaultShutdownTimeout is the default timeout for shutting down the server by Serve.
var DefaultShutdownTimeout = time.Second * 25

// ListenAndServe is a convenience function that calls Serve.
func ListenAndServe(ctx context.Context, svr Server) error { return Serve(ctx, nil, svr) }

// Serve starts the server and waits for the context to be canceled to begin graceful shutdown.
func Serve(ctx context.Context, ln net.Listener, svr Server) error {
	defer svr.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		stopCtx, stopCancel := context.WithTimeout(ctx, DefaultShutdownTimeout)
		defer stopCancel()
		_ = svr.Shutdown(stopCtx)
	}()

	var err error
	if ln == nil {
		err = svr.ListenAndServe()
	} else {
		err = svr.Serve(ln)
	}
	if errors.Is(err, http.ErrServerClosed) || errors.Is(err, net.ErrClosed) /* underlying listener closed */ {
		err = nil
	}
	return err
}

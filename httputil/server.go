package httputil

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

var DefaultShutdownTimeout = time.Second * 25

type Server interface {
	io.Closer
	Shutdown(context.Context) error
	ListenAndServe() error
	Serve(net.Listener) error
}

func ListenAndServe(ctx context.Context, svr Server) error { return Serve(ctx, nil, svr) }

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

func NewDefaultServer(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 10,
		WriteTimeout:      time.Second * 30,
	}
}

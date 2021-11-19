package gap

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type options struct {
	ctx     context.Context
	servers []*Server
}

type Option func(*options)

type App struct {
	opts   options
	ctx    context.Context
	cancel func()
}

// New App
func NewApp(opts ...Option) *App {
	options := options{
		ctx: context.Background(),
	}

	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithCancel(options.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   options,
	}
}

// Servers 添加Server
func Servers(srv ...*Server) Option {
	return func(o *options) { o.servers = srv }
}

// Stop gracefully stops the application.
func (a *App) Stop() error {
	//todo

	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func (app *App) Run() error {
	ctx, cancel := context.WithCancel(app.ctx)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}

	for _, srv := range app.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // wait for stop signal
			return srv.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start()
		})
	}

	wg.Wait()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				app.Stop()
			}
		}
	})
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

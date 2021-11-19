package gap

import (
	"context"
	"errors"
	"log"
	"net/http"
)

type Server struct {
	name        string
	HttpServer  *http.Server
	HttpHandler *HttpHandler
}

// Default	Server默认实现
func Default(name string, addr string) *Server {
	srv := RawSrv(name, addr)

	srv.Use(Router(), Logger())
	return srv
}

// Default	Server默认实现
func RawSrv(name string, addr string) *Server {
	srv := &Server{
		name:        name,
		HttpHandler: NewHttpHandler(name),
		HttpServer:  &http.Server{Addr: addr},
	}
	return srv
}

// AddRouter	添加路由
func (s *Server) AddRouter(path string, f HandlerFunc) {
	s.HttpHandler.AddRouter(path, f)
}

// UseSimple	创建管道简单中间件
func (s *Server) UseSimple(h func(*Context)) {
	middleware := func(next MiddlewareHandler) MiddlewareHandler {
		return func(c *Context) {
			h(c)
			if next != nil {
				next(c)
			}
		}
	}
	s.HttpHandler.Use(middleware)
}

// Use	使用中间件
func (s *Server) Use(m ...Middleware) {
	s.HttpHandler.Use(m...)
}

// Start 启动Server监听
func (s *Server) Start() error {

	s.HttpHandler.BuildPipeline()
	s.HttpServer.Handler = s.HttpHandler

	log.Printf("[HTTP] server [%s] listening on: %s", s.name, s.HttpServer.Addr)
	if err := s.HttpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop
func (s *Server) Stop(ctx context.Context) error {
	err := s.HttpServer.Shutdown(ctx)
	log.Printf("Server [%s] Exited Properly", s.name)
	return err
}

// RegisterOnShutdown
func (s *Server) RegisterOnShutdown(f func()) {
	s.HttpServer.RegisterOnShutdown(f)
}

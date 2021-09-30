package gap

import (
	"log"
	"net/http"
)

type Server struct {
	Context  *Context
	pipeline []Middleware
	entry    MiddlewareHandler
}

// Default	Server默认实现
func Default() *Server {
	srv := RawSrv()

	srv.Use(Router(), Logger())
	return srv
}

// Default	Server默认实现
func RawSrv() *Server {
	srv := &Server{
		pipeline: make([]Middleware, 0, 2),
		Context: &Context{
			Router: make(map[string]HandlerFunc),
		},
	}
	return srv
}

// ServeHTTP	Http请求处理入口
func (s *Server) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	s.Context.Request, s.Context.ResponseWriter = req, &ResponseWriter{
		ResponseWriter: writer,
		StatusCode:     http.StatusOK,
	}

	s.entry(s.Context)
}

// AddRouter	添加路由
func (s *Server) AddRouter(path string, h HandlerFunc) {
	s.Context.Router[path] = h
}

// UseSimple	创建管道简单中间件
func (s *Server) UseSimple(h func(*Context)) {
	middleware := func(next MiddlewareHandler) MiddlewareHandler {
		return func(c *Context) {
			h(s.Context)
			if next != nil {
				next(s.Context)
			}
		}
	}
	s.pipeline = append(s.pipeline, middleware)
}

// Use	使用中间件
func (s *Server) Use(m ...Middleware) {
	s.pipeline = append(s.pipeline, m...)
}

// Run 启动Server监听
func (s *Server) Run(addr string) {

	// build pipeline
	for i := len(s.pipeline); i > 0; i-- {
		s.entry = s.pipeline[i-1](s.entry)
	}
	log.Fatal(http.ListenAndServe(addr, s))
}

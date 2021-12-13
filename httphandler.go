package gap

import (
	"net/http"

	logx "github.com/hzhhong/gap/log"
)

type HttpHandler struct {
	Context  *Context
	pipeline []Middleware
	entry    MiddlewareHandler
}

func NewHttpHandler(srvname string, logger logx.Logger) *HttpHandler {
	h := &HttpHandler{
		pipeline: make([]Middleware, 0, 2),
		Context: &Context{
			Router:  make(map[string]HandlerFunc),
			SrvName: srvname,
			Logger:  logger,
		},
	}
	return h
}

// ServeHTTP	Http请求处理入口
func (h *HttpHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	h.Context.Request, h.Context.ResponseWriter = req, &ResponseWriter{
		ResponseWriter: writer,
		StatusCode:     http.StatusOK,
	}

	h.entry(h.Context)
}

func (h *HttpHandler) Use(m ...Middleware) {
	h.pipeline = append(h.pipeline, m...)
}

// AddRouter	添加路由
func (h *HttpHandler) AddRouter(path string, f HandlerFunc) {
	h.Context.Router[path] = f
}

func (h *HttpHandler) BuildPipeline() {

	for i := len(h.pipeline); i > 0; i-- {
		h.entry = h.pipeline[i-1](h.entry)
	}
}

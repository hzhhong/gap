package gap

import (
	"net/http"
	"sync"

	logx "github.com/hzhhong/gap/log"
)

type HttpHandler struct {
	pipeline []Middleware
	entry    MiddlewareHandler

	SrvName string
	Router  map[string]HandlerFunc
	Logger  logx.Logger
	mu      sync.Mutex
}

func NewHttpHandler(srvname string, logger logx.Logger) *HttpHandler {
	h := &HttpHandler{
		pipeline: make([]Middleware, 0, 2),
		SrvName:  srvname,
		Logger:   logger,
		Router:   make(map[string]HandlerFunc),
	}
	return h
}

// ServeHTTP	Http请求处理入口
func (h *HttpHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	ctx := &Context{
		Router:  h.Router,
		SrvName: h.SrvName,
		Logger:  h.Logger,
		Request: req,
		ResponseWriter: &ResponseWriter{
			ResponseWriter: writer,
			StatusCode:     http.StatusOK,
		},
	}

	h.entry(ctx)
}

func (h *HttpHandler) Use(m ...Middleware) {
	h.pipeline = append(h.pipeline, m...)
}

// AddRouter	添加路由
func (h *HttpHandler) AddRouter(path string, f HandlerFunc) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Router[path] = f
}

func (h *HttpHandler) BuildPipeline() {

	for i := len(h.pipeline); i > 0; i-- {
		h.entry = h.pipeline[i-1](h.entry)
	}
}

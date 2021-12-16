package gap

import (
	"net/http"
	"sync"

	logx "github.com/hzhhong/gap/log"
)

type HandlerFunc func(*Context)

type Context struct {
	SrvName string
	Router  map[string]HandlerFunc

	ResponseWriter *ResponseWriter
	Request        *http.Request
	Logger         logx.Logger
	mu             sync.Mutex
}

// GetRouteHandler
func (ctx *Context) GetRouteHandler(path string) (h HandlerFunc, ok bool) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	h, ok = ctx.Router[path]

	return h, ok
}

// SetRouteHandler
func (ctx *Context) SetRouteHandler(path string, handler HandlerFunc) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.Router[path] = handler
}

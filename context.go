package gap

import (
	"net/http"

	logx "github.com/hzhhong/gap/log"
)

type HandlerFunc func(*Context)

type Context struct {
	SrvName string
	Router  map[string]HandlerFunc

	ResponseWriter *ResponseWriter
	Request        *http.Request
	logger         logx.Logger
}

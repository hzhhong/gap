package gap

import "net/http"

type HandlerFunc func(*Context)

type Context struct {
	Router map[string]HandlerFunc

	ResponseWriter *ResponseWriter
	Request        *http.Request
}

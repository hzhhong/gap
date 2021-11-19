package gap

import (
	"log"
)

func Logger() Middleware {
	return func(next MiddlewareHandler) MiddlewareHandler {
		return func(ctx *Context) {

			log.Printf("clientIp: %s; server: %s; path: %s; statuscode: %d\n", ctx.Request.RemoteAddr, ctx.SrvName, ctx.Request.URL.Path, ctx.ResponseWriter.StatusCode)
			if next != nil {
				next(ctx)
			}
		}
	}
}

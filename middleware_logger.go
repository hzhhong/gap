package gap

import (
	"time"

	logx "github.com/hzhhong/gap/log"
)

func LoggerProcessor() Middleware {
	return func(next MiddlewareHandler) MiddlewareHandler {
		return func(ctx *Context) {
			begintime := time.Now()
			// log.Printf("clientIp: %s; server: %s; path: %s; statuscode: %d\n", ctx.Request.RemoteAddr, ctx.SrvName, ctx.Request.URL.Path, ctx.ResponseWriter.StatusCode)
			if next != nil {
				next(ctx)
			}

			logx.With(ctx.Logger,
				"TimeStamp", time.Now().Format(time.RFC3339),
				"server", ctx.SrvName,
			).Log(logx.LevelInfo,
				"clientIp", ctx.Request.RemoteAddr,
				"path", ctx.Request.URL.Path,
				"statuscode", ctx.ResponseWriter.StatusCode,
				"latency", time.Since(begintime).Seconds(),
			)

		}
	}
}

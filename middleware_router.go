package gap

func RouterProcessor() Middleware {
	return func(next MiddlewareHandler) MiddlewareHandler {
		return func(ctx *Context) {

			if h, ok := ctx.GetRouteHandler(ctx.Request.URL.Path); ok {
				h(ctx)
			}

			if next != nil {
				next(ctx)
			}
		}
	}
}

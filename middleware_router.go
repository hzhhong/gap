package gap

func Router() Middleware {
	return func(next MiddlewareHandler) MiddlewareHandler {
		return func(ctx *Context) {

			if h, ok := ctx.Router[ctx.Request.URL.Path]; ok {
				h(ctx)
			}

			if next != nil {
				next(ctx)
			}
		}
	}
}
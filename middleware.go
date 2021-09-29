package gap

type MiddlewareHandler func(*Context)

type Middleware func(next MiddlewareHandler) MiddlewareHandler

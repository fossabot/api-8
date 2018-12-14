package api

type middleware func(Handler) Handler

func withMiddleware(origin Handler, middlewares ...middleware) Handler {
	h := origin
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

package middleware

import (
	"context"
	"errors"
	"github.com/hqdem/go-api-template/lib/xweb"
	"net/http"
)

func TimeoutMiddleware(timeoutHandlersConfig *xweb.HandlersConfig) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			handlerTimeout := timeoutHandlersConfig.GetHandlerTimeout(r.RequestURI, r.Method)
			ctx, cancel := context.WithTimeout(r.Context(), handlerTimeout)

			defer func() {
				cancel()
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()
			*r = *r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

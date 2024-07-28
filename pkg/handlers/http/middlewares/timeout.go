package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"
)

func TimeoutMiddleware(timeoutDuration time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeoutDuration)

			defer func() {
				cancel()
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					w.WriteHeader(http.StatusRequestTimeout)
				}
			}()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

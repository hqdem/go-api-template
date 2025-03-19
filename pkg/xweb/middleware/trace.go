package middleware

import (
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

func TracingMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, "", otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			method := r.Method
			url := r.URL.Path
			return fmt.Sprintf("%s %s", method, url)
		}))
	}
}

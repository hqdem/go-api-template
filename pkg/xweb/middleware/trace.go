package middleware

import (
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

func TracingMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			handler := otelhttp.NewHandler(next, "",
				otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
					method := r.Method
					url := r.URL.Path
					return fmt.Sprintf("%s %s", method, url)
				}),
				otelhttp.WithSpanOptions(
					trace.WithAttributes(
						attribute.String("http.request_id", getRequestID(r.Context())),
					),
				),
			)

			handler.ServeHTTP(w, r)

		}
		return http.HandlerFunc(fn)
	}
}

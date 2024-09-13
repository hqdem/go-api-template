package middleware

import (
	"github.com/hqdem/go-api-template/pkg/xlog"
	"github.com/hqdem/go-api-template/pkg/xweb/middleware"
	"go.uber.org/zap"
	"net/http"
)

const requestIDLogField = "request_id"

func LogRequestIDMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestID := w.Header().Get(middleware.RequestIDHeader)
			ctx := xlog.WithFields(r.Context(), zap.String(requestIDLogField, requestID))

			*r = *r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

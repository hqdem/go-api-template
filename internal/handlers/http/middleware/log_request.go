package middleware

import (
	"bytes"
	"fmt"
	"github.com/hqdem/go-api-template/pkg/xlog"
	"github.com/hqdem/go-api-template/pkg/xweb/middleware"
	"io"
	"net/http"
)

func LogRequestMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			xlog.Info(ctx, fmt.Sprintf("request uri: %s", r.RequestURI))

			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				xlog.Error(ctx, fmt.Sprintf("can't read request body: %v", err))
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			xlog.Info(ctx, fmt.Sprintf("request body: %s", requestBody))
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

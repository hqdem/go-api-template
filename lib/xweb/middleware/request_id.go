package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"strings"
)

const (
	RequestIDHeader = "X-Request-ID"
	RequestIDSize   = 32
)

type requestIDCtxKey struct{}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"0123456789"

func generateRandString(n int) string {
	randString := strings.Builder{}
	randString.Grow(n)

	for i := 0; i < n; i++ {
		_ = randString.WriteByte(charset[rand.Intn(len(charset))])
	}
	return randString.String()
}

func RequestIDMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = generateRandString(RequestIDSize)
			}

			ctx := context.WithValue(r.Context(), requestIDCtxKey{}, requestID)
			r = r.WithContext(ctx)

			w.Header().Set(RequestIDHeader, requestID)
			next.ServeHTTP(w, r)

		}

		return http.HandlerFunc(fn)
	}
}

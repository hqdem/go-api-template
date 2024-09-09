package middleware

import (
	"fmt"
	"github.com/hqdem/go-api-template/lib/xlog"
	"github.com/hqdem/go-api-template/lib/xweb/middleware"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LogRequestIDMiddlewareTestSuite struct {
	suite.Suite
}

func (s *LogRequestIDMiddlewareTestSuite) TestMiddleware() {
	testCases := []struct {
		Name           string
		HandlerFn      http.HandlerFunc
		RequestHeaders map[string]string
	}{
		{
			Name: "success_no_header",
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			},
			RequestHeaders: map[string]string{},
		},
		{
			Name: "success_with_header",
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			},
			RequestHeaders: map[string]string{
				middleware.RequestIDHeader: "123",
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			handler := LogRequestIDMiddleware()(testCase.HandlerFn)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "localhost:8080/operation", nil)

			for header, value := range testCase.RequestHeaders {
				w.Header().Set(header, value)
			}

			handler.ServeHTTP(w, r)

			ctx := r.Context()
			fields := xlog.GetContextFields(ctx)
			for _, f := range fields {
				if f.Key == requestIDLogField {
					s.Require().Equal(testCase.RequestHeaders[middleware.RequestIDHeader], f.String)
					return
				}
			}
			s.Fail(fmt.Sprintf("no %s was found in ctx fields", requestIDLogField))
		})
	}
}

func TestLogRequestIDMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(LogRequestIDMiddlewareTestSuite))
}

package middleware

import (
	"context"
	"encoding/json"
	"github.com/hqdem/go-api-template/pkg/xotel"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TracingMiddlewareTestSuite struct {
	suite.Suite
	shutdownFn func(context.Context) error
}

func (s *TracingMiddlewareTestSuite) SetupSuite() {
	shutdownOtel, err := xotel.SetupOTelSDK(context.Background())
	s.Require().NoError(err)
	s.shutdownFn = shutdownOtel
}

func (s *TracingMiddlewareTestSuite) TearDownSuite() {
	err := s.shutdownFn(context.Background())
	s.Require().NoError(err)
}

func (s *TracingMiddlewareTestSuite) TestTracingMiddleware() {
	testCases := []struct {
		Name      string
		HandlerFn http.HandlerFunc // returns json span ctx
	}{
		{
			Name: "success_trace_request",
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				spanCtx := trace.SpanFromContext(ctx).SpanContext()

				w.WriteHeader(http.StatusOK)
				isValid, _ := json.Marshal(spanCtx.IsValid())
				_, _ = w.Write(isValid)

			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			handler := TracingMiddleware()(testCase.HandlerFn)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "http://localhost:8080/test_span_operation", nil)

			handler.ServeHTTP(w, r)

			s.Require().Equal(http.StatusOK, w.Code)

			body, _ := io.ReadAll(w.Body)
			var isValidSpan bool
			_ = json.Unmarshal(body, &isValidSpan)

			s.Require().True(isValidSpan)

		})
	}
}

func TestTracingMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(TracingMiddlewareTestSuite))
}

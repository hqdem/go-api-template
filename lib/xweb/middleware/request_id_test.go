package middleware

import (
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RequestIDMiddlewareTestSuite struct {
	suite.Suite
}

func (s *RequestIDMiddlewareTestSuite) TestMiddleware() {
	testCases := []struct {
		Name           string
		HandlerFn      http.HandlerFunc
		RequestHeaders map[string]string
	}{
		{
			Name: "success_no_initial_header",
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			},
			RequestHeaders: nil,
		},
		{
			Name: "success_with_initial_header",
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			},
			RequestHeaders: map[string]string{
				RequestIDHeader: "123",
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			handler := RequestIDMiddleware()(testCase.HandlerFn)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "localhost:8080/operation", nil)

			for header, value := range testCase.RequestHeaders {
				r.Header.Set(header, value)
			}

			handler.ServeHTTP(w, r)

			s.Require().NotEmpty(w.Header().Get(RequestIDHeader))

			requestHeaderRequestID := testCase.RequestHeaders[RequestIDHeader]
			if requestHeaderRequestID != "" {
				s.Require().Equal(requestHeaderRequestID, w.Header().Get(RequestIDHeader))
			}

		})
	}
}

func TestRequestIDMiddlewareTest(t *testing.T) {
	suite.Run(t, new(RequestIDMiddlewareTestSuite))
}

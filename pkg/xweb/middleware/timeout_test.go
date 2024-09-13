package middleware

import (
	"context"
	"github.com/hqdem/go-api-template/pkg/xweb"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TimeoutMiddlewareTestSuite struct {
	suite.Suite
}

func (s *TimeoutMiddlewareTestSuite) TestMiddleware() {
	testCases := []struct {
		Name               string
		TimeoutConfig      *xweb.HandlersConfig
		HandlerFn          http.HandlerFunc
		ExpectedStatusCode int
		ExpectedError      error
	}{
		{
			Name: "success_no_timeout",
			TimeoutConfig: &xweb.HandlersConfig{
				DefaultTimeoutSecs: 5,
				HandlersTimeouts:   nil,
			},
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedError:      context.Canceled,
		},
		{
			Name: "success_handler_default_timeout",
			TimeoutConfig: &xweb.HandlersConfig{
				DefaultTimeoutSecs: 1,
				HandlersTimeouts:   nil,
			},
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
			},
			ExpectedStatusCode: http.StatusGatewayTimeout,
			ExpectedError:      context.DeadlineExceeded,
		},
		{
			Name: "success_handler_timeout_matched_config",
			TimeoutConfig: &xweb.HandlersConfig{
				DefaultTimeoutSecs: 1000,
				HandlersTimeouts: []xweb.HandlerTimeoutConfig{
					{
						Method:      "GET",
						Endpoint:    "/path",
						TimeoutSecs: 1,
					},
				},
			},
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
			},
			ExpectedStatusCode: http.StatusGatewayTimeout,
			ExpectedError:      context.DeadlineExceeded,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			handler := TimeoutMiddleware(testCase.TimeoutConfig)(testCase.HandlerFn)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/path", nil)
			handler.ServeHTTP(w, r)

			ctx := r.Context()
			err := ctx.Err()
			s.Require().Error(err)
			s.Require().ErrorIs(testCase.ExpectedError, err)
			s.Require().Equal(testCase.ExpectedStatusCode, w.Code)
		})
	}
}

func TestTimeoutMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(TimeoutMiddlewareTestSuite))
}

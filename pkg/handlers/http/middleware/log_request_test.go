package middleware

import (
	"bytes"
	"errors"
	"github.com/hqdem/go-api-template/pkg/tests"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LogRequestMiddlewareTestSuite struct {
	tests.BaseTestSuite
}

func (s *LogRequestMiddlewareTestSuite) TestMiddleware() {
	testCases := []struct {
		Name               string
		HandlerFn          http.HandlerFunc
		Request            *http.Request
		ExpectedStatusCode int
	}{
		{
			Name: "success_log",
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				body, _ := io.ReadAll(r.Body)
				_, _ = w.Write(body)
			},
			Request:            httptest.NewRequest("GET", "localhost:8080/operation", bytes.NewBuffer([]byte("i'm the body"))),
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name: "read_body_error",
			HandlerFn: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				body, _ := io.ReadAll(r.Body)
				_, _ = w.Write(body)
			},
			Request:            createRequestWithBodyError(),
			ExpectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			handler := LogRequestMiddleware()(testCase.HandlerFn)
			w := httptest.NewRecorder()
			r := testCase.Request

			handler.ServeHTTP(w, r)
			s.Require().Equal(testCase.ExpectedStatusCode, w.Code)
		})
	}
}

func createRequestWithBodyError() *http.Request {
	r := httptest.NewRequest("GET", "localhost:8080/operation", nil)
	r.Body = &fakeReqBody{}
	return r
}

type fakeReqBody struct {
}

func (b *fakeReqBody) Read(p []byte) (n int, err error) {
	err = errors.New("unexpected error happened")
	return n, err
}

func (b *fakeReqBody) Close() error {
	return nil
}

func TestLogRequestMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(LogRequestMiddlewareTestSuite))
}

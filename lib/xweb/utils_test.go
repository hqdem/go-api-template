package xweb

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type WebRequestsUtilsTestSuite struct {
	suite.Suite
}

func (s *WebRequestsUtilsTestSuite) TestWriteAPIErrorResponse() {
	testCases := []struct {
		Name             string
		Error            error
		ExpectedHTTPCode int
		ExpectedRes      *apiErrorResponse
	}{
		{
			Name:             "success_write_coded_error",
			Error:            NewInternalError(errors.New("some error occurred")),
			ExpectedHTTPCode: http.StatusInternalServerError,
			ExpectedRes: &apiErrorResponse{
				Message: "some error occurred",
				Code:    InternalErrorCode,
				Details: map[string]any{},
			},
		},
		{
			Name:             "success_write_not_coded_error",
			Error:            errors.New("some error occurred"),
			ExpectedHTTPCode: http.StatusInternalServerError,
			ExpectedRes: &apiErrorResponse{
				Message: "some error occurred",
				Code:    InternalErrorCode,
				Details: map[string]any{},
			},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			w := httptest.NewRecorder()
			respWrapper := &ResponseHeaders{
				ResponseWriter: w,
			}
			err := writeAPIErrorResponse(respWrapper, testCase.Error)
			s.Require().NoError(err)

			s.Require().Equal(testCase.ExpectedHTTPCode, w.Code)
			actualContentType := w.Header().Get("Content-Type")
			s.Require().Equal("application/json", strings.ToLower(actualContentType))

			expectedBytes, err := json.Marshal(testCase.ExpectedRes)
			s.Require().NoError(err)

			responseBytes, err := io.ReadAll(w.Body)
			s.Require().NoError(err)
			s.Require().Equal(expectedBytes, responseBytes)
		})
	}
}

func (s *WebRequestsUtilsTestSuite) TestWriteAPIOKResponse() {
	testCases := []struct {
		Name             string
		ExpectedHTTPCode int
		Entity           any
	}{
		{
			Name:             "success_api_ok_write",
			ExpectedHTTPCode: http.StatusOK,
			Entity:           map[string]any{"field1": "test1"},
		},
		{
			Name:             "success_api_ok_write_not_default_status_code",
			ExpectedHTTPCode: http.StatusCreated,
			Entity:           map[string]any{"field2": "test2"},
		},
		{
			Name:             "success_api_ok_write_nil_entity",
			ExpectedHTTPCode: http.StatusOK,
			Entity:           nil,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			w := httptest.NewRecorder()
			respWrapper := &ResponseHeaders{
				ResponseWriter: w,
			}
			respWrapper.SetHTTPCode(testCase.ExpectedHTTPCode)

			err := writeAPIOKResponse(respWrapper, testCase.Entity)
			s.Require().NoError(err)
			s.Require().Equal(testCase.ExpectedHTTPCode, w.Code)
			actualContentType := w.Header().Get("Content-Type")
			s.Require().Equal("application/json", strings.ToLower(actualContentType))

			expectedResp := &apiOKResponse{Data: testCase.Entity}
			expectedBytes, err := json.Marshal(expectedResp)
			s.Require().NoError(err)

			responseBytes, err := io.ReadAll(w.Body)
			s.Require().NoError(err)
			s.Require().Equal(expectedBytes, responseBytes)
		})
	}
}

func TestWebRequestsUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(WebRequestsUtilsTestSuite))
}

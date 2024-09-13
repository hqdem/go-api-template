package xweb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func (s *WebRequestsUtilsTestSuite) TestFacadeHandlerAdapterFlow() {
	type nopFacadeType struct{}
	type nopResType struct{}
	defaultCtxTimeout := time.Millisecond * 500

	testCases := []struct {
		Name               string
		HandlerFn          func(ctx context.Context, w *ResponseHeaders, r *http.Request, facade *nopFacadeType) (*nopResType, error)
		OnCtxDoneHook      OnCtxDoneHookT
		OnPanicHook        OnPanicFnHookT
		OnHandlerDoneHook  OnHandlerDoneHookT
		ExpectedStatusCode int
		IsPanic            bool
		IsCtxDone          bool
		IsHandlerDone      bool
	}{
		{
			Name: "success_handle_without_error",
			HandlerFn: func(ctx context.Context, w *ResponseHeaders, r *http.Request, facade *nopFacadeType) (*nopResType, error) {
				return &nopResType{}, nil
			},
			OnCtxDoneHook:      nopOnCtxDoneHook,
			OnPanicHook:        nopOnPanicHook,
			OnHandlerDoneHook:  nopOnHandlerDoneHook,
			ExpectedStatusCode: http.StatusOK,
			IsPanic:            false,
			IsCtxDone:          false,
			IsHandlerDone:      true,
		},
		{
			Name: "success_handle_with_error",
			HandlerFn: func(ctx context.Context, w *ResponseHeaders, r *http.Request, facade *nopFacadeType) (*nopResType, error) {
				return nil, errors.New("some unexpected error")
			},
			OnCtxDoneHook:      nopOnCtxDoneHook,
			OnPanicHook:        nopOnPanicHook,
			OnHandlerDoneHook:  nopOnHandlerDoneHook,
			ExpectedStatusCode: http.StatusInternalServerError,
			IsPanic:            false,
			IsCtxDone:          false,
			IsHandlerDone:      true,
		},
		{
			Name: "panic_while_handle",
			HandlerFn: func(ctx context.Context, w *ResponseHeaders, r *http.Request, facade *nopFacadeType) (*nopResType, error) {
				panic("123")
			},
			OnCtxDoneHook:      nopOnCtxDoneHook,
			OnPanicHook:        nopOnPanicHook,
			OnHandlerDoneHook:  nopOnHandlerDoneHook,
			ExpectedStatusCode: http.StatusInternalServerError,
			IsPanic:            true,
			IsCtxDone:          false,
			IsHandlerDone:      false,
		},
		{
			Name: "ctx_timeout_while_handle",
			HandlerFn: func(ctx context.Context, w *ResponseHeaders, r *http.Request, facade *nopFacadeType) (*nopResType, error) {
				time.Sleep(defaultCtxTimeout * 2)
				return nil, nil
			},
			OnCtxDoneHook:      nopOnCtxDoneHook,
			OnPanicHook:        nopOnPanicHook,
			OnHandlerDoneHook:  nopOnHandlerDoneHook,
			ExpectedStatusCode: http.StatusOK, // We are not setting any status code when ctx deadline is exceeded. It's logic for timeout middleware
			IsPanic:            false,
			IsCtxDone:          true,
			IsHandlerDone:      false,
		},
	}

	for _, testCase := range testCases {
		isOnPanicHookCalled := false
		isOnCtxDoneHookCalled := false
		isOnHandlerDoneHookCalled := false

		s.Run(testCase.Name, func() {
			SetPanicFnHook(func(ctx context.Context, panicErr error, panicStack []byte) {
				isOnPanicHookCalled = true
				testCase.OnPanicHook(ctx, panicErr, panicStack)
			})
			SetCtxDoneHook(func(ctx context.Context) {
				isOnCtxDoneHookCalled = true
				testCase.OnCtxDoneHook(ctx)
			})
			SetHandlerDoneHook(func(ctx context.Context, res any, err error) {
				isOnHandlerDoneHookCalled = true
				testCase.OnHandlerDoneHook(ctx, res, err)
			})

			handler := FacadeHandlerAdapter(&nopFacadeType{}, testCase.HandlerFn)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "localhost:8080/operation", nil)
			ctx, cancel := context.WithTimeout(r.Context(), defaultCtxTimeout)
			defer cancel()
			r = r.WithContext(ctx)

			handler.ServeHTTP(w, r)

			s.Require().Equal(testCase.ExpectedStatusCode, w.Code)
			s.Require().Equal(testCase.IsPanic, isOnPanicHookCalled)
			s.Require().Equal(testCase.IsCtxDone, isOnCtxDoneHookCalled)
			s.Require().Equal(testCase.IsHandlerDone, isOnHandlerDoneHookCalled)
		})
	}
}

func TestWebRequestsUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(WebRequestsUtilsTestSuite))
}

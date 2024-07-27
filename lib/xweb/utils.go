package xweb

import (
	"context"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
)

func FacadeHandlerAdapter[FacadeT any, RespT any](
	facade FacadeT,
	f func(ctx context.Context, w *ResponseHeaders, facade FacadeT) (RespT, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		//	TODO: tracing here
		ctx := request.Context()
		responseWrapper := &ResponseHeaders{
			ResponseWriter: w,
		}

		res, err := f(ctx, responseWrapper, facade)
		if err != nil {
			_ = writeAPIErrorResponse(responseWrapper, err)
			return
		}
		_ = writeAPIOKResponse(responseWrapper, res)
	}
}

type ResponseHeaders struct {
	http.ResponseWriter
	httpCode int
}

func (w *ResponseHeaders) SetHTTPCode(httpCode int) {
	w.httpCode = httpCode
}

type apiErrorResponse struct {
	Message string         `json:"message"`
	Code    string         `json:"code"`
	Details map[string]any `json:"details"`
}

func writeAPIErrorResponse(w *ResponseHeaders, err error) error {
	var codedErr CodedError
	if !errors.As(err, &codedErr) {
		codedErr = NewInternalError(err)
	}

	resp := &apiErrorResponse{
		Message: codedErr.Error(),
		Code:    codedErr.CharCode(),
		Details: codedErr.Details(),
	}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error("can not marshall coded error", zap.Error(err))
		return err
	}

	w.WriteHeader(codedErr.HTTPCode())
	_, writingErr := w.Write(jsonBytes)
	if writingErr != nil {
		zap.L().Error("error writing error response", zap.Error(writingErr))
		return writingErr
	}
	return nil
}

type apiOKResponse struct {
	Data any `json:"data"`
}

func writeAPIOKResponse(w *ResponseHeaders, entity any) error {
	resp := &apiOKResponse{
		Data: entity,
	}
	content, err := json.Marshal(resp)
	if err != nil {
		zap.L().Error("can not marshall response entity", zap.Error(err))
		return err
	}
	if w.httpCode == 0 {
		w.httpCode = http.StatusOK
	}
	w.WriteHeader(w.httpCode)
	_, err = w.Write(content)
	if err != nil {
		zap.L().Error("error writing ok api response", zap.Error(err))
		return err
	}
	return nil
}

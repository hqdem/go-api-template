package xweb

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/hqdem/go-api-template/lib/xlog"
	"go.uber.org/zap"
	"net/http"
)

func FacadeHandlerAdapter[FacadeT any, RespT any](
	facade FacadeT,
	f func(ctx context.Context, facade FacadeT) (RespT, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		//	TODO: tracing here
		ctx := request.Context()

		res, err := f(ctx, facade)
		if err != nil {
			var codedErr CodedError
			if !errors.As(err, &codedErr) {
				codedErr = NewInternalError(err)
			}
			_ = writeAPIErrorResponse(w, codedErr)
			return
		}
		_ = writeAPIOKResponse(w, res)
	}
}

func writeAPIErrorResponse(w http.ResponseWriter, codedErr CodedError) error {
	type apiErrorResponse struct {
		Message string         `json:"message"`
		Code    string         `json:"code"`
		Details map[string]any `json:"details"`
	}
	resp := apiErrorResponse{
		Message: codedErr.Error(),
		Code:    codedErr.CharCode(),
		Details: codedErr.Details(),
	}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		xlog.Error("can not marshall coded error", zap.Error(err))
		return err
	}

	w.WriteHeader(codedErr.HTTPCode())
	_, writingErr := w.Write(jsonBytes)
	if writingErr != nil {
		xlog.Error("error writing error response", zap.Error(writingErr))
		return writingErr
	}
	return nil
}

func writeAPIOKResponse(w http.ResponseWriter, entity any) error {
	content, err := json.Marshal(entity)
	if err != nil {
		xlog.Error("can not marshall response entity", zap.Error(err))
		return err
	}
	// TODO: custom http code
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(content)
	if err != nil {
		xlog.Error("error writing ok api response", zap.Error(err))
		return err
	}
	return nil
}

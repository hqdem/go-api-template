package xweb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

type OnPanicFnHookT func(ctx context.Context, panicErr error, panicStack []byte)
type OnCtxDoneHookT func(ctx context.Context)

var (
	nopOnPanicHook   OnPanicFnHookT = func(ctx context.Context, panicErr error, panicStack []byte) {}
	nopOnCtxDoneHook OnCtxDoneHookT = func(ctx context.Context) {}
	onPanicHook                     = nopOnPanicHook
	onCtxDoneHook                   = nopOnCtxDoneHook
)

// SetPanicFnHook is not thread safe and should be only called on application start
func SetPanicFnHook(fn OnPanicFnHookT) {
	if fn != nil {
		onPanicHook = fn
	}
}

// SetCtxDoneHook is not thread safe and should be only called on application start
func SetCtxDoneHook(fn OnCtxDoneHookT) {
	if fn != nil {
		onCtxDoneHook = fn
	}
}

type panicMessage struct {
	panicErr   error
	panicStack []byte
}

func FacadeHandlerAdapter[FacadeT any, RespT any](
	facade FacadeT,
	f func(ctx context.Context, w *ResponseHeaders, facade FacadeT) (RespT, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		// TODO: tracing here
		ctx := request.Context()

		var (
			res RespT
			err error
		)
		responseWrapper := &ResponseHeaders{
			ResponseWriter: w,
		}
		doneFnCh := make(chan struct{})
		panicCh := make(chan panicMessage)

		go func() {
			defer func() {
				v := recover()
				if v != nil {
					panicErr := fmt.Errorf("unexpected panic happened: %v", v)
					panicCh <- panicMessage{
						panicErr:   panicErr,
						panicStack: debug.Stack(),
					}
				}
			}()
			res, err = f(ctx, responseWrapper, facade)
			doneFnCh <- struct{}{}
		}()

		for {
			select {
			case <-ctx.Done():
				onCtxDoneHook(ctx)
				return
			case <-doneFnCh:
				handleWriteResponse(responseWrapper, res, err)
				return
			case panicMsg := <-panicCh:
				onPanicHook(ctx, panicMsg.panicErr, panicMsg.panicStack)
				handleWriteResponse(responseWrapper, res, NewInternalError(panicMsg.panicErr))
				return
			}

		}
	}
}

func handleWriteResponse[RespT any](responseWrapper *ResponseHeaders, res RespT, err error) {
	if err != nil {
		writingErr := writeAPIErrorResponse(responseWrapper, err)
		if writingErr != nil {
			panic(writingErr)
		}
		return
	}
	writingErr := writeAPIOKResponse(responseWrapper, res)
	if writingErr != nil {
		panic(writingErr)
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
		return fmt.Errorf("can not marshall coded error: %w", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(codedErr.HTTPCode())
	_, writingErr := w.Write(jsonBytes)
	if writingErr != nil {
		return fmt.Errorf("error writing error response: %w", writingErr)
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
		return fmt.Errorf("can not marshall response entity: %w", err)
	}
	if w.httpCode == 0 {
		w.httpCode = http.StatusOK
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(w.httpCode)
	_, err = w.Write(content)
	if err != nil {
		return fmt.Errorf("error writing ok api response: %w", err)
	}
	return nil
}

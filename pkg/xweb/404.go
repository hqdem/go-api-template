package xweb

import (
	"encoding/json"
	"net/http"
)

var (
	notFoundErr = apiErrorResponse{
		Message: "requested source was not found",
		Code:    NotFoundErrorCode,
		Details: nil,
	}
)

func NotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	jsonBytes, _ := json.Marshal(notFoundErr)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write(jsonBytes)
}

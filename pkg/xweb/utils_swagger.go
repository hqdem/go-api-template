package xweb

// APIErrorResponse is alias for internal apiErrorResponse
// it should be in sync with internal type for correct swagger generation
type APIErrorResponse = apiErrorResponse

// ApiOKResponse is alias for internal apiOkResponse
// it should be in sync with internal type for correct swagger generation
type ApiOKResponse[T any] struct {
	Data T `json:"data"`
}

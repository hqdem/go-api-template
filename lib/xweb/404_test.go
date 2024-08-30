package xweb

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type NotFoundHandlerTestSuite struct {
	suite.Suite
}

func (s *NotFoundHandlerTestSuite) TestHandler() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "localhost:8080/some/unexisting/path", nil)
	expectedRes, err := json.Marshal(notFoundErr)
	s.Require().NoError(err)

	NotFoundHandler(w, r)
	s.Require().Equal(w.Code, http.StatusNotFound)
	s.Require().Equal(expectedRes, w.Body.Bytes())
}

func TestNotFoundHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(NotFoundHandlerTestSuite))
}

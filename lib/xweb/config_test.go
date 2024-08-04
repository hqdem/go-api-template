package xweb

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type GetHandlersTimeoutTestSuite struct {
	suite.Suite
}

func (s *GetHandlersTimeoutTestSuite) TestGetHandlerTimeout() {
	defaultHandlersConfig := &HandlersConfig{
		DefaultTimeoutSecs: 60,
		HandlersTimeouts: []HandlerTimeoutConfig{
			{
				Endpoint:    "/ping",
				Method:      "GET",
				TimeoutSecs: 1,
			},
		},
	}

	testCases := []struct {
		Name            string
		HandlersConfig  *HandlersConfig
		RequestURI      string
		Method          string
		ExpectedTimeout time.Duration
	}{
		{
			Name:            "default_matching_handler",
			HandlersConfig:  defaultHandlersConfig,
			RequestURI:      "/ping",
			Method:          "GET",
			ExpectedTimeout: time.Second * 1,
		},
		{
			Name:            "matching_with_query_params",
			HandlersConfig:  defaultHandlersConfig,
			RequestURI:      "/ping?x=6",
			Method:          "GET",
			ExpectedTimeout: time.Second * 1,
		},
		{
			Name:            "matching_method_case_insensitivity",
			HandlersConfig:  defaultHandlersConfig,
			RequestURI:      "/ping",
			Method:          "get",
			ExpectedTimeout: time.Second * 1,
		},
		{
			Name:            "not_matched_request_uri",
			HandlersConfig:  defaultHandlersConfig,
			RequestURI:      "/pingooooooo",
			Method:          "GET",
			ExpectedTimeout: time.Second * 60,
		},
		{
			Name:            "not_matched_request_method",
			HandlersConfig:  defaultHandlersConfig,
			RequestURI:      "/ping",
			Method:          "POST",
			ExpectedTimeout: time.Second * 60,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			timeout := testCase.HandlersConfig.GetHandlerTimeout(testCase.RequestURI, testCase.Method)
			s.Require().Equal(testCase.ExpectedTimeout, timeout)
		})
	}
}

func TestGetHandlersTimeoutTestSuite(t *testing.T) {
	suite.Run(t, new(GetHandlersTimeoutTestSuite))
}

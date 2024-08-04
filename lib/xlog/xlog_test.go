package xlog

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zapcore"
	"testing"
)

type BaseLoggerTestSuite struct {
	suite.Suite
}

func (s *BaseLoggerTestSuite) TestDispatch() {
	testCases := []struct {
		Name             string
		Level            string
		ExpectedZapLevel zapcore.Level
		Panics           bool
	}{
		{
			Name:             "success_dispatch",
			Level:            "FATAL",
			ExpectedZapLevel: zapcore.FatalLevel,
			Panics:           false,
		},
		{
			Name:             "success_dispatch_case_sensitivity",
			Level:            "FaTaL",
			ExpectedZapLevel: zapcore.FatalLevel,
			Panics:           false,
		},
		{
			Name:   "unknown_level",
			Level:  "crazy",
			Panics: true,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			if testCase.Panics {
				s.Panics(func() {
					getZapLogLevel(testCase.Level)
				})
			} else {
				level := getZapLogLevel(testCase.Level)
				s.Require().Equal(testCase.ExpectedZapLevel, level)
			}
		})
	}
}

func TestBaseLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(BaseLoggerTestSuite))
}

package xlog

import (
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type BaseLoggerTestSuite struct {
	suite.Suite
}

func (s *BaseLoggerTestSuite) TestSetGetLogger() {
	logger, err := zap.NewDevelopment()
	s.Require().NoError(err)
	SetGlobalLogger(logger)
	s.Equal(logger, GetGlobalLogger())
}

func TestBaseLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(BaseLoggerTestSuite))
}

package tests

import (
	"github.com/hqdem/go-api-template/lib/xlog"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
}

func (s *BaseTestSuite) SetupSuite() {
	err := xlog.SetDefaultLogger(xlog.DEBUG, true)
	if err != nil {
		panic(err)
	}
}

package xlog

import (
	"context"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type CtxLogTestSuite struct {
	suite.Suite
}

func (s *CtxLogTestSuite) TestMergeFields() {
	testCases := []struct {
		Name           string
		OldFields      []zap.Field
		NewFields      []zap.Field
		ExpectedFields []zap.Field
	}{
		{
			Name:           "success_new_and_old_are_not_empty",
			OldFields:      []zap.Field{zap.String("test1", "test1"), zap.String("test2", "test2")},
			NewFields:      []zap.Field{zap.String("test3", "test3"), zap.String("test4", "test4")},
			ExpectedFields: []zap.Field{zap.String("test1", "test1"), zap.String("test2", "test2"), zap.String("test3", "test3"), zap.String("test4", "test4")},
		},
		{
			Name:           "success_new_empty",
			OldFields:      []zap.Field{zap.String("test1", "test1"), zap.String("test2", "test2")},
			NewFields:      nil,
			ExpectedFields: []zap.Field{zap.String("test1", "test1"), zap.String("test2", "test2")},
		},
		{
			Name:           "success_old_empty",
			OldFields:      nil,
			NewFields:      []zap.Field{zap.String("test3", "test3"), zap.String("test4", "test4")},
			ExpectedFields: []zap.Field{zap.String("test3", "test3"), zap.String("test4", "test4")},
		},
		{
			Name:           "success_both_empty",
			OldFields:      nil,
			NewFields:      nil,
			ExpectedFields: nil,
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			actualFields := mergeFields(testCase.OldFields, testCase.NewFields)
			s.Require().Equal(testCase.ExpectedFields, actualFields)
		})
	}
}

func (s *CtxLogTestSuite) TestContextFields() {
	testCases := []struct {
		Name           string
		SetupContextFn func(ctx context.Context) context.Context
		NewFields      []zap.Field
		ExpectedFields []zap.Field
	}{
		{
			Name: "success_no_fields_in_base_ctx",
			SetupContextFn: func(ctx context.Context) context.Context {
				return ctx
			},
			NewFields:      []zap.Field{zap.Int("int", 1), zap.Float64("float64", 1.0)},
			ExpectedFields: []zap.Field{zap.Int("int", 1), zap.Float64("float64", 1.0)},
		},
		{
			Name: "success_some_fields_in_base_ctx",
			SetupContextFn: func(ctx context.Context) context.Context {
				ctx = WithFields(ctx, zap.Int("test", 1), zap.String("test", "test"))
				return ctx
			},
			NewFields:      []zap.Field{zap.Int("int", 1), zap.Float64("float64", 1.0)},
			ExpectedFields: []zap.Field{zap.Int("test", 1), zap.String("test", "test"), zap.Int("int", 1), zap.Float64("float64", 1.0)},
		},
		{
			Name: "success_with_zero_array_new_fields",
			SetupContextFn: func(ctx context.Context) context.Context {
				ctx = WithFields(ctx, zap.Int("test", 1), zap.String("test", "test"))
				return ctx
			},
			NewFields:      nil,
			ExpectedFields: []zap.Field{zap.Int("test", 1), zap.String("test", "test")},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			ctx := context.Background()
			ctx = testCase.SetupContextFn(ctx)
			ctx = WithFields(ctx, testCase.NewFields...)
			actualFields := GetContextFields(ctx)
			s.Require().Equal(testCase.ExpectedFields, actualFields)
		})
	}
}

func TestCtxLogTestSuite(t *testing.T) {
	suite.Run(t, new(CtxLogTestSuite))
}

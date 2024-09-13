package xweb

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type GenericCodedErrorTestSuite struct {
	suite.Suite
}

func (s *GenericCodedErrorTestSuite) TestGetProperties() {
	msg := "some error"
	err := errors.New("some error")
	charCode := "CHAR_CODE"
	httpCode := http.StatusInternalServerError
	details := map[string]any{
		"test": "test",
	}
	genericCodedErr := NewGenericCodedError(err, httpCode, charCode, details)

	s.Run("test_get_properties", func() {
		var codedErr CodedError
		s.True(errors.As(genericCodedErr, &codedErr))

		s.Equal(msg, codedErr.Error())
		s.Equal(charCode, codedErr.CharCode())
		s.Equal(httpCode, codedErr.HTTPCode())
		s.Equal(details, codedErr.Details())
	})
}

func (s *GenericCodedErrorTestSuite) TestSetDetail() {
	testCases := []struct {
		Name            string
		OriginalDetails map[string]any
		ChangeDetailsFn func(codedErr CodedError)
		ExpectedDetails map[string]any
	}{
		{
			Name: "set_new_detail",
			OriginalDetails: map[string]any{
				"test1": "test1",
			},
			ChangeDetailsFn: func(codedErr CodedError) {
				codedErr.SetDetail("test2", "test2")
			},
			ExpectedDetails: map[string]any{
				"test1": "test1",
				"test2": "test2",
			},
		},
		{
			Name: "reset_existing_detail",
			OriginalDetails: map[string]any{
				"test1": "test1",
			},
			ChangeDetailsFn: func(codedErr CodedError) {
				codedErr.SetDetail("test1", "test2")
			},
			ExpectedDetails: map[string]any{
				"test1": "test2",
			},
		},
		{
			Name:            "set_detail_with_nil_original_details",
			OriginalDetails: nil,
			ChangeDetailsFn: func(codedErr CodedError) {
				codedErr.SetDetail("test1", "test1")

			},
			ExpectedDetails: map[string]any{
				"test1": "test1",
			},
		},
		//{
		//	Name: "delete_existing_detail",
		//},
		//{
		//	Name: "delete_non_existing_detail",
		//},
		//{
		//	Name: "delete_nil_detail",
		//},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			err := errors.New("unexpected error")
			httpCode := http.StatusInternalServerError
			charCode := "CHAR_CODE"
			genericCodedError := NewGenericCodedError(err, httpCode, charCode, testCase.OriginalDetails)

			testCase.ChangeDetailsFn(genericCodedError)

			s.Equal(testCase.ExpectedDetails, genericCodedError.Details())
		})
	}
}

func (s *GenericCodedErrorTestSuite) TestRemoveDetail() {
	testCases := []struct {
		Name            string
		OriginalDetails map[string]any
		ChangeDetailsFn func(codedErr CodedError)
		ExpectedDetails map[string]any
	}{
		{
			Name: "remove_existing_detail",
			OriginalDetails: map[string]any{
				"test1": "test1",
				"test2": "test2",
			},
			ChangeDetailsFn: func(codedErr CodedError) {
				codedErr.RemoveDetail("test2")
			},
			ExpectedDetails: map[string]any{
				"test1": "test1",
			},
		},
		{
			Name: "remove_non_existing_detail",
			OriginalDetails: map[string]any{
				"test1": "test1",
			},
			ChangeDetailsFn: func(codedErr CodedError) {
				codedErr.RemoveDetail("test2")
			},
			ExpectedDetails: map[string]any{
				"test1": "test1",
			},
		},
		{
			Name:            "remove_nil_detail",
			OriginalDetails: nil,
			ChangeDetailsFn: func(codedErr CodedError) {
				codedErr.RemoveDetail("test2")
			},
			ExpectedDetails: map[string]any{},
		},
	}

	for _, testCase := range testCases {
		s.Run(testCase.Name, func() {
			err := errors.New("unexpected error")
			httpCode := http.StatusInternalServerError
			charCode := "CHAR_CODE"
			genericCodedError := NewGenericCodedError(err, httpCode, charCode, testCase.OriginalDetails)

			testCase.ChangeDetailsFn(genericCodedError)

			s.Equal(testCase.ExpectedDetails, genericCodedError.Details())
		})
	}
}

func TestGenericCodedErrorTestSuite(t *testing.T) {
	suite.Run(t, new(GenericCodedErrorTestSuite))
}

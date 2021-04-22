// Code generated by mockery v2.2.1. DO NOT EDIT.

package rates

import (
	mock "github.com/stretchr/testify/mock"
)

// RateAndMarginSource is an autogenerated mock type for the RateAndMarginSource type
type RateAndMarginSourceMock struct {
	mock.Mock
}

// FindRateAndMargin provides a mock function with given fields: base, reference
func (_m *RateAndMarginSourceMock) FindRateAndMargin(base string, reference string) (*RateAndMargin, error) {
	ret := _m.Called(base, reference)

	var r0 *RateAndMargin
	if rf, ok := ret.Get(0).(func(string, string) *RateAndMargin); ok {
		r0 = rf(base, reference)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RateAndMargin)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(base, reference)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

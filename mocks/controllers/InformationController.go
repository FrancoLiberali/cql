// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	httperrors "github.com/ditrit/badaas/httperrors"
	mock "github.com/stretchr/testify/mock"
)

// InformationController is an autogenerated mock type for the InformationController type
type InformationController struct {
	mock.Mock
}

// Info provides a mock function with given fields: response, r
func (_m *InformationController) Info(response http.ResponseWriter, r *http.Request) (interface{}, httperrors.HTTPError) {
	ret := _m.Called(response, r)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) interface{}); ok {
		r0 = rf(response, r)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 httperrors.HTTPError
	if rf, ok := ret.Get(1).(func(http.ResponseWriter, *http.Request) httperrors.HTTPError); ok {
		r1 = rf(response, r)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(httperrors.HTTPError)
		}
	}

	return r0, r1
}

type mockConstructorTestingTNewInformationController interface {
	mock.TestingT
	Cleanup(func())
}

// NewInformationController creates a new instance of InformationController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewInformationController(t mockConstructorTestingTNewInformationController) *InformationController {
	mock := &InformationController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

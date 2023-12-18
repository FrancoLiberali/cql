// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ModelID is an autogenerated mock type for the ModelID type
type ModelID struct {
	mock.Mock
}

// IsNil provides a mock function with given fields:
func (_m *ModelID) IsNil() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewModelID interface {
	mock.TestingT
	Cleanup(func())
}

// NewModelID creates a new instance of ModelID. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewModelID(t mockConstructorTestingTNewModelID) *ModelID {
	mock := &ModelID{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Model is an autogenerated mock type for the Model type
type Model struct {
	mock.Mock
}

// IsLoaded provides a mock function with given fields:
func (_m *Model) IsLoaded() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewModel interface {
	mock.TestingT
	Cleanup(func())
}

// NewModel creates a new instance of Model. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewModel(t mockConstructorTestingTNewModel) *Model {
	mock := &Model{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
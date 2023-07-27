// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	zap "go.uber.org/zap"
)

// PaginationConfiguration is an autogenerated mock type for the PaginationConfiguration type
type PaginationConfiguration struct {
	mock.Mock
}

// GetMaxElemPerPage provides a mock function with given fields:
func (_m *PaginationConfiguration) GetMaxElemPerPage() uint {
	ret := _m.Called()

	var r0 uint
	if rf, ok := ret.Get(0).(func() uint); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint)
	}

	return r0
}

// Log provides a mock function with given fields: logger
func (_m *PaginationConfiguration) Log(logger *zap.Logger) {
	_m.Called(logger)
}

// Reload provides a mock function with given fields:
func (_m *PaginationConfiguration) Reload() {
	_m.Called()
}

type mockConstructorTestingTNewPaginationConfiguration interface {
	mock.TestingT
	Cleanup(func())
}

// NewPaginationConfiguration creates a new instance of PaginationConfiguration. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPaginationConfiguration(t mockConstructorTestingTNewPaginationConfiguration) *PaginationConfiguration {
	mock := &PaginationConfiguration{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
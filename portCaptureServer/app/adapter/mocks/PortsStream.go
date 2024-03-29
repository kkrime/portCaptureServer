// Code generated by mockery v3.0.0-alpha.0. DO NOT EDIT.

package mocks

import (
	entity "portCaptureServer/app/entity"

	mock "github.com/stretchr/testify/mock"
)

// PortsStream is an autogenerated mock type for the PortsStream type
type PortsStream struct {
	mock.Mock
}

// Recv provides a mock function with given fields:
func (_m *PortsStream) Recv() (*entity.Port, error) {
	ret := _m.Called()

	var r0 *entity.Port
	if rf, ok := ret.Get(0).(func() *entity.Port); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Port)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPortsStream interface {
	mock.TestingT
	Cleanup(func())
}

// NewPortsStream creates a new instance of PortsStream. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPortsStream(t mockConstructorTestingTNewPortsStream) *PortsStream {
	mock := &PortsStream{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

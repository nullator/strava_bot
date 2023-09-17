// Code generated by mockery v2.33.2. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Base is an autogenerated mock type for the Base type
type Base struct {
	mock.Mock
}

// Get provides a mock function with given fields: key, bucket
func (_m *Base) Get(key string, bucket string) (string, error) {
	ret := _m.Called(key, bucket)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (string, error)); ok {
		return rf(key, bucket)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(key, bucket)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(key, bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: key, value, bucket
func (_m *Base) Save(key string, value string, bucket string) error {
	ret := _m.Called(key, value, bucket)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(key, value, bucket)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewBase creates a new instance of Base. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBase(t interface {
	mock.TestingT
	Cleanup(func())
}) *Base {
	mock := &Base{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

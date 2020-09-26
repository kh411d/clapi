// Code generated by mockery v2.2.2. DO NOT EDIT.

package mocks

import (
	context "context"

	db "github.com/kh411d/clapi/db"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// KV is an autogenerated mock type for the KV type
type KV struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0, _a1, _a2
func (_m *KV) Add(_a0 string, _a1 []byte, _a2 time.Duration) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, time.Duration) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: _a0
func (_m *KV) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: _a0
func (_m *KV) Get(_a0 string) ([]byte, error) {
	ret := _m.Called(_a0)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Incr provides a mock function with given fields: _a0
func (_m *KV) Incr(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IncrBy provides a mock function with given fields: _a0, _a1
func (_m *KV) IncrBy(_a0 string, _a1 int64) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int64) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Set provides a mock function with given fields: _a0, _a1, _a2
func (_m *KV) Set(_a0 string, _a1 []byte, _a2 time.Duration) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte, time.Duration) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithContext provides a mock function with given fields: _a0
func (_m *KV) WithContext(_a0 context.Context) db.KV {
	ret := _m.Called(_a0)

	var r0 db.KV
	if rf, ok := ret.Get(0).(func(context.Context) db.KV); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.KV)
		}
	}

	return r0
}

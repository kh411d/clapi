// Code generated by mockery v2.2.2. DO NOT EDIT.

package mocks

import (
	context "context"

	db "github.com/kh411d/clapi/db"
	mock "github.com/stretchr/testify/mock"
)

// Clapper is an autogenerated mock type for the Clapper type
type Clapper struct {
	mock.Mock
}

// AddClap provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *Clapper) AddClap(_a0 context.Context, _a1 db.KV, _a2 string, _a3 int64) {
	_m.Called(_a0, _a1, _a2, _a3)
}

// GetClap provides a mock function with given fields: _a0, _a1, _a2
func (_m *Clapper) GetClap(_a0 context.Context, _a1 db.KV, _a2 string) string {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, db.KV, string) string); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify
// Copyright (c) Abstract Machines

// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"context"

	"github.com/absmach/magistrala/alarms"
	mock "github.com/stretchr/testify/mock"
)

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// CreateAlarm provides a mock function for the type Repository
func (_mock *Repository) CreateAlarm(ctx context.Context, alarm alarms.Alarm) (alarms.Alarm, error) {
	ret := _mock.Called(ctx, alarm)

	if len(ret) == 0 {
		panic("no return value specified for CreateAlarm")
	}

	var r0 alarms.Alarm
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, alarms.Alarm) (alarms.Alarm, error)); ok {
		return returnFunc(ctx, alarm)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, alarms.Alarm) alarms.Alarm); ok {
		r0 = returnFunc(ctx, alarm)
	} else {
		r0 = ret.Get(0).(alarms.Alarm)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, alarms.Alarm) error); ok {
		r1 = returnFunc(ctx, alarm)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Repository_CreateAlarm_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAlarm'
type Repository_CreateAlarm_Call struct {
	*mock.Call
}

// CreateAlarm is a helper method to define mock.On call
//   - ctx
//   - alarm
func (_e *Repository_Expecter) CreateAlarm(ctx interface{}, alarm interface{}) *Repository_CreateAlarm_Call {
	return &Repository_CreateAlarm_Call{Call: _e.mock.On("CreateAlarm", ctx, alarm)}
}

func (_c *Repository_CreateAlarm_Call) Run(run func(ctx context.Context, alarm alarms.Alarm)) *Repository_CreateAlarm_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(alarms.Alarm))
	})
	return _c
}

func (_c *Repository_CreateAlarm_Call) Return(alarm1 alarms.Alarm, err error) *Repository_CreateAlarm_Call {
	_c.Call.Return(alarm1, err)
	return _c
}

func (_c *Repository_CreateAlarm_Call) RunAndReturn(run func(ctx context.Context, alarm alarms.Alarm) (alarms.Alarm, error)) *Repository_CreateAlarm_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteAlarm provides a mock function for the type Repository
func (_mock *Repository) DeleteAlarm(ctx context.Context, id string) error {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlarm")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = returnFunc(ctx, id)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// Repository_DeleteAlarm_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAlarm'
type Repository_DeleteAlarm_Call struct {
	*mock.Call
}

// DeleteAlarm is a helper method to define mock.On call
//   - ctx
//   - id
func (_e *Repository_Expecter) DeleteAlarm(ctx interface{}, id interface{}) *Repository_DeleteAlarm_Call {
	return &Repository_DeleteAlarm_Call{Call: _e.mock.On("DeleteAlarm", ctx, id)}
}

func (_c *Repository_DeleteAlarm_Call) Run(run func(ctx context.Context, id string)) *Repository_DeleteAlarm_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Repository_DeleteAlarm_Call) Return(err error) *Repository_DeleteAlarm_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *Repository_DeleteAlarm_Call) RunAndReturn(run func(ctx context.Context, id string) error) *Repository_DeleteAlarm_Call {
	_c.Call.Return(run)
	return _c
}

// ListAlarms provides a mock function for the type Repository
func (_mock *Repository) ListAlarms(ctx context.Context, pm alarms.PageMetadata) (alarms.AlarmsPage, error) {
	ret := _mock.Called(ctx, pm)

	if len(ret) == 0 {
		panic("no return value specified for ListAlarms")
	}

	var r0 alarms.AlarmsPage
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, alarms.PageMetadata) (alarms.AlarmsPage, error)); ok {
		return returnFunc(ctx, pm)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, alarms.PageMetadata) alarms.AlarmsPage); ok {
		r0 = returnFunc(ctx, pm)
	} else {
		r0 = ret.Get(0).(alarms.AlarmsPage)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, alarms.PageMetadata) error); ok {
		r1 = returnFunc(ctx, pm)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Repository_ListAlarms_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAlarms'
type Repository_ListAlarms_Call struct {
	*mock.Call
}

// ListAlarms is a helper method to define mock.On call
//   - ctx
//   - pm
func (_e *Repository_Expecter) ListAlarms(ctx interface{}, pm interface{}) *Repository_ListAlarms_Call {
	return &Repository_ListAlarms_Call{Call: _e.mock.On("ListAlarms", ctx, pm)}
}

func (_c *Repository_ListAlarms_Call) Run(run func(ctx context.Context, pm alarms.PageMetadata)) *Repository_ListAlarms_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(alarms.PageMetadata))
	})
	return _c
}

func (_c *Repository_ListAlarms_Call) Return(alarmsPage alarms.AlarmsPage, err error) *Repository_ListAlarms_Call {
	_c.Call.Return(alarmsPage, err)
	return _c
}

func (_c *Repository_ListAlarms_Call) RunAndReturn(run func(ctx context.Context, pm alarms.PageMetadata) (alarms.AlarmsPage, error)) *Repository_ListAlarms_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateAlarm provides a mock function for the type Repository
func (_mock *Repository) UpdateAlarm(ctx context.Context, alarm alarms.Alarm) (alarms.Alarm, error) {
	ret := _mock.Called(ctx, alarm)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAlarm")
	}

	var r0 alarms.Alarm
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, alarms.Alarm) (alarms.Alarm, error)); ok {
		return returnFunc(ctx, alarm)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, alarms.Alarm) alarms.Alarm); ok {
		r0 = returnFunc(ctx, alarm)
	} else {
		r0 = ret.Get(0).(alarms.Alarm)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, alarms.Alarm) error); ok {
		r1 = returnFunc(ctx, alarm)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Repository_UpdateAlarm_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateAlarm'
type Repository_UpdateAlarm_Call struct {
	*mock.Call
}

// UpdateAlarm is a helper method to define mock.On call
//   - ctx
//   - alarm
func (_e *Repository_Expecter) UpdateAlarm(ctx interface{}, alarm interface{}) *Repository_UpdateAlarm_Call {
	return &Repository_UpdateAlarm_Call{Call: _e.mock.On("UpdateAlarm", ctx, alarm)}
}

func (_c *Repository_UpdateAlarm_Call) Run(run func(ctx context.Context, alarm alarms.Alarm)) *Repository_UpdateAlarm_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(alarms.Alarm))
	})
	return _c
}

func (_c *Repository_UpdateAlarm_Call) Return(alarm1 alarms.Alarm, err error) *Repository_UpdateAlarm_Call {
	_c.Call.Return(alarm1, err)
	return _c
}

func (_c *Repository_UpdateAlarm_Call) RunAndReturn(run func(ctx context.Context, alarm alarms.Alarm) (alarms.Alarm, error)) *Repository_UpdateAlarm_Call {
	_c.Call.Return(run)
	return _c
}

// ViewAlarm provides a mock function for the type Repository
func (_mock *Repository) ViewAlarm(ctx context.Context, alarmID string, domainID string) (alarms.Alarm, error) {
	ret := _mock.Called(ctx, alarmID, domainID)

	if len(ret) == 0 {
		panic("no return value specified for ViewAlarm")
	}

	var r0 alarms.Alarm
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string) (alarms.Alarm, error)); ok {
		return returnFunc(ctx, alarmID, domainID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string) alarms.Alarm); ok {
		r0 = returnFunc(ctx, alarmID, domainID)
	} else {
		r0 = ret.Get(0).(alarms.Alarm)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = returnFunc(ctx, alarmID, domainID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Repository_ViewAlarm_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ViewAlarm'
type Repository_ViewAlarm_Call struct {
	*mock.Call
}

// ViewAlarm is a helper method to define mock.On call
//   - ctx
//   - alarmID
//   - domainID
func (_e *Repository_Expecter) ViewAlarm(ctx interface{}, alarmID interface{}, domainID interface{}) *Repository_ViewAlarm_Call {
	return &Repository_ViewAlarm_Call{Call: _e.mock.On("ViewAlarm", ctx, alarmID, domainID)}
}

func (_c *Repository_ViewAlarm_Call) Run(run func(ctx context.Context, alarmID string, domainID string)) *Repository_ViewAlarm_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *Repository_ViewAlarm_Call) Return(alarm alarms.Alarm, err error) *Repository_ViewAlarm_Call {
	_c.Call.Return(alarm, err)
	return _c
}

func (_c *Repository_ViewAlarm_Call) RunAndReturn(run func(ctx context.Context, alarmID string, domainID string) (alarms.Alarm, error)) *Repository_ViewAlarm_Call {
	_c.Call.Return(run)
	return _c
}

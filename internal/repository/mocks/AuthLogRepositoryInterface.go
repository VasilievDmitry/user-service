// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	pkg "github.com/lotproject/user-service/pkg"
	mock "github.com/stretchr/testify/mock"
)

// AuthLogRepositoryInterface is an autogenerated mock type for the AuthLogRepositoryInterface type
type AuthLogRepositoryInterface struct {
	mock.Mock
}

// GetByAccessToken provides a mock function with given fields: ctx, token
func (_m *AuthLogRepositoryInterface) GetByAccessToken(ctx context.Context, token string) (*pkg.AuthLog, error) {
	ret := _m.Called(ctx, token)

	var r0 *pkg.AuthLog
	if rf, ok := ret.Get(0).(func(context.Context, string) *pkg.AuthLog); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.AuthLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByRefreshToken provides a mock function with given fields: ctx, token
func (_m *AuthLogRepositoryInterface) GetByRefreshToken(ctx context.Context, token string) (*pkg.AuthLog, error) {
	ret := _m.Called(ctx, token)

	var r0 *pkg.AuthLog
	if rf, ok := ret.Get(0).(func(context.Context, string) *pkg.AuthLog); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.AuthLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: ctx, log
func (_m *AuthLogRepositoryInterface) Insert(ctx context.Context, log *pkg.AuthLog) error {
	ret := _m.Called(ctx, log)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.AuthLog) error); ok {
		r0 = rf(ctx, log)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, log
func (_m *AuthLogRepositoryInterface) Update(ctx context.Context, log *pkg.AuthLog) error {
	ret := _m.Called(ctx, log)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.AuthLog) error); ok {
		r0 = rf(ctx, log)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

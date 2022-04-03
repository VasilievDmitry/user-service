// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	client "github.com/micro/go-micro/client"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	mock "github.com/stretchr/testify/mock"

	pkg "github.com/lotproject/user-service/pkg"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// ConfirmLogin provides a mock function with given fields: ctx, in, opts
func (_m *UserService) ConfirmLogin(ctx context.Context, in *pkg.ConfirmLoginRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.ConfirmLoginRequest, ...client.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.ConfirmLoginRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateAuthToken provides a mock function with given fields: ctx, in, opts
func (_m *UserService) CreateAuthToken(ctx context.Context, in *pkg.CreateAuthTokenRequest, opts ...client.CallOption) (*pkg.ResponseWithAuthToken, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pkg.ResponseWithAuthToken
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.CreateAuthTokenRequest, ...client.CallOption) *pkg.ResponseWithAuthToken); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.ResponseWithAuthToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.CreateAuthTokenRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreatePasswordRecoveryCode provides a mock function with given fields: ctx, in, opts
func (_m *UserService) CreatePasswordRecoveryCode(ctx context.Context, in *pkg.CreatePasswordRecoveryCodeRequest, opts ...client.CallOption) (*pkg.CreatePasswordRecoveryCodeResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pkg.CreatePasswordRecoveryCodeResponse
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.CreatePasswordRecoveryCodeRequest, ...client.CallOption) *pkg.CreatePasswordRecoveryCodeResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.CreatePasswordRecoveryCodeResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.CreatePasswordRecoveryCodeRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUserByWallet provides a mock function with given fields: ctx, in, opts
func (_m *UserService) CreateUserByWallet(ctx context.Context, in *pkg.CreateUserByWalletRequest, opts ...client.CallOption) (*pkg.ResponseWithUserProfile, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pkg.ResponseWithUserProfile
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.CreateUserByWalletRequest, ...client.CallOption) *pkg.ResponseWithUserProfile); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.ResponseWithUserProfile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.CreateUserByWalletRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeactivateAuthToken provides a mock function with given fields: ctx, in, opts
func (_m *UserService) DeactivateAuthToken(ctx context.Context, in *pkg.DeactivateAuthTokenRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.DeactivateAuthTokenRequest, ...client.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.DeactivateAuthTokenRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByAccessToken provides a mock function with given fields: ctx, in, opts
func (_m *UserService) GetUserByAccessToken(ctx context.Context, in *pkg.GetUserByAccessTokenRequest, opts ...client.CallOption) (*pkg.ResponseWithUserProfile, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pkg.ResponseWithUserProfile
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.GetUserByAccessTokenRequest, ...client.CallOption) *pkg.ResponseWithUserProfile); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.ResponseWithUserProfile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.GetUserByAccessTokenRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserById provides a mock function with given fields: ctx, in, opts
func (_m *UserService) GetUserById(ctx context.Context, in *pkg.GetUserByIdRequest, opts ...client.CallOption) (*pkg.ResponseWithUserProfile, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pkg.ResponseWithUserProfile
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.GetUserByIdRequest, ...client.CallOption) *pkg.ResponseWithUserProfile); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.ResponseWithUserProfile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.GetUserByIdRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByLogin provides a mock function with given fields: ctx, in, opts
func (_m *UserService) GetUserByLogin(ctx context.Context, in *pkg.GetUserByLoginRequest, opts ...client.CallOption) (*pkg.ResponseWithUserProfile, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pkg.ResponseWithUserProfile
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.GetUserByLoginRequest, ...client.CallOption) *pkg.ResponseWithUserProfile); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.ResponseWithUserProfile)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.GetUserByLoginRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Ping provides a mock function with given fields: ctx, in, opts
func (_m *UserService) Ping(ctx context.Context, in *emptypb.Empty, opts ...client.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *emptypb.Empty, ...client.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *emptypb.Empty, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RefreshAccessToken provides a mock function with given fields: ctx, in, opts
func (_m *UserService) RefreshAccessToken(ctx context.Context, in *pkg.RefreshAccessTokenRequest, opts ...client.CallOption) (*pkg.ResponseWithAuthToken, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pkg.ResponseWithAuthToken
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.RefreshAccessTokenRequest, ...client.CallOption) *pkg.ResponseWithAuthToken); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.ResponseWithAuthToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.RefreshAccessTokenRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetLogin provides a mock function with given fields: ctx, in, opts
func (_m *UserService) SetLogin(ctx context.Context, in *pkg.SetLoginRequest, opts ...client.CallOption) (*pkg.SetLoginResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *pkg.SetLoginResponse
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.SetLoginRequest, ...client.CallOption) *pkg.SetLoginResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pkg.SetLoginResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.SetLoginRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetPassword provides a mock function with given fields: ctx, in, opts
func (_m *UserService) SetPassword(ctx context.Context, in *pkg.SetPasswordRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.SetPasswordRequest, ...client.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.SetPasswordRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetUsername provides a mock function with given fields: ctx, in, opts
func (_m *UserService) SetUsername(ctx context.Context, in *pkg.SetUsernameRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.SetUsernameRequest, ...client.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.SetUsernameRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UsePasswordRecoveryCode provides a mock function with given fields: ctx, in, opts
func (_m *UserService) UsePasswordRecoveryCode(ctx context.Context, in *pkg.UsePasswordRecoveryCodeRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.UsePasswordRecoveryCodeRequest, ...client.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.UsePasswordRecoveryCodeRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VerifyPassword provides a mock function with given fields: ctx, in, opts
func (_m *UserService) VerifyPassword(ctx context.Context, in *pkg.VerifyPasswordRequest, opts ...client.CallOption) (*emptypb.Empty, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *emptypb.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *pkg.VerifyPasswordRequest, ...client.CallOption) *emptypb.Empty); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*emptypb.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *pkg.VerifyPasswordRequest, ...client.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
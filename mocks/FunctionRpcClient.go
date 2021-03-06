// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import grpc "google.golang.org/grpc"
import mock "github.com/stretchr/testify/mock"
import rpc "github.com/graphql-editor/azure-functions-golang-worker/rpc"

// FunctionRpcClient is an autogenerated mock type for the FunctionRpcClient type
type FunctionRpcClient struct {
	mock.Mock
}

// EventStream provides a mock function with given fields: ctx, opts
func (_m *FunctionRpcClient) EventStream(ctx context.Context, opts ...grpc.CallOption) (rpc.FunctionRpc_EventStreamClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 rpc.FunctionRpc_EventStreamClient
	if rf, ok := ret.Get(0).(func(context.Context, ...grpc.CallOption) rpc.FunctionRpc_EventStreamClient); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rpc.FunctionRpc_EventStreamClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

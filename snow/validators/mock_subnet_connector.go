// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/luxdefi/node/snow/validators (interfaces: SubnetConnector)

// Package validators is a generated GoMock package.
package validators

import (
	context "context"
	reflect "reflect"

	ids "github.com/luxdefi/node/ids"
	gomock "go.uber.org/mock/gomock"
)

// MockSubnetConnector is a mock of SubnetConnector interface.
type MockSubnetConnector struct {
	ctrl     *gomock.Controller
	recorder *MockSubnetConnectorMockRecorder
}

// MockSubnetConnectorMockRecorder is the mock recorder for MockSubnetConnector.
type MockSubnetConnectorMockRecorder struct {
	mock *MockSubnetConnector
}

// NewMockSubnetConnector creates a new mock instance.
func NewMockSubnetConnector(ctrl *gomock.Controller) *MockSubnetConnector {
	mock := &MockSubnetConnector{ctrl: ctrl}
	mock.recorder = &MockSubnetConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubnetConnector) EXPECT() *MockSubnetConnectorMockRecorder {
	return m.recorder
}

// ConnectedSubnet mocks base method.
func (m *MockSubnetConnector) ConnectedSubnet(arg0 context.Context, arg1 ids.NodeID, arg2 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectedSubnet", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConnectedSubnet indicates an expected call of ConnectedSubnet.
func (mr *MockSubnetConnectorMockRecorder) ConnectedSubnet(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectedSubnet", reflect.TypeOf((*MockSubnetConnector)(nil).ConnectedSubnet), arg0, arg1, arg2)
}

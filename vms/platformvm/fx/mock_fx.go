// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/skychains/chain/vms/platformvm/fx (interfaces: Fx,Owner)
//
// Generated by this command:
//
//	mockgen -package=fx -destination=vms/platformvm/fx/mock_fx.go github.com/skychains/chain/vms/platformvm/fx Fx,Owner
//

// Package fx is a generated GoMock package.
package fx

import (
	reflect "reflect"

	snow "github.com/skychains/chain/snow"
	verify "github.com/skychains/chain/vms/components/verify"
	gomock "go.uber.org/mock/gomock"
)

// MockFx is a mock of Fx interface.
type MockFx struct {
	ctrl     *gomock.Controller
	recorder *MockFxMockRecorder
}

// MockFxMockRecorder is the mock recorder for MockFx.
type MockFxMockRecorder struct {
	mock *MockFx
}

// NewMockFx creates a new mock instance.
func NewMockFx(ctrl *gomock.Controller) *MockFx {
	mock := &MockFx{ctrl: ctrl}
	mock.recorder = &MockFxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFx) EXPECT() *MockFxMockRecorder {
	return m.recorder
}

// Bootstrapped mocks base method.
func (m *MockFx) Bootstrapped() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bootstrapped")
	ret0, _ := ret[0].(error)
	return ret0
}

// Bootstrapped indicates an expected call of Bootstrapped.
func (mr *MockFxMockRecorder) Bootstrapped() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bootstrapped", reflect.TypeOf((*MockFx)(nil).Bootstrapped))
}

// Bootstrapping mocks base method.
func (m *MockFx) Bootstrapping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bootstrapping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Bootstrapping indicates an expected call of Bootstrapping.
func (mr *MockFxMockRecorder) Bootstrapping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bootstrapping", reflect.TypeOf((*MockFx)(nil).Bootstrapping))
}

// CreateOutput mocks base method.
func (m *MockFx) CreateOutput(arg0 uint64, arg1 any) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOutput", arg0, arg1)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateOutput indicates an expected call of CreateOutput.
func (mr *MockFxMockRecorder) CreateOutput(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOutput", reflect.TypeOf((*MockFx)(nil).CreateOutput), arg0, arg1)
}

// Initialize mocks base method.
func (m *MockFx) Initialize(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockFxMockRecorder) Initialize(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockFx)(nil).Initialize), arg0)
}

// VerifyPermission mocks base method.
func (m *MockFx) VerifyPermission(arg0, arg1, arg2, arg3 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyPermission", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyPermission indicates an expected call of VerifyPermission.
func (mr *MockFxMockRecorder) VerifyPermission(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyPermission", reflect.TypeOf((*MockFx)(nil).VerifyPermission), arg0, arg1, arg2, arg3)
}

// VerifyTransfer mocks base method.
func (m *MockFx) VerifyTransfer(arg0, arg1, arg2, arg3 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyTransfer", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyTransfer indicates an expected call of VerifyTransfer.
func (mr *MockFxMockRecorder) VerifyTransfer(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyTransfer", reflect.TypeOf((*MockFx)(nil).VerifyTransfer), arg0, arg1, arg2, arg3)
}

// MockOwner is a mock of Owner interface.
type MockOwner struct {
	verify.IsNotState

	ctrl     *gomock.Controller
	recorder *MockOwnerMockRecorder
}

// MockOwnerMockRecorder is the mock recorder for MockOwner.
type MockOwnerMockRecorder struct {
	mock *MockOwner
}

// NewMockOwner creates a new mock instance.
func NewMockOwner(ctrl *gomock.Controller) *MockOwner {
	mock := &MockOwner{ctrl: ctrl}
	mock.recorder = &MockOwnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOwner) EXPECT() *MockOwnerMockRecorder {
	return m.recorder
}

// InitCtx mocks base method.
func (m *MockOwner) InitCtx(arg0 *snow.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InitCtx", arg0)
}

// InitCtx indicates an expected call of InitCtx.
func (mr *MockOwnerMockRecorder) InitCtx(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitCtx", reflect.TypeOf((*MockOwner)(nil).InitCtx), arg0)
}

// Verify mocks base method.
func (m *MockOwner) Verify() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify")
	ret0, _ := ret[0].(error)
	return ret0
}

// Verify indicates an expected call of Verify.
func (mr *MockOwnerMockRecorder) Verify() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockOwner)(nil).Verify))
}

// isState mocks base method.
func (m *MockOwner) isState() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "isState")
	ret0, _ := ret[0].(error)
	return ret0
}

// isState indicates an expected call of isState.
func (mr *MockOwnerMockRecorder) isState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isState", reflect.TypeOf((*MockOwner)(nil).isState))
}

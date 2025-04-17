// Code generated by MockGen. DO NOT EDIT.
// Source: vms/avm/txs/tx.go
//
// Generated by this command:
//
//	mockgen -source=vms/avm/txs/tx.go -destination=vms/avm/txs/mock_unsigned_tx.go -package=txs -exclude_interfaces=
//

// Package txs is a generated GoMock package.
package txs

import (
	reflect "reflect"

	ids "github.com/SkyChains/chain/ids"
	snow "github.com/SkyChains/chain/snow"
	set "github.com/SkyChains/chain/utils/set"
	lux "github.com/SkyChains/chain/vms/components/lux"
	gomock "go.uber.org/mock/gomock"
)

// MockUnsignedTx is a mock of UnsignedTx interface.
type MockUnsignedTx struct {
	ctrl     *gomock.Controller
	recorder *MockUnsignedTxMockRecorder
}

// MockUnsignedTxMockRecorder is the mock recorder for MockUnsignedTx.
type MockUnsignedTxMockRecorder struct {
	mock *MockUnsignedTx
}

// NewMockUnsignedTx creates a new mock instance.
func NewMockUnsignedTx(ctrl *gomock.Controller) *MockUnsignedTx {
	mock := &MockUnsignedTx{ctrl: ctrl}
	mock.recorder = &MockUnsignedTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsignedTx) EXPECT() *MockUnsignedTxMockRecorder {
	return m.recorder
}

// Bytes mocks base method.
func (m *MockUnsignedTx) Bytes() []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Bytes")
	ret0, _ := ret[0].([]byte)
	return ret0
}

// Bytes indicates an expected call of Bytes.
func (mr *MockUnsignedTxMockRecorder) Bytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Bytes", reflect.TypeOf((*MockUnsignedTx)(nil).Bytes))
}

// InitCtx mocks base method.
func (m *MockUnsignedTx) InitCtx(ctx *snow.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InitCtx", ctx)
}

// InitCtx indicates an expected call of InitCtx.
func (mr *MockUnsignedTxMockRecorder) InitCtx(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitCtx", reflect.TypeOf((*MockUnsignedTx)(nil).InitCtx), ctx)
}

// InputIDs mocks base method.
func (m *MockUnsignedTx) InputIDs() set.Set[ids.ID] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InputIDs")
	ret0, _ := ret[0].(set.Set[ids.ID])
	return ret0
}

// InputIDs indicates an expected call of InputIDs.
func (mr *MockUnsignedTxMockRecorder) InputIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InputIDs", reflect.TypeOf((*MockUnsignedTx)(nil).InputIDs))
}

// InputUTXOs mocks base method.
func (m *MockUnsignedTx) InputUTXOs() []*lux.UTXOID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InputUTXOs")
	ret0, _ := ret[0].([]*lux.UTXOID)
	return ret0
}

// InputUTXOs indicates an expected call of InputUTXOs.
func (mr *MockUnsignedTxMockRecorder) InputUTXOs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InputUTXOs", reflect.TypeOf((*MockUnsignedTx)(nil).InputUTXOs))
}

// NumCredentials mocks base method.
func (m *MockUnsignedTx) NumCredentials() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NumCredentials")
	ret0, _ := ret[0].(int)
	return ret0
}

// NumCredentials indicates an expected call of NumCredentials.
func (mr *MockUnsignedTxMockRecorder) NumCredentials() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NumCredentials", reflect.TypeOf((*MockUnsignedTx)(nil).NumCredentials))
}

// SetBytes mocks base method.
func (m *MockUnsignedTx) SetBytes(unsignedBytes []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBytes", unsignedBytes)
}

// SetBytes indicates an expected call of SetBytes.
func (mr *MockUnsignedTxMockRecorder) SetBytes(unsignedBytes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBytes", reflect.TypeOf((*MockUnsignedTx)(nil).SetBytes), unsignedBytes)
}

// Visit mocks base method.
func (m *MockUnsignedTx) Visit(visitor Visitor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Visit", visitor)
	ret0, _ := ret[0].(error)
	return ret0
}

// Visit indicates an expected call of Visit.
func (mr *MockUnsignedTxMockRecorder) Visit(visitor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Visit", reflect.TypeOf((*MockUnsignedTx)(nil).Visit), visitor)
}

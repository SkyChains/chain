// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/skychains/chain/vms/avm/txs/mempool (interfaces: Mempool)
//
// Generated by this command:
//
//	mockgen -package=mempool -destination=vms/avm/txs/mempool/mock_mempool.go github.com/skychains/chain/vms/avm/txs/mempool Mempool
//

// Package mempool is a generated GoMock package.
package mempool

import (
	reflect "reflect"

	ids "github.com/skychains/chain/ids"
	txs "github.com/skychains/chain/vms/avm/txs"
	gomock "go.uber.org/mock/gomock"
)

// MockMempool is a mock of Mempool interface.
type MockMempool struct {
	ctrl     *gomock.Controller
	recorder *MockMempoolMockRecorder
}

// MockMempoolMockRecorder is the mock recorder for MockMempool.
type MockMempoolMockRecorder struct {
	mock *MockMempool
}

// NewMockMempool creates a new mock instance.
func NewMockMempool(ctrl *gomock.Controller) *MockMempool {
	mock := &MockMempool{ctrl: ctrl}
	mock.recorder = &MockMempoolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMempool) EXPECT() *MockMempoolMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockMempool) Add(arg0 *txs.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockMempoolMockRecorder) Add(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockMempool)(nil).Add), arg0)
}

// Get mocks base method.
func (m *MockMempool) Get(arg0 ids.ID) (*txs.Tx, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*txs.Tx)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMempoolMockRecorder) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMempool)(nil).Get), arg0)
}

// GetDropReason mocks base method.
func (m *MockMempool) GetDropReason(arg0 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDropReason", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetDropReason indicates an expected call of GetDropReason.
func (mr *MockMempoolMockRecorder) GetDropReason(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDropReason", reflect.TypeOf((*MockMempool)(nil).GetDropReason), arg0)
}

// Iterate mocks base method.
func (m *MockMempool) Iterate(arg0 func(*txs.Tx) bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Iterate", arg0)
}

// Iterate indicates an expected call of Iterate.
func (mr *MockMempoolMockRecorder) Iterate(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Iterate", reflect.TypeOf((*MockMempool)(nil).Iterate), arg0)
}

// Len mocks base method.
func (m *MockMempool) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MockMempoolMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockMempool)(nil).Len))
}

// MarkDropped mocks base method.
func (m *MockMempool) MarkDropped(arg0 ids.ID, arg1 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MarkDropped", arg0, arg1)
}

// MarkDropped indicates an expected call of MarkDropped.
func (mr *MockMempoolMockRecorder) MarkDropped(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkDropped", reflect.TypeOf((*MockMempool)(nil).MarkDropped), arg0, arg1)
}

// Peek mocks base method.
func (m *MockMempool) Peek() (*txs.Tx, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peek")
	ret0, _ := ret[0].(*txs.Tx)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Peek indicates an expected call of Peek.
func (mr *MockMempoolMockRecorder) Peek() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peek", reflect.TypeOf((*MockMempool)(nil).Peek))
}

// Remove mocks base method.
func (m *MockMempool) Remove(arg0 ...*txs.Tx) {
	m.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Remove", varargs...)
}

// Remove indicates an expected call of Remove.
func (mr *MockMempoolMockRecorder) Remove(arg0 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockMempool)(nil).Remove), arg0...)
}

// RequestBuildBlock mocks base method.
func (m *MockMempool) RequestBuildBlock() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RequestBuildBlock")
}

// RequestBuildBlock indicates an expected call of RequestBuildBlock.
func (mr *MockMempoolMockRecorder) RequestBuildBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestBuildBlock", reflect.TypeOf((*MockMempool)(nil).RequestBuildBlock))
}

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SkyChains/chain/vms/platformvm/utxo (interfaces: Verifier)
//
// Generated by this command:
//
//	mockgen -package=utxo -destination=vms/platformvm/utxo/mock_verifier.go github.com/SkyChains/chain/vms/platformvm/utxo Verifier
//

// Package utxo is a generated GoMock package.
package utxo

import (
	reflect "reflect"

	ids "github.com/SkyChains/chain/ids"
	lux "github.com/SkyChains/chain/vms/components/lux"
	verify "github.com/SkyChains/chain/vms/components/verify"
	txs "github.com/SkyChains/chain/vms/platformvm/txs"
	gomock "go.uber.org/mock/gomock"
)

// MockVerifier is a mock of Verifier interface.
type MockVerifier struct {
	ctrl     *gomock.Controller
	recorder *MockVerifierMockRecorder
}

// MockVerifierMockRecorder is the mock recorder for MockVerifier.
type MockVerifierMockRecorder struct {
	mock *MockVerifier
}

// NewMockVerifier creates a new mock instance.
func NewMockVerifier(ctrl *gomock.Controller) *MockVerifier {
	mock := &MockVerifier{ctrl: ctrl}
	mock.recorder = &MockVerifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVerifier) EXPECT() *MockVerifierMockRecorder {
	return m.recorder
}

// VerifySpend mocks base method.
func (m *MockVerifier) VerifySpend(arg0 txs.UnsignedTx, arg1 lux.UTXOGetter, arg2 []*lux.TransferableInput, arg3 []*lux.TransferableOutput, arg4 []verify.Verifiable, arg5 map[ids.ID]uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySpend", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifySpend indicates an expected call of VerifySpend.
func (mr *MockVerifierMockRecorder) VerifySpend(arg0, arg1, arg2, arg3, arg4, arg5 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySpend", reflect.TypeOf((*MockVerifier)(nil).VerifySpend), arg0, arg1, arg2, arg3, arg4, arg5)
}

// VerifySpendUTXOs mocks base method.
func (m *MockVerifier) VerifySpendUTXOs(arg0 txs.UnsignedTx, arg1 []*lux.UTXO, arg2 []*lux.TransferableInput, arg3 []*lux.TransferableOutput, arg4 []verify.Verifiable, arg5 map[ids.ID]uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifySpendUTXOs", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifySpendUTXOs indicates an expected call of VerifySpendUTXOs.
func (mr *MockVerifierMockRecorder) VerifySpendUTXOs(arg0, arg1, arg2, arg3, arg4, arg5 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifySpendUTXOs", reflect.TypeOf((*MockVerifier)(nil).VerifySpendUTXOs), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Code generated by MockGen. DO NOT EDIT.
// Source: vms/platformvm/txs/staker_tx.go
//
// Generated by this command:
//
//	mockgen -source=vms/platformvm/txs/staker_tx.go -destination=vms/platformvm/txs/mock_staker_tx.go -package=txs -exclude_interfaces=ValidatorTx,DelegatorTx,StakerTx,PermissionlessStaker
//

// Package txs is a generated GoMock package.
package txs

import (
	reflect "reflect"
	time "time"

	ids "github.com/skychains/chain/ids"
	bls "github.com/skychains/chain/utils/crypto/bls"
	gomock "go.uber.org/mock/gomock"
)

// MockStaker is a mock of Staker interface.
type MockStaker struct {
	ctrl     *gomock.Controller
	recorder *MockStakerMockRecorder
}

// MockStakerMockRecorder is the mock recorder for MockStaker.
type MockStakerMockRecorder struct {
	mock *MockStaker
}

// NewMockStaker creates a new mock instance.
func NewMockStaker(ctrl *gomock.Controller) *MockStaker {
	mock := &MockStaker{ctrl: ctrl}
	mock.recorder = &MockStakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStaker) EXPECT() *MockStakerMockRecorder {
	return m.recorder
}

// CurrentPriority mocks base method.
func (m *MockStaker) CurrentPriority() Priority {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentPriority")
	ret0, _ := ret[0].(Priority)
	return ret0
}

// CurrentPriority indicates an expected call of CurrentPriority.
func (mr *MockStakerMockRecorder) CurrentPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentPriority", reflect.TypeOf((*MockStaker)(nil).CurrentPriority))
}

// EndTime mocks base method.
func (m *MockStaker) EndTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// EndTime indicates an expected call of EndTime.
func (mr *MockStakerMockRecorder) EndTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndTime", reflect.TypeOf((*MockStaker)(nil).EndTime))
}

// NodeID mocks base method.
func (m *MockStaker) NodeID() ids.NodeID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodeID")
	ret0, _ := ret[0].(ids.NodeID)
	return ret0
}

// NodeID indicates an expected call of NodeID.
func (mr *MockStakerMockRecorder) NodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeID", reflect.TypeOf((*MockStaker)(nil).NodeID))
}

// PublicKey mocks base method.
func (m *MockStaker) PublicKey() (*bls.PublicKey, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicKey")
	ret0, _ := ret[0].(*bls.PublicKey)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PublicKey indicates an expected call of PublicKey.
func (mr *MockStakerMockRecorder) PublicKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicKey", reflect.TypeOf((*MockStaker)(nil).PublicKey))
}

// SubnetID mocks base method.
func (m *MockStaker) SubnetID() ids.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubnetID")
	ret0, _ := ret[0].(ids.ID)
	return ret0
}

// SubnetID indicates an expected call of SubnetID.
func (mr *MockStakerMockRecorder) SubnetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubnetID", reflect.TypeOf((*MockStaker)(nil).SubnetID))
}

// Weight mocks base method.
func (m *MockStaker) Weight() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Weight")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Weight indicates an expected call of Weight.
func (mr *MockStakerMockRecorder) Weight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Weight", reflect.TypeOf((*MockStaker)(nil).Weight))
}

// MockScheduledStaker is a mock of ScheduledStaker interface.
type MockScheduledStaker struct {
	ctrl     *gomock.Controller
	recorder *MockScheduledStakerMockRecorder
}

// MockScheduledStakerMockRecorder is the mock recorder for MockScheduledStaker.
type MockScheduledStakerMockRecorder struct {
	mock *MockScheduledStaker
}

// NewMockScheduledStaker creates a new mock instance.
func NewMockScheduledStaker(ctrl *gomock.Controller) *MockScheduledStaker {
	mock := &MockScheduledStaker{ctrl: ctrl}
	mock.recorder = &MockScheduledStakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduledStaker) EXPECT() *MockScheduledStakerMockRecorder {
	return m.recorder
}

// CurrentPriority mocks base method.
func (m *MockScheduledStaker) CurrentPriority() Priority {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentPriority")
	ret0, _ := ret[0].(Priority)
	return ret0
}

// CurrentPriority indicates an expected call of CurrentPriority.
func (mr *MockScheduledStakerMockRecorder) CurrentPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentPriority", reflect.TypeOf((*MockScheduledStaker)(nil).CurrentPriority))
}

// EndTime mocks base method.
func (m *MockScheduledStaker) EndTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// EndTime indicates an expected call of EndTime.
func (mr *MockScheduledStakerMockRecorder) EndTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndTime", reflect.TypeOf((*MockScheduledStaker)(nil).EndTime))
}

// NodeID mocks base method.
func (m *MockScheduledStaker) NodeID() ids.NodeID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NodeID")
	ret0, _ := ret[0].(ids.NodeID)
	return ret0
}

// NodeID indicates an expected call of NodeID.
func (mr *MockScheduledStakerMockRecorder) NodeID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NodeID", reflect.TypeOf((*MockScheduledStaker)(nil).NodeID))
}

// PendingPriority mocks base method.
func (m *MockScheduledStaker) PendingPriority() Priority {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PendingPriority")
	ret0, _ := ret[0].(Priority)
	return ret0
}

// PendingPriority indicates an expected call of PendingPriority.
func (mr *MockScheduledStakerMockRecorder) PendingPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PendingPriority", reflect.TypeOf((*MockScheduledStaker)(nil).PendingPriority))
}

// PublicKey mocks base method.
func (m *MockScheduledStaker) PublicKey() (*bls.PublicKey, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicKey")
	ret0, _ := ret[0].(*bls.PublicKey)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PublicKey indicates an expected call of PublicKey.
func (mr *MockScheduledStakerMockRecorder) PublicKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicKey", reflect.TypeOf((*MockScheduledStaker)(nil).PublicKey))
}

// StartTime mocks base method.
func (m *MockScheduledStaker) StartTime() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTime")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// StartTime indicates an expected call of StartTime.
func (mr *MockScheduledStakerMockRecorder) StartTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTime", reflect.TypeOf((*MockScheduledStaker)(nil).StartTime))
}

// SubnetID mocks base method.
func (m *MockScheduledStaker) SubnetID() ids.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubnetID")
	ret0, _ := ret[0].(ids.ID)
	return ret0
}

// SubnetID indicates an expected call of SubnetID.
func (mr *MockScheduledStakerMockRecorder) SubnetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubnetID", reflect.TypeOf((*MockScheduledStaker)(nil).SubnetID))
}

// Weight mocks base method.
func (m *MockScheduledStaker) Weight() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Weight")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Weight indicates an expected call of Weight.
func (mr *MockScheduledStakerMockRecorder) Weight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Weight", reflect.TypeOf((*MockScheduledStaker)(nil).Weight))
}

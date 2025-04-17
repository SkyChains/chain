// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SkyChains/chain/vms/proposervm/state (interfaces: State)
//
// Generated by this command:
//
//	mockgen -package=state -destination=vms/proposervm/state/mock_state.go github.com/SkyChains/chain/vms/proposervm/state State
//

// Package state is a generated GoMock package.
package state

import (
	reflect "reflect"

	ids "github.com/SkyChains/chain/ids"
	choices "github.com/SkyChains/chain/snow/choices"
	block "github.com/SkyChains/chain/vms/proposervm/block"
	gomock "go.uber.org/mock/gomock"
)

// MockState is a mock of State interface.
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState.
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance.
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// DeleteBlock mocks base method.
func (m *MockState) DeleteBlock(arg0 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBlock", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBlock indicates an expected call of DeleteBlock.
func (mr *MockStateMockRecorder) DeleteBlock(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBlock", reflect.TypeOf((*MockState)(nil).DeleteBlock), arg0)
}

// DeleteBlockIDAtHeight mocks base method.
func (m *MockState) DeleteBlockIDAtHeight(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBlockIDAtHeight", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBlockIDAtHeight indicates an expected call of DeleteBlockIDAtHeight.
func (mr *MockStateMockRecorder) DeleteBlockIDAtHeight(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBlockIDAtHeight", reflect.TypeOf((*MockState)(nil).DeleteBlockIDAtHeight), arg0)
}

// DeleteLastAccepted mocks base method.
func (m *MockState) DeleteLastAccepted() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLastAccepted")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteLastAccepted indicates an expected call of DeleteLastAccepted.
func (mr *MockStateMockRecorder) DeleteLastAccepted() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLastAccepted", reflect.TypeOf((*MockState)(nil).DeleteLastAccepted))
}

// GetBlock mocks base method.
func (m *MockState) GetBlock(arg0 ids.ID) (block.Block, choices.Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock", arg0)
	ret0, _ := ret[0].(block.Block)
	ret1, _ := ret[1].(choices.Status)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockStateMockRecorder) GetBlock(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockState)(nil).GetBlock), arg0)
}

// GetBlockIDAtHeight mocks base method.
func (m *MockState) GetBlockIDAtHeight(arg0 uint64) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockIDAtHeight", arg0)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockIDAtHeight indicates an expected call of GetBlockIDAtHeight.
func (mr *MockStateMockRecorder) GetBlockIDAtHeight(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockIDAtHeight", reflect.TypeOf((*MockState)(nil).GetBlockIDAtHeight), arg0)
}

// GetForkHeight mocks base method.
func (m *MockState) GetForkHeight() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForkHeight")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForkHeight indicates an expected call of GetForkHeight.
func (mr *MockStateMockRecorder) GetForkHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForkHeight", reflect.TypeOf((*MockState)(nil).GetForkHeight))
}

// GetLastAccepted mocks base method.
func (m *MockState) GetLastAccepted() (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastAccepted")
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastAccepted indicates an expected call of GetLastAccepted.
func (mr *MockStateMockRecorder) GetLastAccepted() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastAccepted", reflect.TypeOf((*MockState)(nil).GetLastAccepted))
}

// GetMinimumHeight mocks base method.
func (m *MockState) GetMinimumHeight() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMinimumHeight")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMinimumHeight indicates an expected call of GetMinimumHeight.
func (mr *MockStateMockRecorder) GetMinimumHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMinimumHeight", reflect.TypeOf((*MockState)(nil).GetMinimumHeight))
}

// PutBlock mocks base method.
func (m *MockState) PutBlock(arg0 block.Block, arg1 choices.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutBlock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutBlock indicates an expected call of PutBlock.
func (mr *MockStateMockRecorder) PutBlock(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutBlock", reflect.TypeOf((*MockState)(nil).PutBlock), arg0, arg1)
}

// SetBlockIDAtHeight mocks base method.
func (m *MockState) SetBlockIDAtHeight(arg0 uint64, arg1 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetBlockIDAtHeight", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetBlockIDAtHeight indicates an expected call of SetBlockIDAtHeight.
func (mr *MockStateMockRecorder) SetBlockIDAtHeight(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBlockIDAtHeight", reflect.TypeOf((*MockState)(nil).SetBlockIDAtHeight), arg0, arg1)
}

// SetForkHeight mocks base method.
func (m *MockState) SetForkHeight(arg0 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetForkHeight", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetForkHeight indicates an expected call of SetForkHeight.
func (mr *MockStateMockRecorder) SetForkHeight(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetForkHeight", reflect.TypeOf((*MockState)(nil).SetForkHeight), arg0)
}

// SetLastAccepted mocks base method.
func (m *MockState) SetLastAccepted(arg0 ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLastAccepted", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLastAccepted indicates an expected call of SetLastAccepted.
func (mr *MockStateMockRecorder) SetLastAccepted(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastAccepted", reflect.TypeOf((*MockState)(nil).SetLastAccepted), arg0)
}

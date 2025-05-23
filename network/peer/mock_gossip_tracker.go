// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/skychains/chain/network/peer (interfaces: GossipTracker)

// Package peer is a generated GoMock package.
package peer

import (
	reflect "reflect"

	ids "github.com/skychains/chain/ids"
	gomock "go.uber.org/mock/gomock"
)

// MockGossipTracker is a mock of GossipTracker interface.
type MockGossipTracker struct {
	ctrl     *gomock.Controller
	recorder *MockGossipTrackerMockRecorder
}

// MockGossipTrackerMockRecorder is the mock recorder for MockGossipTracker.
type MockGossipTrackerMockRecorder struct {
	mock *MockGossipTracker
}

// NewMockGossipTracker creates a new mock instance.
func NewMockGossipTracker(ctrl *gomock.Controller) *MockGossipTracker {
	mock := &MockGossipTracker{ctrl: ctrl}
	mock.recorder = &MockGossipTrackerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGossipTracker) EXPECT() *MockGossipTrackerMockRecorder {
	return m.recorder
}

// AddKnown mocks base method.
func (m *MockGossipTracker) AddKnown(arg0 ids.NodeID, arg1, arg2 []ids.ID) ([]ids.ID, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddKnown", arg0, arg1, arg2)
	ret0, _ := ret[0].([]ids.ID)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// AddKnown indicates an expected call of AddKnown.
func (mr *MockGossipTrackerMockRecorder) AddKnown(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddKnown", reflect.TypeOf((*MockGossipTracker)(nil).AddKnown), arg0, arg1, arg2)
}

// AddValidator mocks base method.
func (m *MockGossipTracker) AddValidator(arg0 ValidatorID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddValidator", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// AddValidator indicates an expected call of AddValidator.
func (mr *MockGossipTrackerMockRecorder) AddValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddValidator", reflect.TypeOf((*MockGossipTracker)(nil).AddValidator), arg0)
}

// GetNodeID mocks base method.
func (m *MockGossipTracker) GetNodeID(arg0 ids.ID) (ids.NodeID, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeID", arg0)
	ret0, _ := ret[0].(ids.NodeID)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetNodeID indicates an expected call of GetNodeID.
func (mr *MockGossipTrackerMockRecorder) GetNodeID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeID", reflect.TypeOf((*MockGossipTracker)(nil).GetNodeID), arg0)
}

// GetUnknown mocks base method.
func (m *MockGossipTracker) GetUnknown(arg0 ids.NodeID) ([]ValidatorID, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnknown", arg0)
	ret0, _ := ret[0].([]ValidatorID)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetUnknown indicates an expected call of GetUnknown.
func (mr *MockGossipTrackerMockRecorder) GetUnknown(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnknown", reflect.TypeOf((*MockGossipTracker)(nil).GetUnknown), arg0)
}

// RemoveValidator mocks base method.
func (m *MockGossipTracker) RemoveValidator(arg0 ids.NodeID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveValidator", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// RemoveValidator indicates an expected call of RemoveValidator.
func (mr *MockGossipTrackerMockRecorder) RemoveValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveValidator", reflect.TypeOf((*MockGossipTracker)(nil).RemoveValidator), arg0)
}

// ResetValidator mocks base method.
func (m *MockGossipTracker) ResetValidator(arg0 ids.NodeID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetValidator", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ResetValidator indicates an expected call of ResetValidator.
func (mr *MockGossipTrackerMockRecorder) ResetValidator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetValidator", reflect.TypeOf((*MockGossipTracker)(nil).ResetValidator), arg0)
}

// StartTrackingPeer mocks base method.
func (m *MockGossipTracker) StartTrackingPeer(arg0 ids.NodeID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTrackingPeer", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// StartTrackingPeer indicates an expected call of StartTrackingPeer.
func (mr *MockGossipTrackerMockRecorder) StartTrackingPeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTrackingPeer", reflect.TypeOf((*MockGossipTracker)(nil).StartTrackingPeer), arg0)
}

// StopTrackingPeer mocks base method.
func (m *MockGossipTracker) StopTrackingPeer(arg0 ids.NodeID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopTrackingPeer", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// StopTrackingPeer indicates an expected call of StopTrackingPeer.
func (mr *MockGossipTrackerMockRecorder) StopTrackingPeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopTrackingPeer", reflect.TypeOf((*MockGossipTracker)(nil).StopTrackingPeer), arg0)
}

// Tracked mocks base method.
func (m *MockGossipTracker) Tracked(arg0 ids.NodeID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Tracked", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Tracked indicates an expected call of Tracked.
func (mr *MockGossipTrackerMockRecorder) Tracked(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tracked", reflect.TypeOf((*MockGossipTracker)(nil).Tracked), arg0)
}

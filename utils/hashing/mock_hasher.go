// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/luxdefi/node/utils/hashing (interfaces: Hasher)

// Package hashing is a generated GoMock package.
package hashing

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHasher is a mock of Hasher interface.
type MockHasher struct {
	ctrl     *gomock.Controller
	recorder *MockHasherMockRecorder
}

// MockHasherMockRecorder is the mock recorder for MockHasher.
type MockHasherMockRecorder struct {
	mock *MockHasher
}

// NewMockHasher creates a new mock instance.
func NewMockHasher(ctrl *gomock.Controller) *MockHasher {
	mock := &MockHasher{ctrl: ctrl}
	mock.recorder = &MockHasherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHasher) EXPECT() *MockHasherMockRecorder {
	return m.recorder
}

// Hash mocks base method.
func (m *MockHasher) Hash(arg0 []byte) uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", arg0)
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Hash indicates an expected call of Hash.
func (mr *MockHasherMockRecorder) Hash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockHasher)(nil).Hash), arg0)
}

// Copyright (C) 2019-2023, Lux Partners Limited. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/skychains/chain/utils/logging (interfaces: Logger)

// Package logging is a generated GoMock package.
package logging

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	zapcore "go.uber.org/zap/zapcore"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), varargs...)
}

// Enabled mocks base method.
func (m *MockLogger) Enabled(arg0 Level) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enabled", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Enabled indicates an expected call of Enabled.
func (mr *MockLoggerMockRecorder) Enabled(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enabled", reflect.TypeOf((*MockLogger)(nil).Enabled), arg0)
}

// Error mocks base method.
func (m *MockLogger) Error(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), varargs...)
}

// Fatal mocks base method.
func (m *MockLogger) Fatal(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatal", varargs...)
}

// Fatal indicates an expected call of Fatal.
func (mr *MockLoggerMockRecorder) Fatal(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*MockLogger)(nil).Fatal), varargs...)
}

// Info mocks base method.
func (m *MockLogger) Info(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), varargs...)
}

// RecoverAndExit mocks base method.
func (m *MockLogger) RecoverAndExit(arg0, arg1 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RecoverAndExit", arg0, arg1)
}

// RecoverAndExit indicates an expected call of RecoverAndExit.
func (mr *MockLoggerMockRecorder) RecoverAndExit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecoverAndExit", reflect.TypeOf((*MockLogger)(nil).RecoverAndExit), arg0, arg1)
}

// RecoverAndPanic mocks base method.
func (m *MockLogger) RecoverAndPanic(arg0 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RecoverAndPanic", arg0)
}

// RecoverAndPanic indicates an expected call of RecoverAndPanic.
func (mr *MockLoggerMockRecorder) RecoverAndPanic(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecoverAndPanic", reflect.TypeOf((*MockLogger)(nil).RecoverAndPanic), arg0)
}

// SetLevel mocks base method.
func (m *MockLogger) SetLevel(arg0 Level) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLevel", arg0)
}

// SetLevel indicates an expected call of SetLevel.
func (mr *MockLoggerMockRecorder) SetLevel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLevel", reflect.TypeOf((*MockLogger)(nil).SetLevel), arg0)
}

// Stop mocks base method.
func (m *MockLogger) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockLoggerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockLogger)(nil).Stop))
}

// StopOnPanic mocks base method.
func (m *MockLogger) StopOnPanic() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopOnPanic")
}

// StopOnPanic indicates an expected call of StopOnPanic.
func (mr *MockLoggerMockRecorder) StopOnPanic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopOnPanic", reflect.TypeOf((*MockLogger)(nil).StopOnPanic))
}

// Trace mocks base method.
func (m *MockLogger) Trace(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Trace", varargs...)
}

// Trace indicates an expected call of Trace.
func (mr *MockLoggerMockRecorder) Trace(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trace", reflect.TypeOf((*MockLogger)(nil).Trace), varargs...)
}

// Verbo mocks base method.
func (m *MockLogger) Verbo(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Verbo", varargs...)
}

// Verbo indicates an expected call of Verbo.
func (mr *MockLoggerMockRecorder) Verbo(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verbo", reflect.TypeOf((*MockLogger)(nil).Verbo), varargs...)
}

// Warn mocks base method.
func (m *MockLogger) Warn(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockLoggerMockRecorder) Warn(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLogger)(nil).Warn), varargs...)
}

// Write mocks base method.
func (m *MockLogger) Write(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockLoggerMockRecorder) Write(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockLogger)(nil).Write), arg0)
}

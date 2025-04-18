// Code generated by MockGen. DO NOT EDIT.
// Source: x/merkledb/db.go
//
// Generated by this command:
//
//	mockgen -source=x/merkledb/db.go -destination=x/merkledb/mock_db.go -package=merkledb -exclude_interfaces=ChangeProofer,RangeProofer,Clearer,Prefetcher
//

// Package merkledb is a generated GoMock package.
package merkledb

import (
	context "context"
	reflect "reflect"

	database "github.com/skychains/chain/database"
	ids "github.com/skychains/chain/ids"
	maybe "github.com/skychains/chain/utils/maybe"
	gomock "go.uber.org/mock/gomock"
)

// MockMerkleDB is a mock of MerkleDB interface.
type MockMerkleDB struct {
	ctrl     *gomock.Controller
	recorder *MockMerkleDBMockRecorder
}

// MockMerkleDBMockRecorder is the mock recorder for MockMerkleDB.
type MockMerkleDBMockRecorder struct {
	mock *MockMerkleDB
}

// NewMockMerkleDB creates a new mock instance.
func NewMockMerkleDB(ctrl *gomock.Controller) *MockMerkleDB {
	mock := &MockMerkleDB{ctrl: ctrl}
	mock.recorder = &MockMerkleDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMerkleDB) EXPECT() *MockMerkleDBMockRecorder {
	return m.recorder
}

// Clear mocks base method.
func (m *MockMerkleDB) Clear() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear")
	ret0, _ := ret[0].(error)
	return ret0
}

// Clear indicates an expected call of Clear.
func (mr *MockMerkleDBMockRecorder) Clear() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockMerkleDB)(nil).Clear))
}

// Close mocks base method.
func (m *MockMerkleDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockMerkleDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMerkleDB)(nil).Close))
}

// CommitChangeProof mocks base method.
func (m *MockMerkleDB) CommitChangeProof(ctx context.Context, proof *ChangeProof) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitChangeProof", ctx, proof)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitChangeProof indicates an expected call of CommitChangeProof.
func (mr *MockMerkleDBMockRecorder) CommitChangeProof(ctx, proof any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitChangeProof", reflect.TypeOf((*MockMerkleDB)(nil).CommitChangeProof), ctx, proof)
}

// CommitRangeProof mocks base method.
func (m *MockMerkleDB) CommitRangeProof(ctx context.Context, start, end maybe.Maybe[[]byte], proof *RangeProof) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommitRangeProof", ctx, start, end, proof)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommitRangeProof indicates an expected call of CommitRangeProof.
func (mr *MockMerkleDBMockRecorder) CommitRangeProof(ctx, start, end, proof any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommitRangeProof", reflect.TypeOf((*MockMerkleDB)(nil).CommitRangeProof), ctx, start, end, proof)
}

// Compact mocks base method.
func (m *MockMerkleDB) Compact(start, limit []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compact", start, limit)
	ret0, _ := ret[0].(error)
	return ret0
}

// Compact indicates an expected call of Compact.
func (mr *MockMerkleDBMockRecorder) Compact(start, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compact", reflect.TypeOf((*MockMerkleDB)(nil).Compact), start, limit)
}

// Delete mocks base method.
func (m *MockMerkleDB) Delete(key []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockMerkleDBMockRecorder) Delete(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMerkleDB)(nil).Delete), key)
}

// Get mocks base method.
func (m *MockMerkleDB) Get(key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMerkleDBMockRecorder) Get(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMerkleDB)(nil).Get), key)
}

// GetChangeProof mocks base method.
func (m *MockMerkleDB) GetChangeProof(ctx context.Context, startRootID, endRootID ids.ID, start, end maybe.Maybe[[]byte], maxLength int) (*ChangeProof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChangeProof", ctx, startRootID, endRootID, start, end, maxLength)
	ret0, _ := ret[0].(*ChangeProof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChangeProof indicates an expected call of GetChangeProof.
func (mr *MockMerkleDBMockRecorder) GetChangeProof(ctx, startRootID, endRootID, start, end, maxLength any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChangeProof", reflect.TypeOf((*MockMerkleDB)(nil).GetChangeProof), ctx, startRootID, endRootID, start, end, maxLength)
}

// GetMerkleRoot mocks base method.
func (m *MockMerkleDB) GetMerkleRoot(ctx context.Context) (ids.ID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMerkleRoot", ctx)
	ret0, _ := ret[0].(ids.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMerkleRoot indicates an expected call of GetMerkleRoot.
func (mr *MockMerkleDBMockRecorder) GetMerkleRoot(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMerkleRoot", reflect.TypeOf((*MockMerkleDB)(nil).GetMerkleRoot), ctx)
}

// GetProof mocks base method.
func (m *MockMerkleDB) GetProof(ctx context.Context, keyBytes []byte) (*Proof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProof", ctx, keyBytes)
	ret0, _ := ret[0].(*Proof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProof indicates an expected call of GetProof.
func (mr *MockMerkleDBMockRecorder) GetProof(ctx, keyBytes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProof", reflect.TypeOf((*MockMerkleDB)(nil).GetProof), ctx, keyBytes)
}

// GetRangeProof mocks base method.
func (m *MockMerkleDB) GetRangeProof(ctx context.Context, start, end maybe.Maybe[[]byte], maxLength int) (*RangeProof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRangeProof", ctx, start, end, maxLength)
	ret0, _ := ret[0].(*RangeProof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRangeProof indicates an expected call of GetRangeProof.
func (mr *MockMerkleDBMockRecorder) GetRangeProof(ctx, start, end, maxLength any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRangeProof", reflect.TypeOf((*MockMerkleDB)(nil).GetRangeProof), ctx, start, end, maxLength)
}

// GetRangeProofAtRoot mocks base method.
func (m *MockMerkleDB) GetRangeProofAtRoot(ctx context.Context, rootID ids.ID, start, end maybe.Maybe[[]byte], maxLength int) (*RangeProof, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRangeProofAtRoot", ctx, rootID, start, end, maxLength)
	ret0, _ := ret[0].(*RangeProof)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRangeProofAtRoot indicates an expected call of GetRangeProofAtRoot.
func (mr *MockMerkleDBMockRecorder) GetRangeProofAtRoot(ctx, rootID, start, end, maxLength any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRangeProofAtRoot", reflect.TypeOf((*MockMerkleDB)(nil).GetRangeProofAtRoot), ctx, rootID, start, end, maxLength)
}

// GetValue mocks base method.
func (m *MockMerkleDB) GetValue(ctx context.Context, key []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", ctx, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValue indicates an expected call of GetValue.
func (mr *MockMerkleDBMockRecorder) GetValue(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockMerkleDB)(nil).GetValue), ctx, key)
}

// GetValues mocks base method.
func (m *MockMerkleDB) GetValues(ctx context.Context, keys [][]byte) ([][]byte, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValues", ctx, keys)
	ret0, _ := ret[0].([][]byte)
	ret1, _ := ret[1].([]error)
	return ret0, ret1
}

// GetValues indicates an expected call of GetValues.
func (mr *MockMerkleDBMockRecorder) GetValues(ctx, keys any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValues", reflect.TypeOf((*MockMerkleDB)(nil).GetValues), ctx, keys)
}

// Has mocks base method.
func (m *MockMerkleDB) Has(key []byte) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", key)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Has indicates an expected call of Has.
func (mr *MockMerkleDBMockRecorder) Has(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockMerkleDB)(nil).Has), key)
}

// HealthCheck mocks base method.
func (m *MockMerkleDB) HealthCheck(arg0 context.Context) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *MockMerkleDBMockRecorder) HealthCheck(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockMerkleDB)(nil).HealthCheck), arg0)
}

// NewBatch mocks base method.
func (m *MockMerkleDB) NewBatch() database.Batch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBatch")
	ret0, _ := ret[0].(database.Batch)
	return ret0
}

// NewBatch indicates an expected call of NewBatch.
func (mr *MockMerkleDBMockRecorder) NewBatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBatch", reflect.TypeOf((*MockMerkleDB)(nil).NewBatch))
}

// NewIterator mocks base method.
func (m *MockMerkleDB) NewIterator() database.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIterator")
	ret0, _ := ret[0].(database.Iterator)
	return ret0
}

// NewIterator indicates an expected call of NewIterator.
func (mr *MockMerkleDBMockRecorder) NewIterator() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIterator", reflect.TypeOf((*MockMerkleDB)(nil).NewIterator))
}

// NewIteratorWithPrefix mocks base method.
func (m *MockMerkleDB) NewIteratorWithPrefix(prefix []byte) database.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIteratorWithPrefix", prefix)
	ret0, _ := ret[0].(database.Iterator)
	return ret0
}

// NewIteratorWithPrefix indicates an expected call of NewIteratorWithPrefix.
func (mr *MockMerkleDBMockRecorder) NewIteratorWithPrefix(prefix any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIteratorWithPrefix", reflect.TypeOf((*MockMerkleDB)(nil).NewIteratorWithPrefix), prefix)
}

// NewIteratorWithStart mocks base method.
func (m *MockMerkleDB) NewIteratorWithStart(start []byte) database.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIteratorWithStart", start)
	ret0, _ := ret[0].(database.Iterator)
	return ret0
}

// NewIteratorWithStart indicates an expected call of NewIteratorWithStart.
func (mr *MockMerkleDBMockRecorder) NewIteratorWithStart(start any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIteratorWithStart", reflect.TypeOf((*MockMerkleDB)(nil).NewIteratorWithStart), start)
}

// NewIteratorWithStartAndPrefix mocks base method.
func (m *MockMerkleDB) NewIteratorWithStartAndPrefix(start, prefix []byte) database.Iterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewIteratorWithStartAndPrefix", start, prefix)
	ret0, _ := ret[0].(database.Iterator)
	return ret0
}

// NewIteratorWithStartAndPrefix indicates an expected call of NewIteratorWithStartAndPrefix.
func (mr *MockMerkleDBMockRecorder) NewIteratorWithStartAndPrefix(start, prefix any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewIteratorWithStartAndPrefix", reflect.TypeOf((*MockMerkleDB)(nil).NewIteratorWithStartAndPrefix), start, prefix)
}

// NewView mocks base method.
func (m *MockMerkleDB) NewView(ctx context.Context, changes ViewChanges) (View, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewView", ctx, changes)
	ret0, _ := ret[0].(View)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewView indicates an expected call of NewView.
func (mr *MockMerkleDBMockRecorder) NewView(ctx, changes any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewView", reflect.TypeOf((*MockMerkleDB)(nil).NewView), ctx, changes)
}

// PrefetchPath mocks base method.
func (m *MockMerkleDB) PrefetchPath(key []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrefetchPath", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrefetchPath indicates an expected call of PrefetchPath.
func (mr *MockMerkleDBMockRecorder) PrefetchPath(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrefetchPath", reflect.TypeOf((*MockMerkleDB)(nil).PrefetchPath), key)
}

// PrefetchPaths mocks base method.
func (m *MockMerkleDB) PrefetchPaths(keys [][]byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrefetchPaths", keys)
	ret0, _ := ret[0].(error)
	return ret0
}

// PrefetchPaths indicates an expected call of PrefetchPaths.
func (mr *MockMerkleDBMockRecorder) PrefetchPaths(keys any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrefetchPaths", reflect.TypeOf((*MockMerkleDB)(nil).PrefetchPaths), keys)
}

// Put mocks base method.
func (m *MockMerkleDB) Put(key, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockMerkleDBMockRecorder) Put(key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockMerkleDB)(nil).Put), key, value)
}

// VerifyChangeProof mocks base method.
func (m *MockMerkleDB) VerifyChangeProof(ctx context.Context, proof *ChangeProof, start, end maybe.Maybe[[]byte], expectedEndRootID ids.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyChangeProof", ctx, proof, start, end, expectedEndRootID)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyChangeProof indicates an expected call of VerifyChangeProof.
func (mr *MockMerkleDBMockRecorder) VerifyChangeProof(ctx, proof, start, end, expectedEndRootID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyChangeProof", reflect.TypeOf((*MockMerkleDB)(nil).VerifyChangeProof), ctx, proof, start, end, expectedEndRootID)
}

// getEditableNode mocks base method.
func (m *MockMerkleDB) getEditableNode(key Key, hasValue bool) (*node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getEditableNode", key, hasValue)
	ret0, _ := ret[0].(*node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getEditableNode indicates an expected call of getEditableNode.
func (mr *MockMerkleDBMockRecorder) getEditableNode(key, hasValue any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getEditableNode", reflect.TypeOf((*MockMerkleDB)(nil).getEditableNode), key, hasValue)
}

// getNode mocks base method.
func (m *MockMerkleDB) getNode(key Key, hasValue bool) (*node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getNode", key, hasValue)
	ret0, _ := ret[0].(*node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getNode indicates an expected call of getNode.
func (mr *MockMerkleDBMockRecorder) getNode(key, hasValue any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getNode", reflect.TypeOf((*MockMerkleDB)(nil).getNode), key, hasValue)
}

// getRoot mocks base method.
func (m *MockMerkleDB) getRoot() maybe.Maybe[*node] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getRoot")
	ret0, _ := ret[0].(maybe.Maybe[*node])
	return ret0
}

// getRoot indicates an expected call of getRoot.
func (mr *MockMerkleDBMockRecorder) getRoot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getRoot", reflect.TypeOf((*MockMerkleDB)(nil).getRoot))
}

// getTokenSize mocks base method.
func (m *MockMerkleDB) getTokenSize() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getTokenSize")
	ret0, _ := ret[0].(int)
	return ret0
}

// getTokenSize indicates an expected call of getTokenSize.
func (mr *MockMerkleDBMockRecorder) getTokenSize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getTokenSize", reflect.TypeOf((*MockMerkleDB)(nil).getTokenSize))
}

// getValue mocks base method.
func (m *MockMerkleDB) getValue(key Key) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getValue", key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getValue indicates an expected call of getValue.
func (mr *MockMerkleDBMockRecorder) getValue(key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getValue", reflect.TypeOf((*MockMerkleDB)(nil).getValue), key)
}

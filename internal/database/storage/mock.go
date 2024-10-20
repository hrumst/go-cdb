// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	os "os"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	compute "github.com/hrumst/go-cdb/internal/database/compute"
)

// MockstorageEngine is a mock of storageEngine interface.
type MockstorageEngine struct {
	ctrl     *gomock.Controller
	recorder *MockstorageEngineMockRecorder
}

// MockstorageEngineMockRecorder is the mock recorder for MockstorageEngine.
type MockstorageEngineMockRecorder struct {
	mock *MockstorageEngine
}

// NewMockstorageEngine creates a new mock instance.
func NewMockstorageEngine(ctrl *gomock.Controller) *MockstorageEngine {
	mock := &MockstorageEngine{ctrl: ctrl}
	mock.recorder = &MockstorageEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockstorageEngine) EXPECT() *MockstorageEngineMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockstorageEngine) Del(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockstorageEngineMockRecorder) Del(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockstorageEngine)(nil).Del), ctx, key)
}

// Get mocks base method.
func (m *MockstorageEngine) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockstorageEngineMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockstorageEngine)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockstorageEngine) Set(ctx context.Context, key, val string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockstorageEngineMockRecorder) Set(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockstorageEngine)(nil).Set), ctx, key, val)
}

// Mockwal is a mock of wal interface.
type Mockwal struct {
	ctrl     *gomock.Controller
	recorder *MockwalMockRecorder
}

// MockwalMockRecorder is the mock recorder for Mockwal.
type MockwalMockRecorder struct {
	mock *Mockwal
}

// NewMockwal creates a new mock instance.
func NewMockwal(ctrl *gomock.Controller) *Mockwal {
	mock := &Mockwal{ctrl: ctrl}
	mock.recorder = &MockwalMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockwal) EXPECT() *MockwalMockRecorder {
	return m.recorder
}

// AddLogRecord mocks base method.
func (m *Mockwal) AddLogRecord(ctx context.Context, cmdType compute.CommandType, args []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLogRecord", ctx, cmdType, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLogRecord indicates an expected call of AddLogRecord.
func (mr *MockwalMockRecorder) AddLogRecord(ctx, cmdType, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLogRecord", reflect.TypeOf((*Mockwal)(nil).AddLogRecord), ctx, cmdType, args)
}

// MockfsDir is a mock of fsDir interface.
type MockfsDir struct {
	ctrl     *gomock.Controller
	recorder *MockfsDirMockRecorder
}

// MockfsDirMockRecorder is the mock recorder for MockfsDir.
type MockfsDirMockRecorder struct {
	mock *MockfsDir
}

// NewMockfsDir creates a new mock instance.
func NewMockfsDir(ctrl *gomock.Controller) *MockfsDir {
	mock := &MockfsDir{ctrl: ctrl}
	mock.recorder = &MockfsDirMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockfsDir) EXPECT() *MockfsDirMockRecorder {
	return m.recorder
}

// FilesStats mocks base method.
func (m *MockfsDir) FilesStats() ([]os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilesStats")
	ret0, _ := ret[0].([]os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilesStats indicates an expected call of FilesStats.
func (mr *MockfsDirMockRecorder) FilesStats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilesStats", reflect.TypeOf((*MockfsDir)(nil).FilesStats))
}

// ReadFile mocks base method.
func (m *MockfsDir) ReadFile(filename string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFile", filename)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadFile indicates an expected call of ReadFile.
func (mr *MockfsDirMockRecorder) ReadFile(filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFile", reflect.TypeOf((*MockfsDir)(nil).ReadFile), filename)
}

// WriteSync mocks base method.
func (m *MockfsDir) WriteSync(filename string, data []byte) (os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteSync", filename, data)
	ret0, _ := ret[0].(os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteSync indicates an expected call of WriteSync.
func (mr *MockfsDirMockRecorder) WriteSync(filename, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteSync", reflect.TypeOf((*MockfsDir)(nil).WriteSync), filename, data)
}

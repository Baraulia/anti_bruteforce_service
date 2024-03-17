// Code generated by MockGen. DO NOT EDIT.
// Source: storageInterface.go

// Package mockstorage is a generated GoMock package.
package mockstorage

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddToBlackList mocks base method.
func (m *MockStorage) AddToBlackList(ctx context.Context, ip string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToBlackList", ctx, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToBlackList indicates an expected call of AddToBlackList.
func (mr *MockStorageMockRecorder) AddToBlackList(ctx, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToBlackList", reflect.TypeOf((*MockStorage)(nil).AddToBlackList), ctx, ip)
}

// AddToWhiteList mocks base method.
func (m *MockStorage) AddToWhiteList(ctx context.Context, ip string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToWhiteList", ctx, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddToWhiteList indicates an expected call of AddToWhiteList.
func (mr *MockStorageMockRecorder) AddToWhiteList(ctx, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToWhiteList", reflect.TypeOf((*MockStorage)(nil).AddToWhiteList), ctx, ip)
}

// CheckIPInBlackList mocks base method.
func (m *MockStorage) CheckIPInBlackList(ctx context.Context, ip string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIPInBlackList", ctx, ip)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIPInBlackList indicates an expected call of CheckIPInBlackList.
func (mr *MockStorageMockRecorder) CheckIPInBlackList(ctx, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIPInBlackList", reflect.TypeOf((*MockStorage)(nil).CheckIPInBlackList), ctx, ip)
}

// CheckIPInWhiteList mocks base method.
func (m *MockStorage) CheckIPInWhiteList(ctx context.Context, ip string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIPInWhiteList", ctx, ip)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIPInWhiteList indicates an expected call of CheckIPInWhiteList.
func (mr *MockStorageMockRecorder) CheckIPInWhiteList(ctx, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIPInWhiteList", reflect.TypeOf((*MockStorage)(nil).CheckIPInWhiteList), ctx, ip)
}

// RemoveFromBlackList mocks base method.
func (m *MockStorage) RemoveFromBlackList(ctx context.Context, ip string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromBlackList", ctx, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromBlackList indicates an expected call of RemoveFromBlackList.
func (mr *MockStorageMockRecorder) RemoveFromBlackList(ctx, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromBlackList", reflect.TypeOf((*MockStorage)(nil).RemoveFromBlackList), ctx, ip)
}

// RemoveFromWhiteList mocks base method.
func (m *MockStorage) RemoveFromWhiteList(ctx context.Context, ip string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromWhiteList", ctx, ip)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromWhiteList indicates an expected call of RemoveFromWhiteList.
func (mr *MockStorageMockRecorder) RemoveFromWhiteList(ctx, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromWhiteList", reflect.TypeOf((*MockStorage)(nil).RemoveFromWhiteList), ctx, ip)
}

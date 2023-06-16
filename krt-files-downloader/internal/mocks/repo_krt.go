// Code generated by MockGen. DO NOT EDIT.
// Source: krt.go

// Package mocks is a generated GoMock package.
package mocks

import (
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockKRTRepository is a mock of KRTRepository interface.
type MockKRTRepository struct {
	ctrl     *gomock.Controller
	recorder *MockKRTRepositoryMockRecorder
}

// MockKRTRepositoryMockRecorder is the mock recorder for MockKRTRepository.
type MockKRTRepositoryMockRecorder struct {
	mock *MockKRTRepository
}

// NewMockKRTRepository creates a new mock instance.
func NewMockKRTRepository(ctrl *gomock.Controller) *MockKRTRepository {
	mock := &MockKRTRepository{ctrl: ctrl}
	mock.recorder = &MockKRTRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKRTRepository) EXPECT() *MockKRTRepositoryMockRecorder {
	return m.recorder
}

// DownloadKRT mocks base method.
func (m *MockKRTRepository) DownloadKRT(runtimeID, versionID string) (io.Reader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadKRT", runtimeID, versionID)
	ret0, _ := ret[0].(io.Reader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadKRT indicates an expected call of DownloadKRT.
func (mr *MockKRTRepositoryMockRecorder) DownloadKRT(runtimeID, versionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadKRT", reflect.TypeOf((*MockKRTRepository)(nil).DownloadKRT), runtimeID, versionID)
}
// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/antsanchez/go-download-web/pkg/scraper (interfaces: HttpGet)
//
// Generated by this command:
//
//	mockgen -destination=pkg/get/mock_get.go -package=get github.com/antsanchez/go-download-web/pkg/scraper HttpGet
//

// Package get is a generated GoMock package.
package get

import (
	bytes "bytes"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHttpGet is a mock of HttpGet interface.
type MockHttpGet struct {
	ctrl     *gomock.Controller
	recorder *MockHttpGetMockRecorder
}

// MockHttpGetMockRecorder is the mock recorder for MockHttpGet.
type MockHttpGetMockRecorder struct {
	mock *MockHttpGet
}

// NewMockHttpGet creates a new mock instance.
func NewMockHttpGet(ctrl *gomock.Controller) *MockHttpGet {
	mock := &MockHttpGet{ctrl: ctrl}
	mock.recorder = &MockHttpGetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpGet) EXPECT() *MockHttpGetMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockHttpGet) Get(arg0 string) (string, int, *bytes.Buffer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(*bytes.Buffer)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// Get indicates an expected call of Get.
func (mr *MockHttpGetMockRecorder) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockHttpGet)(nil).Get), arg0)
}

// ParseURL mocks base method.
func (m *MockHttpGet) ParseURL(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseURL", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseURL indicates an expected call of ParseURL.
func (mr *MockHttpGetMockRecorder) ParseURL(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseURL", reflect.TypeOf((*MockHttpGet)(nil).ParseURL), arg0, arg1)
}

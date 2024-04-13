// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/antsanchez/go-download-web/pkg/scraper (interfaces: Console)
//
// Generated by this command:
//
//	mockgen -destination=pkg/console/mock_console.go -package=console github.com/antsanchez/go-download-web/pkg/scraper Console
//

// Package console is a generated GoMock package.
package console

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockConsole is a mock of Console interface.
type MockConsole struct {
	ctrl     *gomock.Controller
	recorder *MockConsoleMockRecorder
}

// MockConsoleMockRecorder is the mock recorder for MockConsole.
type MockConsoleMockRecorder struct {
	mock *MockConsole
}

// NewMockConsole creates a new mock instance.
func NewMockConsole(ctrl *gomock.Controller) *MockConsole {
	mock := &MockConsole{ctrl: ctrl}
	mock.recorder = &MockConsoleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsole) EXPECT() *MockConsoleMockRecorder {
	return m.recorder
}

// AddAttachments mocks base method.
func (m *MockConsole) AddAttachments() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddAttachments")
}

// AddAttachments indicates an expected call of AddAttachments.
func (mr *MockConsoleMockRecorder) AddAttachments() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAttachments", reflect.TypeOf((*MockConsole)(nil).AddAttachments))
}

// AddDomain mocks base method.
func (m *MockConsole) AddDomain(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddDomain", arg0)
}

// AddDomain indicates an expected call of AddDomain.
func (mr *MockConsoleMockRecorder) AddDomain(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDomain", reflect.TypeOf((*MockConsole)(nil).AddDomain), arg0)
}

// AddDownloaded mocks base method.
func (m *MockConsole) AddDownloaded() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddDownloaded")
}

// AddDownloaded indicates an expected call of AddDownloaded.
func (mr *MockConsoleMockRecorder) AddDownloaded() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDownloaded", reflect.TypeOf((*MockConsole)(nil).AddDownloaded))
}

// AddDownloading mocks base method.
func (m *MockConsole) AddDownloading() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddDownloading")
}

// AddDownloading indicates an expected call of AddDownloading.
func (mr *MockConsoleMockRecorder) AddDownloading() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDownloading", reflect.TypeOf((*MockConsole)(nil).AddDownloading))
}

// AddErrors mocks base method.
func (m *MockConsole) AddErrors(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddErrors", arg0)
}

// AddErrors indicates an expected call of AddErrors.
func (mr *MockConsoleMockRecorder) AddErrors(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddErrors", reflect.TypeOf((*MockConsole)(nil).AddErrors), arg0)
}

// AddFinished mocks base method.
func (m *MockConsole) AddFinished() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddFinished")
}

// AddFinished indicates an expected call of AddFinished.
func (mr *MockConsoleMockRecorder) AddFinished() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFinished", reflect.TypeOf((*MockConsole)(nil).AddFinished))
}

// AddStarted mocks base method.
func (m *MockConsole) AddStarted() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddStarted")
}

// AddStarted indicates an expected call of AddStarted.
func (mr *MockConsoleMockRecorder) AddStarted() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddStarted", reflect.TypeOf((*MockConsole)(nil).AddStarted))
}

// AddStatus mocks base method.
func (m *MockConsole) AddStatus(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddStatus", arg0)
}

// AddStatus indicates an expected call of AddStatus.
func (mr *MockConsoleMockRecorder) AddStatus(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddStatus", reflect.TypeOf((*MockConsole)(nil).AddStatus), arg0)
}

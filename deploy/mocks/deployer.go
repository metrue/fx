// Code generated by MockGen. DO NOT EDIT.
// Source: deployer.go

// Package mock_deploy is a generated GoMock package.
package mock_deploy

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	types "github.com/metrue/fx/types"
	reflect "reflect"
)

// MockDeployer is a mock of Deployer interface
type MockDeployer struct {
	ctrl     *gomock.Controller
	recorder *MockDeployerMockRecorder
}

// MockDeployerMockRecorder is the mock recorder for MockDeployer
type MockDeployerMockRecorder struct {
	mock *MockDeployer
}

// NewMockDeployer creates a new mock instance
func NewMockDeployer(ctrl *gomock.Controller) *MockDeployer {
	mock := &MockDeployer{ctrl: ctrl}
	mock.recorder = &MockDeployerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDeployer) EXPECT() *MockDeployerMockRecorder {
	return m.recorder
}

// Deploy mocks base method
func (m *MockDeployer) Deploy(ctx context.Context, fn types.Func, name, image string, bindings []types.PortBinding) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deploy", ctx, fn, name, image, bindings)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deploy indicates an expected call of Deploy
func (mr *MockDeployerMockRecorder) Deploy(ctx, fn, name, image, bindings interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deploy", reflect.TypeOf((*MockDeployer)(nil).Deploy), ctx, fn, name, image, bindings)
}

// Destroy mocks base method
func (m *MockDeployer) Destroy(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockDeployerMockRecorder) Destroy(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockDeployer)(nil).Destroy), ctx, name)
}

// Update mocks base method
func (m *MockDeployer) Update(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockDeployerMockRecorder) Update(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDeployer)(nil).Update), ctx, name)
}

// GetStatus mocks base method
func (m *MockDeployer) GetStatus(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetStatus indicates an expected call of GetStatus
func (mr *MockDeployerMockRecorder) GetStatus(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockDeployer)(nil).GetStatus), ctx, name)
}

// List mocks base method
func (m *MockDeployer) List(ctx context.Context, name string) ([]types.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, name)
	ret0, _ := ret[0].([]types.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockDeployerMockRecorder) List(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDeployer)(nil).List), ctx, name)
}

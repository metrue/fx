// Code generated by MockGen. DO NOT EDIT.
// Source: docker.go

// Package mock_driver is a generated GoMock package.
package mock_driver

import (
	context "context"
	types "github.com/docker/docker/api/types"
	container "github.com/docker/docker/api/types/container"
	network "github.com/docker/docker/api/types/network"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
	time "time"
)

// MockIDockerClient is a mock of IDockerClient interface
type MockIDockerClient struct {
	ctrl     *gomock.Controller
	recorder *MockIDockerClientMockRecorder
}

// MockIDockerClientMockRecorder is the mock recorder for MockIDockerClient
type MockIDockerClientMockRecorder struct {
	mock *MockIDockerClient
}

// NewMockIDockerClient creates a new mock instance
func NewMockIDockerClient(ctrl *gomock.Controller) *MockIDockerClient {
	mock := &MockIDockerClient{ctrl: ctrl}
	mock.recorder = &MockIDockerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIDockerClient) EXPECT() *MockIDockerClientMockRecorder {
	return m.recorder
}

// ImageBuild mocks base method
func (m *MockIDockerClient) ImageBuild(ctx context.Context, r io.Reader, opt types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	ret := m.ctrl.Call(m, "ImageBuild", ctx, r, opt)
	ret0, _ := ret[0].(types.ImageBuildResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageBuild indicates an expected call of ImageBuild
func (mr *MockIDockerClientMockRecorder) ImageBuild(ctx, r, opt interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageBuild", reflect.TypeOf((*MockIDockerClient)(nil).ImageBuild), ctx, r, opt)
}

// ContainerCreate mocks base method
func (m *MockIDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error) {
	ret := m.ctrl.Call(m, "ContainerCreate", ctx, config, hostConfig, networkingConfig, containerName)
	ret0, _ := ret[0].(container.ContainerCreateCreatedBody)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ContainerCreate indicates an expected call of ContainerCreate
func (mr *MockIDockerClientMockRecorder) ContainerCreate(ctx, config, hostConfig, networkingConfig, containerName interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerCreate", reflect.TypeOf((*MockIDockerClient)(nil).ContainerCreate), ctx, config, hostConfig, networkingConfig, containerName)
}

// ContainerStart mocks base method
func (m *MockIDockerClient) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	ret := m.ctrl.Call(m, "ContainerStart", ctx, containerID, options)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerStart indicates an expected call of ContainerStart
func (mr *MockIDockerClientMockRecorder) ContainerStart(ctx, containerID, options interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerStart", reflect.TypeOf((*MockIDockerClient)(nil).ContainerStart), ctx, containerID, options)
}

// ContainerStop mocks base method
func (m *MockIDockerClient) ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error {
	ret := m.ctrl.Call(m, "ContainerStop", ctx, containerID, timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// ContainerStop indicates an expected call of ContainerStop
func (mr *MockIDockerClientMockRecorder) ContainerStop(ctx, containerID, timeout interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContainerStop", reflect.TypeOf((*MockIDockerClient)(nil).ContainerStop), ctx, containerID, timeout)
}

// ImagePull mocks base method
func (m *MockIDockerClient) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	ret := m.ctrl.Call(m, "ImagePull", ctx, refStr, options)
	ret0, _ := ret[0].(io.ReadCloser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImagePull indicates an expected call of ImagePull
func (mr *MockIDockerClientMockRecorder) ImagePull(ctx, refStr, options interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImagePull", reflect.TypeOf((*MockIDockerClient)(nil).ImagePull), ctx, refStr, options)
}

// ImageList mocks base method
func (m *MockIDockerClient) ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error) {
	ret := m.ctrl.Call(m, "ImageList", ctx, options)
	ret0, _ := ret[0].([]types.ImageSummary)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ImageList indicates an expected call of ImageList
func (mr *MockIDockerClientMockRecorder) ImageList(ctx, options interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImageList", reflect.TypeOf((*MockIDockerClient)(nil).ImageList), ctx, options)
}

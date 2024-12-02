// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/actor.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "cinema/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockstoreActor is a mock of storeActor interface.
type MockstoreActor struct {
	ctrl     *gomock.Controller
	recorder *MockstoreActorMockRecorder
}

// MockstoreActorMockRecorder is the mock recorder for MockstoreActor.
type MockstoreActorMockRecorder struct {
	mock *MockstoreActor
}

// NewMockstoreActor creates a new mock instance.
func NewMockstoreActor(ctrl *gomock.Controller) *MockstoreActor {
	mock := &MockstoreActor{ctrl: ctrl}
	mock.recorder = &MockstoreActorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockstoreActor) EXPECT() *MockstoreActorMockRecorder {
	return m.recorder
}

// CreateActor mocks base method.
func (m *MockstoreActor) CreateActor(actor models.CreateActor) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActor", actor)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateActor indicates an expected call of CreateActor.
func (mr *MockstoreActorMockRecorder) CreateActor(actor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActor", reflect.TypeOf((*MockstoreActor)(nil).CreateActor), actor)
}

// DeleteActor mocks base method.
func (m *MockstoreActor) DeleteActor(id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActor", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActor indicates an expected call of DeleteActor.
func (mr *MockstoreActorMockRecorder) DeleteActor(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActor", reflect.TypeOf((*MockstoreActor)(nil).DeleteActor), id)
}

// GetActor mocks base method.
func (m *MockstoreActor) GetActor(id uuid.UUID) (map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActor", id)
	ret0, _ := ret[0].(map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActor indicates an expected call of GetActor.
func (mr *MockstoreActorMockRecorder) GetActor(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActor", reflect.TypeOf((*MockstoreActor)(nil).GetActor), id)
}

// GetActorsWithMovies mocks base method.
func (m *MockstoreActor) GetActorsWithMovies(limit, offset int) ([]map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorsWithMovies", limit, offset)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorsWithMovies indicates an expected call of GetActorsWithMovies.
func (mr *MockstoreActorMockRecorder) GetActorsWithMovies(limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorsWithMovies", reflect.TypeOf((*MockstoreActor)(nil).GetActorsWithMovies), limit, offset)
}

// GetAllActors mocks base method.
func (m *MockstoreActor) GetAllActors(limit, offset int) ([]map[string]interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllActors", limit, offset)
	ret0, _ := ret[0].([]map[string]interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllActors indicates an expected call of GetAllActors.
func (mr *MockstoreActorMockRecorder) GetAllActors(limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllActors", reflect.TypeOf((*MockstoreActor)(nil).GetAllActors), limit, offset)
}

// UpdateActor mocks base method.
func (m *MockstoreActor) UpdateActor(id uuid.UUID, actor models.UpdateActor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateActor", id, actor)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateActor indicates an expected call of UpdateActor.
func (mr *MockstoreActorMockRecorder) UpdateActor(id, actor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateActor", reflect.TypeOf((*MockstoreActor)(nil).UpdateActor), id, actor)
}

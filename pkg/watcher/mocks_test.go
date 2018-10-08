// Code generated by MockGen. DO NOT EDIT.
// Source: sigs.k8s.io/controller-runtime/pkg/reconcile (interfaces: Reconciler)

// Package watcher_test is a generated GoMock package.
package watcher_test

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// MockReconciler is a mock of Reconciler interface
type MockReconciler struct {
	ctrl     *gomock.Controller
	recorder *MockReconcilerMockRecorder
}

// MockReconcilerMockRecorder is the mock recorder for MockReconciler
type MockReconcilerMockRecorder struct {
	mock *MockReconciler
}

// NewMockReconciler creates a new mock instance
func NewMockReconciler(ctrl *gomock.Controller) *MockReconciler {
	mock := &MockReconciler{ctrl: ctrl}
	mock.recorder = &MockReconcilerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReconciler) EXPECT() *MockReconcilerMockRecorder {
	return m.recorder
}

// Reconcile mocks base method
func (m *MockReconciler) Reconcile(arg0 reconcile.Request) (reconcile.Result, error) {
	ret := m.ctrl.Call(m, "Reconcile", arg0)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reconcile indicates an expected call of Reconcile
func (mr *MockReconcilerMockRecorder) Reconcile(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reconcile", reflect.TypeOf((*MockReconciler)(nil).Reconcile), arg0)
}
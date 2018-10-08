//go:generate mockgen -destination=mocks_test.go -package=watcher_test sigs.k8s.io/controller-runtime/pkg/reconcile Reconciler

package watcher_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/MYOB-Technology/shell-operator/pkg/watcher"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestDeduplicateExecutesNormally(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inner := NewMockReconciler(mockCtrl)
	inner.EXPECT().Reconcile(gomock.Any()).Return(reconcile.Result{}, nil).Times(1)

	deduper := watcher.NewDeduplicateReconciler(inner, "prefix")
	req := reconcile.Request{types.NamespacedName{Name: "test-object", Namespace: "ns1"}}
	_, err := deduper.Reconcile(req)

	if err != nil {
		t.Error(err)
	}
}

func TestDeduplicatePassesErrorFromInner(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inner := NewMockReconciler(mockCtrl)
	inner.EXPECT().Reconcile(gomock.Any()).Return(reconcile.Result{}, errors.New("testing")).Times(1)

	deduper := watcher.NewDeduplicateReconciler(inner, "prefix")
	req := reconcile.Request{types.NamespacedName{Name: "test-object", Namespace: "ns1"}}
	_, err := deduper.Reconcile(req)

	if err == nil {
		t.Error(err)
		t.FailNow()
	}

	if err.Error() != "testing" {
		t.Error("expecting different error", err.Error())
	}
}

func TestDeduplicateFiltersMultipleCalls(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inner := NewMockReconciler(mockCtrl)
	inner.EXPECT().Reconcile(gomock.Any()).DoAndReturn(func(_ reconcile.Request) (reconcile.Result, error) {
		time.Sleep(1 * time.Second)
		return reconcile.Result{}, nil
	}).Times(2)

	deduper := watcher.NewDeduplicateReconciler(inner, "prefix")

	go func() {
		req := reconcile.Request{types.NamespacedName{Name: "test-object1", Namespace: "ns1"}}
		_, err := deduper.Reconcile(req)

		if err != nil {
			t.Error(err)
		}
	}()

	go func() {
		req := reconcile.Request{types.NamespacedName{Name: "test-object2", Namespace: "ns1"}}
		_, err := deduper.Reconcile(req)

		if err != nil {
			t.Error(err)
		}
	}()

	req := reconcile.Request{types.NamespacedName{Name: "test-object1", Namespace: "ns1"}}
	_, err := deduper.Reconcile(req)

	if err != nil {
		t.Error(err)
	}
}
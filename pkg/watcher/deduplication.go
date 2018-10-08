package watcher

import (
	"fmt"
	"sync"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func NewDeduplicateReconciler(inner reconcile.Reconciler, prefix string) *DeduplicateReconciler {
	return &DeduplicateReconciler{
		inner:        inner,
		syncMap:      sync.Map{},
		dupKeyPrefix: prefix,
	}
}

type DeduplicateReconciler struct {
	inner        reconcile.Reconciler
	syncMap      sync.Map
	dupKeyPrefix string
}

func (d *DeduplicateReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	var err error

	uniq := GetUniqueRequestKey(d.dupKeyPrefix, req)
	_, loaded := d.syncMap.LoadOrStore(uniq, "")

	res, err := reconcile.Result{}, nil

	if !loaded {
		res, err = d.inner.Reconcile(req)
	}

	d.syncMap.Delete(uniq)

	return res, err
}

func GetUniqueRequestKey(prefix string, req reconcile.Request) string {
	return fmt.Sprintf("%s-%s-%s", prefix, req.Namespace, req.Name)
}

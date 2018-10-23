package watcher

import (
	"testing"

	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/types"
)

func TestUniqKeyQueueDedupsOnObject(t *testing.T) {
	inner := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	q := NewUniqKeyQueue(inner)
	defer q.ShutDown()

	a := newReconcileRequest("a", "1")
	b := newReconcileRequest("a", "1")
	c := newReconcileRequest("c", "1")

	for i := 1; i <= 100; i++ {
		q.Add(a)
		q.Add(b)
	}

	q.Add(c)

	first, _ := q.Get()
	second, _ := q.Get()

	if first.(reconcile.Request).Namespace != "a" {
		t.Error("incorrect first request", first)
		t.FailNow()
	}

	if second.(reconcile.Request).Namespace != "c" {
		t.Error("incorrect second request", second)
		t.FailNow()
	}

	q.Forget(first)
	q.Forget(second)

	if q.Len() != 0 {
		t.Error("items in queue when should be empty")
	}
}

func TestQueueNormalOperationForNonRequestObj(t *testing.T) {
	inner := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	q := NewUniqKeyQueue(inner)
	defer q.ShutDown()

	q.Add(1)

	myInt, _ := q.Get()

	if myInt.(int) != 1 {
		t.Error("number not gotten from queue", myInt)
	}

	q.Forget(myInt)

	if q.Len() != 0 {
		t.Error("items in queue when should be empty")
	}
}

func newReconcileRequest(namespace, name string) reconcile.Request {
	return reconcile.Request{
		types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
	}
}

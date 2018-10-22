package watcher

import (
	"testing"

	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/types"
)

func TestRawQueueDedupsOnObjectValue(t *testing.T) {
	q := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	defer q.ShutDown()

	for i := 1; i <= 100; i++ {
		q.Add("test")
	}

	q.Add(2)

	first, _ := q.Get()
	second, _ := q.Get()

	if first.(string) != "test" {
		t.Error("first message in queue was", first)
	}

	if second.(int) != 2 {
		t.Error("first message was not deduped, got", second)
	}
}

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

	q.Done(first)

	q.Add(first)

	first, _ = q.Get()

	if first.(reconcile.Request).Namespace != "a" {
		t.Error("incorrect first request, the second time", first)
		t.FailNow()
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

	q.Done(myInt)
}

func newReconcileRequest(namespace, name string) reconcile.Request {
	return reconcile.Request{
		types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
	}
}

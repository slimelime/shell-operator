package watcher

import (
	"fmt"
	"sync"

	"github.com/golang/glog"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type UniqKeyQueue struct {
	workqueue.RateLimitingInterface

	uniqCond       *sync.Cond
	currentRunning map[string]struct{}
}

func (u *UniqKeyQueue) Add(obj interface{}) {
	u.uniqCond.L.Lock()
	defer u.uniqCond.L.Unlock()

	req, ok := obj.(reconcile.Request)

	if !ok {
		u.RateLimitingInterface.Add(obj)
		return
	}

	uniq := uniqStringFromRequest(req)

	if _, exists := u.currentRunning[uniq]; exists {
		glog.V(4).Infof("Already running req for %v, so discarding new event.", req)
		return
	}

	u.currentRunning[uniq] = struct{}{}

	u.RateLimitingInterface.Add(obj)
}

func (u *UniqKeyQueue) Forget(obj interface{}) {
	u.uniqCond.L.Lock()
	defer u.uniqCond.L.Unlock()

	req, ok := obj.(reconcile.Request)

	if !ok {
		u.RateLimitingInterface.Done(obj)
		return
	}

	uniq := uniqStringFromRequest(req)

	delete(u.currentRunning, uniq)
	u.RateLimitingInterface.Forget(obj)
}

func NewUniqKeyQueue(inner workqueue.RateLimitingInterface) *UniqKeyQueue {
	return &UniqKeyQueue{
		uniqCond:              sync.NewCond(&sync.Mutex{}),
		currentRunning:        map[string]struct{}{},
		RateLimitingInterface: inner,
	}
}

func uniqStringFromRequest(req reconcile.Request) string {
	return fmt.Sprintf("%s-%s", req.Namespace, req.Name)
}

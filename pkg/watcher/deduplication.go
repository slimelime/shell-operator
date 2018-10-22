package watcher

import (
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/client-go/util/workqueue"
)

type SourceDedupeDecorator struct {
	source.Kind
}

func (d *SourceDedupeDecorator) Start(h handler.EventHandler, q workqueue.RateLimitingInterface, p ...predicate.Predicate) error {
	wrappedQ := NewUniqKeyQueue(q)
	return d.Kind.Start(h, wrappedQ, p...)
}

package main

import (
	"log"


	k8sscheme "k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
)


func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		log.Fatal(err)
	}

	sc := k8sscheme.Scheme

	do, _ := CreateAndRegisterWatchObject(sc, "v1", "Pod")

	mgr, err := manager.New(cfg, manager.Options{Scheme: sc})

	if err != nil {
		panic(err)
	}

	c, err := controller.New("shell-controller", mgr,
		controller.Options{
			MaxConcurrentReconciles: 5,
			Reconciler: &ShellController{
				Client: mgr.GetClient(),
			},
		},
	)

	if err != nil {
		panic(err)
	}

	err = c.Watch(
		&source.Kind{Type: do},
		&handler.EnqueueRequestForObject{},
	)

	if err != nil {
		panic(err)
	}

	// Start the Cmd
	log.Fatal(mgr.Start(signals.SetupSignalHandler()))
}

type ShellController struct {
	client.Client
}

func (s *ShellController) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	log.Printf("Reconcile %v\n", req)
	return reconcile.Result{}, nil
}

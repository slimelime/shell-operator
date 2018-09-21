package watcher

import (
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/MYOB-Technology/shell-operator/pkg/executor"
	"github.com/golang/glog"
)

const (
	NameEnvVarKey      = "SHOP_OBJECT_NAME"
	NamespaceEnvVarKey = "SHOP_OBJECT_NAMESPACE"
)

type ShellReconciler struct {
	Command string
}

func (s *ShellReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	glog.V(4).Infof("Received reconcile request for %s/%s.", req.Namespace, req.Name)

	glog.V(4).Infof("Setting up shell command %s for %s/%s...", s.Command, req.Namespace, req.Name)
	cmd := executor.SetupShellCommand(s.Command, map[string]string{
		NameEnvVarKey:      req.Name,
		NamespaceEnvVarKey: req.Namespace,
	})

	glog.V(4).Infof("Running command %s...", s.Command)
	_, err := cmd.CombinedOutput()

	if err != nil {
		glog.V(4).Infof("Command %s failed: %s.", s.Command, err.Error())
		return reconcile.Result{Requeue: true}, nil
	}

	glog.V(4).Infof("Command %s successful.", s.Command)
	return reconcile.Result{Requeue: false}, nil
}

// SetupWatches takes the kubebuilder manager and will setup a kubebuilder controller for
// each watch item in the ShellConfig. It will use the Shell Reconciler to run the fn
// for every item that comes through.
func SetupWatches(mgr manager.Manager, shellConfig *ShellConfig) error {
	for _, watch := range shellConfig.Watch {
		glog.V(1).Infof("Setting up watch for %s/%s", watch.ApiVersion, watch.Kind)

		c, err := controller.New("shell-controller", mgr,
			controller.Options{
				MaxConcurrentReconciles: watch.Concurrency,
				Reconciler: &ShellReconciler{
					Command: watch.Command,
				},
			},
		)

		if err != nil {
			return err
		}

		do, err := CreateAndRegisterWatchObject(mgr.GetScheme(), watch.ApiVersion, watch.Kind)

		if err != nil {
			return err
		}

		c.Watch(&source.Kind{Type: do}, &handler.EnqueueRequestForObject{})
	}

	return nil
}

package watcher

import (
	"context"
	"os"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/MYOB-Technology/shell-operator/pkg/config"
	"github.com/MYOB-Technology/shell-operator/pkg/dynamic"
	"github.com/MYOB-Technology/shell-operator/pkg/shell"
	"github.com/golang/glog"
)

const (
	NameEnvVarKey       = "SHOP_OBJECT_NAME"
	NamespaceEnvVarKey  = "SHOP_OBJECT_NAMESPACE"
	APIVersionEnvVarKey = "SHOP_API_VERSION"
	KindEnvVarKey       = "SHOP_KIND"
)

type ShellReconciler struct {
	Command          string
	ObjectKind       string
	ObjectApiVersion string
	Timeout          time.Duration
	Environment      map[string]string
	ctx              context.Context
}

func (s *ShellReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	glog.V(4).Infof("Received reconcile request for %s", req)

	glog.V(4).Infof("Setting up shell command for %s, %s", s.ObjectKind, req)
	ctx, cancel := context.WithTimeout(s.ctx, s.Timeout*time.Second)
	defer cancel()

	env := map[string]string{
		NameEnvVarKey:       req.Name,
		NamespaceEnvVarKey:  req.Namespace,
		APIVersionEnvVarKey: s.ObjectApiVersion,
		KindEnvVarKey:       s.ObjectKind,
	}

	cmd := shell.New(ctx, s.Command)
	shell.AddEnvironment(cmd, env)
	shell.AddEnvironment(cmd, s.Environment)

	glog.V(4).Infof("Running command for %s %s", s.ObjectKind, req)
	err := shell.RunWithProgress(cmd, os.Stdout, os.Stderr)

	if err != nil {
		glog.V(4).Infof("Command for %s %s failed: %s.", s.ObjectKind, req, err.Error())
		return reconcile.Result{Requeue: true}, nil
	}

	glog.V(4).Infof("Command for %s %s successful.", s.ObjectKind, req)
	return reconcile.Result{Requeue: false}, nil
}

// SetupWatches takes the kubebuilder manager and will setup a kubebuilder controller for
// each watch item in the ShellConfig. It will use the Shell Reconciler to run the fn
// for every item that comes through.
func SetupWatches(ctx context.Context, mgr manager.Manager, shellConfig *config.ShellConfig) error {
	for _, watch := range shellConfig.Watch {
		glog.V(1).Infof("Setting up watch for %s/%s", watch.ApiVersion, watch.Kind)

		c, err := controller.New("shell-controller", mgr,
			controller.Options{
				MaxConcurrentReconciles: watch.Concurrency,
				Reconciler: &ShellReconciler{
					Command:          watch.Command,
					ObjectApiVersion: watch.ApiVersion,
					ObjectKind:       watch.Kind,
					Timeout:          time.Duration(watch.Timeout),
					Environment:      watch.Environment,
					ctx:              ctx,
				},
			},
		)

		if err != nil {
			return err
		}

		do, err := dynamic.CreateAndRegisterWatchObject(mgr.GetScheme(), watch.ApiVersion, watch.Kind)

		if err != nil {
			return err
		}

		c.Watch(&source.Kind{Type: do}, &handler.EnqueueRequestForObject{})
	}

	return nil
}

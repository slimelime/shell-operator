package watcher

import (
	"strings"
	"os"
	"os/exec"
	"log"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/handler"
)

const (
	NameEnvVarKey = "SHOP_OBJECT_NAME"
	NamespaceEnvVarKey = "SHOP_OBJECT_NAMESPACE"
)

type ShellReconciler struct {
	Command string
}

func (s *ShellReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	if err := runCommandForReconcile(s.Command, req.Name, req.Namespace); err != nil {
		log.Println("Error", err)
		return reconcile.Result{Requeue: true}, nil
	}

	return reconcile.Result{Requeue: false}, nil
}

func runCommandForReconcile(shellCommand, name, namespace string) error {
	cmds := strings.Split(shellCommand, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("%s=%s", NameEnvVarKey, name),
		fmt.Sprintf("%s=%s", NamespaceEnvVarKey, namespace),
	)
	out, err := cmd.CombinedOutput()

	log.Println("output:", string(out))
	return err
}

// SetupWatches takes the kubebuilder manager and will setup a kubebuilder controller for
// each watch item in the ShellConfig. It will use the Shell Reconciler to run the fn
// for every item that comes through.
func SetupWatches(mgr manager.Manager, shellConfig *ShellConfig) error {
	for _, watch := range shellConfig.Watch {
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

		c.Watch(&source.Kind{Type: do}, &handler.EnqueueRequestForObject{},)
	}

	return nil
}

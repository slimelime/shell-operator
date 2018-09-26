package executor

import (
	"context"
	"os/exec"
	"time"

	"github.com/MYOB-Technology/shell-operator/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	NameEnvVarKey       = "SHOP_OBJECT_NAME"
	NamespaceEnvVarKey  = "SHOP_OBJECT_NAMESPACE"
	APIVersionEnvVarKey = "SHOP_API_VERSION"
	KindEnvVarKey       = "SHOP_KIND"
)

type ExecFunc func(*exec.Cmd) error

type ShellReconciler struct {
	watch  config.Watch
	execFn ExecFunc
	ctx    context.Context
}

func NewShellReconciler(ctx context.Context, watch config.Watch, execFn ExecFunc) *ShellReconciler {
	return &ShellReconciler{
		watch:  watch,
		ctx:    ctx,
		execFn: execFn,
	}
}

func (s *ShellReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Duration(s.watch.Timeout)*time.Second)
	defer cancel()

	theEnv := map[string]string{}
	for k, v := range s.watch.Environment {
		theEnv[k] = v
	}

	theEnv[NamespaceEnvVarKey] = req.Namespace
	theEnv[NameEnvVarKey] = req.Name
	theEnv[APIVersionEnvVarKey] = s.watch.ApiVersion
	theEnv[KindEnvVarKey] = s.watch.Kind

	cmd := SetupShellCommand(ctx, s.watch.Command, theEnv)

	err := s.execFn(cmd)

	return reconcile.Result{}, err
}

package executor

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"testing"

	"github.com/MYOB-Technology/shell-operator/pkg/config"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func failIfEnvVarMissing(t *testing.T, key, value string, env []string) {
	found := false
	for _, v := range env {
		if v == fmt.Sprintf("%s=%s", key, value) {
			found = true
		}
	}

	if !found {
		t.Error("missing env var", key, value)
	}
}

func TestReconcilerHappy(t *testing.T) {
	w := config.Watch{
		Command:    "echo hello world",
		Kind:       "Pod",
		ApiVersion: "v1",
		Environment: map[string]string{
			"ENV1": "value1",
		},
	}

	funcCalled := 0

	recon := NewShellReconciler(context.Background(), w, func(cmd *exec.Cmd) error {
		funcCalled++
		if cmd.Args[0] != "echo" || cmd.Args[1] != "hello" || cmd.Args[2] != "world" {
			t.Error("incorrect command called:", cmd.Args)
		}

		failIfEnvVarMissing(t, "SHOP_OBJECT_NAMESPACE", "ns1", cmd.Env)
		failIfEnvVarMissing(t, "SHOP_OBJECT_NAME", "test-object", cmd.Env)
		failIfEnvVarMissing(t, "SHOP_KIND", "Pod", cmd.Env)
		failIfEnvVarMissing(t, "SHOP_API_VERSION", "v1", cmd.Env)
		failIfEnvVarMissing(t, "ENV1", "value1", cmd.Env)

		return nil
	})

	req := reconcile.Request{types.NamespacedName{Name: "test-object", Namespace: "ns1"}}
	_, err := recon.Reconcile(req)

	if err != nil {
		t.Error(err)
	}

	if funcCalled != 1 {
		t.Error("exec func not called only once, but called:", funcCalled)
	}
}

func TestReconcileErrorPropogation(t *testing.T) {
	w := config.Watch{
		Command:    "echo hello world",
		Kind:       "Pod",
		ApiVersion: "v1",
	}

	recon := NewShellReconciler(context.Background(), w, func(cmd *exec.Cmd) error {
		return errors.New("testing")
	})

	req := reconcile.Request{types.NamespacedName{Name: "test-object", Namespace: "ns1"}}
	_, err := recon.Reconcile(req)

	if err == nil {
		t.Error("exec func error not propogated")
		t.FailNow()
	}

	if err.Error() != "testing" {
		t.Error("unexpect error:", err)
	}
}

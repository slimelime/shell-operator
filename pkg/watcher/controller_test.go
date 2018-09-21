package watcher

import (
	"testing"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcilerRunsShellCmd(t *testing.T) {
	recon := ShellReconciler{Command: "bash -c echo $SHOP_OBJECT_NAME $SHOP_OBJECT_NAMESPACE > /tmp/shop-test1"}
	req := reconcile.Request{types.NamespacedName{Name: "test-object", Namespace: "ns1"}}

	res, err := recon.Reconcile(req)

	if err != nil {
		t.Error(err)
	}

	if res.Requeue {
		t.Error("Wanting to requeue when should have succeeded.")
	}
}

func TestReconcilerFailedShellCommand(t *testing.T) {
	recon := ShellReconciler{Command: "exit 1"}
	req := reconcile.Request{types.NamespacedName{Name: "test-object", Namespace: "ns1"}}

	res, err := recon.Reconcile(req)

	if err != nil {
		t.Error(err)
	}

	if !res.Requeue {
		t.Error("Not wanting to requeue when it should have failed.")
	}
}

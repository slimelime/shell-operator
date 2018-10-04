package dynamic

import (
	"reflect"
	"testing"

	"k8s.io/client-go/kubernetes/scheme"
)

func TestCreatingDynamicCRD(t *testing.T) {
	s := scheme.Scheme

	do, err := CreateAndRegisterWatchObject(s, "test.crd.io/v1", "MyObject")

	if err != nil {
		t.Error(err)
	}

	acKind := do.GetObjectKind().GroupVersionKind().Kind

	if acKind != "MyObject" {
		t.Error("Wrong kind:", acKind)
	}

	typeName := reflect.TypeOf(do).Elem().Name()
	if typeName != "DynamicObject" {
		t.Error("Unexpected struct type returned", typeName)
	}
}

func TestCreatingDynamick8sObj(t *testing.T) {
	s := scheme.Scheme

	do, err := CreateAndRegisterWatchObject(s, "v1", "Pod")

	if err != nil {
		t.Error(err)
	}

	typeName := reflect.TypeOf(do).Elem().Name()
	if typeName != "Pod" {
		t.Error("Unexpected struct type returned", typeName)
	}
}

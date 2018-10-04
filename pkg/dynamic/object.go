package dynamic

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DynamicObject represents ANY CRD that Kubernetes does not aleady know about.
type DynamicObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// DeepCopyObject is a function to satisfy the runtime.Object interface. This is needed to
// ensure we can use the DynamicObject to represent any CRD as if it were a compiled in struct
func (d *DynamicObject) DeepCopyObject() runtime.Object {
	return &DynamicObject{
		metav1.TypeMeta{Kind: d.Kind, APIVersion: d.APIVersion},
		metav1.ObjectMeta{Name: d.Name, Namespace: d.Namespace},
	}
}

// NewDynamicObject will take a Kubernetes 'kind' (the value in the 'kind' key in your yaml) and
// an api version (the 'apiVersion' key in your yaml) and generate a Go struct that looks like a
// compiled go struct the Kubernetes client recognises.
// This should only be used for CRDs, all native object are loaded into the scheme already for use.
// Example:
//   apiVersion=certmanager.k8s.io/v1alpha1, kind=Certificate
func NewDynamicObject(kind string, apiVersion string) *DynamicObject {
	return &DynamicObject{
		metav1.TypeMeta{Kind: kind, APIVersion: apiVersion},
		metav1.ObjectMeta{},
	}
}

// DynamicObjectList represents of list of 0 or more DynamicObjects
type DynamicObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []DynamicObject `json:"items"`
}

// DeepCopyObject is required by the runtime.Object interface to be able to act like a normal
// Kubernetes Object.
func (d *DynamicObjectList) DeepCopyObject() runtime.Object {
	return &DynamicObjectList{
		metav1.TypeMeta{Kind: d.Kind, APIVersion: d.APIVersion},
		metav1.ListMeta{},
		d.Items,
	}
}

// NewDynamicObjectList is used to represent the listing version of a DynamicObject. This is required
// as the controll will be trying to list objects as well as watch for individual changes on object as
// part of the reconciliation process.
// The kind and apiVersion are the same as a NewDynamicObject except they are always suffixed with 'List'.
// Examples:
//   apiVersion=certmanager.k8s.io/v1alpha1, kind=CertificateList
func NewDynamicObjectList(kind string, apiVersion string) *DynamicObjectList {
	return &DynamicObjectList{
		metav1.TypeMeta{Kind: kind, APIVersion: apiVersion},
		metav1.ListMeta{},
		[]DynamicObject{},
	}
}

// RegisterDynamicObjects will take the Kubernetes Scheme which houses all known Kubernetes Objects that can
// be manipulated and handles version across different api versions, and a list of any Kubernetes object that
// we want registered. This is specically for any CRDs as represented by DynamicObjects or DynamicObjectLists.
// This function will panic for any issues with the objects as there is reflection going onto do runtime checks
// rather than compile time checks underneath in the Kubernetes scheme struct.
func RegisterDynamicObjects(s *runtime.Scheme, dos ...runtime.Object) {
	for _, do := range dos {
		s.AddKnownTypeWithName(do.GetObjectKind().GroupVersionKind(), do)
		metav1.AddToGroupVersion(s, do.GetObjectKind().GroupVersionKind().GroupVersion())
	}
}

// CreateAndRegisterWatchObject will handle newing up a Kubernetes runtime.Object and its list equivalent as well
// as register it with the Kubernetes scheme so everything is wired up and ready to watch. It will distuiguish between
// CRDs and known k8s Objects and return the correct item once registered.
func CreateAndRegisterWatchObject(s *runtime.Scheme, apiVersion, kind string) (runtime.Object, error) {
	var obj runtime.Object

	// lookup type in scheme for gvk to see if its aleady known for k8s objects
	gv, err := schema.ParseGroupVersion(apiVersion)

	if err != nil {
		return nil, err
	}

	gvk := schema.GroupVersionKind{Kind: kind, Group: gv.Group, Version: gv.Version}

	if s.Recognizes(gvk) {
		obj, err = s.New(gvk)

		if err != nil {
			return nil, err
		}
	} else {
		obj = NewDynamicObject(kind, apiVersion)
		dol := NewDynamicObjectList(fmt.Sprintf("%sList", kind), apiVersion)

		RegisterDynamicObjects(s, obj, dol)
	}

	return obj, nil
}

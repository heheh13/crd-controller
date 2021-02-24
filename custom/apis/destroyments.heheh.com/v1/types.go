package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//_ "k8s.io/code-generator"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Destroyment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DestroymentSpec   `json:"spec,omitempty"`
	Status            DestroymentStatus `json:"status,omitempty"'`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DestroymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Destroyments    []Destroyment `json:"destroyments"`
}

type DestroymentSpec struct {
	//may  be need to define as a pointer
	// * Replicas int32
	Replicas  int32
	Container ContainerSpec
}

type ContainerSpec struct {
	Image string
	port  int32
}

type DestroymentStatus struct {
	Phase string
}

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	///_ "k8s.io/code-generator"
)

// +genclient
// +groupName=heheh.com
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
	// Items is the list of Deployments.
	Items []Destroyment `json:"items,omitempty"`

	//Destroymentss []Destroyment `json:"destroymentss,omitempty"`
}

type DestroymentSpec struct {
	//may  be need to define as a pointer
	// * Replicas int32
	Replicas  *int32        `json:"replicas,omitempty"`
	Container ContainerSpec `json:"container,omitempty"`
}

type ContainerSpec struct {
	Image string `json:"image,omitempty"`
	Port  int32  `json:"port,omitempty"`
}

type DestroymentStatus struct {
	Phase string `json:"phase,omitempty"`
}

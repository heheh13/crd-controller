package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	///_ "k8s.io/code-generator"
)

// +genclient
// +groupName=heheh.com
//+kubebuilder:object:root=true
//+kubebuilder:object:generator=true
// +kubebuilder:resource:path=destroyments,singular=destroyment,shortName=ds,categories={}
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:JSONPath=".metadata.creationTimestamp",name=Age,type=date
// +kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=integer
// +kubebuilder:printcolumn:JSONPath=".status.phase",name=Status,type=string
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Destroyment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DestroymentSpec `json:"spec,omitempty"`
	// +optional
	Status DestroymentStatus `json:"status"`
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
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Replicas    *int32        `json:"replicas,omitempty"`
	Container   ContainerSpec `json:"container,omitempty"`
	ServiceSpec ServiceSpec   `json:"serviceSpec,omitempty"`
}

type ContainerSpec struct {
	// +kubebuilder:validation:MaxLength=50
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image,omitempty"`
	Port  int32  `json:"port,omitempty"`
}

type ServiceSpec struct {
	//+kubebuilder:default=ClusterIP
	ServiceType string `json:"serviceType,omitempty"`
}

type DestroymentStatus struct {
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1
	Phase             string `json:"phase,omitempty"`
	AvailableReplicas int32  `json:"availableReplicas,omitempty"`
	Replicas          int32  `json:"replicas,omitempty"`
}

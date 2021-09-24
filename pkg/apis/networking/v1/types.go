package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// FederationDNS records
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:categories=dns,singular=federationdns,path=federationdns,shortName=fddns,scope=Namespaced
// +kubebuilder:subresource:status
type FederationDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FederationDNSSpec   `json:"spec,omitempty"`
	Status            FederationDNSStatus `json:"status,omitempty"`
}

type FederationDNSStatus struct {
	Process string
}

type FederationDNSSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=128
	Host string `json:"host"`
	// +listType=map
	// +optional
	RR []string `json:"rr"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FederationDNSList is a list of Hero resources.
type FederationDNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FederationDNS `json:"items"`
}

// FederationDNSSlice records
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:singular=federationdnsslice,path=federationdnsslice,shortName=fddnsslice;sh,scope=Namespaced,categories=dns;federationdns
// +kubebuilder:subresource:status
type FederationDNSSlice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FederationDNSSliceSpec `json:"spec,omitempty"`
	Status            FederationDNSStatus    `json:"status,omitempty"`
}

type FederationDNSSliceSpec struct {
	Items []FederationDNSSpec `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FederationDNSSliceList is a list of Hero resources.
type FederationDNSSliceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []FederationDNSSlice `json:"items"`
}

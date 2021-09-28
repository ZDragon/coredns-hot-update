package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// HostEntry records
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:categories=dns,singular=hostentry,path=hostentries,shortName=fdhe,scope=Namespaced
// +kubebuilder:subresource:status
type HostEntry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HostEntrySpec   `json:"spec,omitempty"`
	Status            HostEntryStatus `json:"status,omitempty"`
}

type HostEntryStatus struct {
	Process    string      `json:"process"`
	LastUpdate metav1.Time `json:"lastUpdate"`
}

type HostEntrySpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=128
	Host string `json:"host"`
	// +listType=map
	// +optional
	RR []string `json:"rr"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HostEntryList is a list of Hero resources.
type HostEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []HostEntry `json:"items"`
}

// HostEntriesSlice records
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:singular=hostentriesslice,path=hostentriesslice,shortName=fdhes,scope=Namespaced,categories=dns
// +kubebuilder:subresource:status
type HostEntriesSlice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HostEntriesSliceSpec `json:"spec,omitempty"`
	Status            HostEntryStatus      `json:"status,omitempty"`
}

type HostEntriesSliceSpec struct {
	Items []HostEntrySpec `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HostEntriesSliceList is a list of Hero resources.
type HostEntriesSliceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []HostEntriesSlice `json:"items"`
}

package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// FederationDNS records
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:singular=hero,path=heroes,shortName=he;sh,scope=Namespaced,categories=heroes;superheroes
type FederationDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FederationDNSSpec `json:"spec,omitempty"`
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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EntitlementSpec defines the desired state of Entitlement
type EntitlementSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Service    string  `json:"service,omitempty"` // Appstudio
	QuotaItems []Quota `json:"quotaItems"`
}

type Quota struct {
	Resource string `json:"resource,omitempty"` // apibindings
	Limits   string `json:"limits,omitempty"`   // count
}

// EntitlementStatus defines the observed state of Entitlement
type EntitlementStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +crd
// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories=kcp
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Entitlement is the Schema for the entitlements API
type Entitlement struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EntitlementSpec   `json:"spec,omitempty"`
	Status EntitlementStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EntitlementList contains a list of Entitlement
type EntitlementList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Entitlement `json:"items"`
}

func (e EntitlementList) DeepCopyObject() runtime.Object {
	//TODO implement me
	panic("implement me")
}

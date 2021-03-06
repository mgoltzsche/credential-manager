package v1alpha1

import (
	//cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/operator-framework/operator-sdk/pkg/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionSynced  = status.ConditionType("Synced")
	ConditionReady   = status.ConditionType("Ready")
	ReasonFailedSync = status.ConditionReason("FailedSync")
	ReasonUpdating   = status.ConditionReason("Updating")
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ImageRegistrySpec defines the desired state of ImageRegistry
type ImageRegistrySpec struct {
	Replicas              *int32                    `json:"replicas,omitempty"`
	PersistentVolumeClaim PersistentVolumeClaimSpec `json:"persistentVolumeClaim"`
	TLS                   CertificateSpec           `json:"tls,omitempty"`
	Auth                  AuthSpec                  `json:"auth,omitempty"`
}

// PersistentVolumeClaimSpec specifies the PersistentVolumeClaim that should be maintained
type PersistentVolumeClaimSpec struct {
	StorageClassName *string                             `json:"storageClassName,omitempty"`
	AccessModes      []corev1.PersistentVolumeAccessMode `json:"accessModes,omitempty" protobuf:"bytes,1,rep,name=accessModes,casttype=PersistentVolumeAccessMode"`
	Resources        corev1.ResourceRequirements         `json:"resources,omitempty" protobuf:"bytes,2,opt,name=resources"`
	DeleteClaim      bool                                `json:"deleteClaim,omitempty"`
}

// AuthSpec specifies the CA certificate and optional docker_auth ConfigMap name
type AuthSpec struct {
	ConfigMapName *string         `json:"configMapName,omitempty"`
	CA            CertificateSpec `json:"ca"`
}

// CertificateSpec refers to a secret and an optional issuer to generate it
type CertificateSpec struct {
	IssuerRef  *CertIssuerRefSpec `json:"issuerRef,omitempty"`
	SecretName *string            `json:"secretName,omitempty"`
}

// CertificateIssuerSpec refers to a certificate issuer
type CertIssuerRefSpec struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// ImageRegistryStatus defines the observed state of ImageRegistry
type ImageRegistryStatus struct {
	ObservedGeneration int64             `json:"observedGeneration,omitempty"`
	Conditions         status.Conditions `json:"conditions,omitempty"`
	Hostname           string            `json:"hostname,omitempty"`
	TLSSecretName      string            `json:"tlsSecretName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageRegistry is the Schema for the imageregistries API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=imageregistries,scope=Namespaced
type ImageRegistry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageRegistrySpec   `json:"spec,omitempty"`
	Status ImageRegistryStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageRegistryList contains a list of ImageRegistry
type ImageRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageRegistry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImageRegistry{}, &ImageRegistryList{})
}

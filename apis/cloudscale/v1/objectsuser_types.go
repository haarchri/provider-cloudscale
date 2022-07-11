package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
)

// ObjectsUserSpec defines the desired state of an ObjectsUser.
type ObjectsUserSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`

	// SecretRef contains the name of the Secret where the credentials of the ObjectsUser are stored.
	// Must be a name that Kubernetes accepts as Secret name (lowercase RFC 1123 subdomain).
	SecretRef string `json:"secretRef"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced,categories={appcat,s3}

// ObjectsUser is the API for creating S3 Objects users on cloudscale.ch.
type ObjectsUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ObjectsUserSpec   `json:"spec"`
	Status ObjectsUserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ObjectsUserList contains a list of ObjectsUser
type ObjectsUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ObjectsUser `json:"items"`
}

// ObjectsUser type metadata.
var (
	ObjectsUserKind             = reflect.TypeOf(ObjectsUser{}).Name()
	ObjectsUserGroupKind        = schema.GroupKind{Group: Group, Kind: ObjectsUserKind}.String()
	ObjectsUserKindAPIVersion   = ObjectsUserKind + "." + SchemeGroupVersion.String()
	ObjectsUserGroupVersionKind = SchemeGroupVersion.WithKind(ObjectsUserKind)
)

func init() {
	SchemeBuilder.Register(&ObjectsUser{}, &ObjectsUserList{})
}

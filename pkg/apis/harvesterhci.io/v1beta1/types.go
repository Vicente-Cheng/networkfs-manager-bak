package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetworkFSState string

const (
	// NetworkFSStateEnabled indicates the networkFS endpoint is enabled
	NetworkFSStateEnabled NetworkFSState = "Enabled"
	// NetworkFSStateEnabling indicates the networkFS endpoint is enabling
	NetworkFSStateNotEnabling NetworkFSState = "Enabling"
	// NetworkFSStateDisabling indicates the networkFS endpoint is disabling
	NetworkFSStateDisabling NetworkFSState = "Disabling"
	// NetworkFSStateDisabled indicates the networkFS endpoint is disabled
	NetworkFSStateDisabled NetworkFSState = "Disabled"
	// NetworkFSStateUnknown indicates the networkFS endpoint state is unknown (initial state)
	NetworkFSStateUnknown NetworkFSState = "Unknown"

	// NetworkFSTypeNFS indicates the networkFS endpoint is NFS
	NetworkFSTypeNFS string = "NFS"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:shortName=netfilesystem;netfilesystems,scope=Namespaced
// +kubebuilder:printcolumn:name="DesiredState",type="string",JSONPath=`.spec.desiredState`
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=`.status.endpoint`
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=`.status.state`
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=`.status.type`

type NetworkFilesystem struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              NetworkFSSpec   `json:"spec"`
	Status            NetworkFSStatus `json:"status"`
}

type NetworkFSSpec struct {
	// name of the networkFS to which the endpoint is exported
	// +kubebuilder:validation:Required
	NetworkFSName string `json:"networkFSName"`

	// desired state of the networkFS endpoint, options are "Disabled", "Enabling", "Enabled", "Disabling", or "Unknown"
	// +kubebuilder:validation:Required:Enum:=Disabled;Enabling;Enabled;Disabling;Unknown
	DesiredState NetworkFSState `json:"desiredState"`
}

type NetworkFSStatus struct {
	// the current Endpoint of the networkFS
	// +kubebuilder:validation:
	// +kubebuilder:default:=""
	Endpoint string `json:"endpoint"`

	// the current state of the networkFS endpoint, options are "Enabled", "Enabling", "Disabling", "Disabled", or "Unknown"
	// +kubebuilder:validation:Enum:=Enabled;Enabling;Disabling;Disabled;Unknown
	// +kubebuilder:default:=Disabled
	State NetworkFSState `json:"state"`

	// the type of the networkFS endpoint, options are "NFS", or "Unknown"
	// +kubebuilder:validation:Enum:=NFS;Unknown
	// +kubebuilder:default:=NFS
	Type string `json:"type"`
}

/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SandboxEnvironmentSpec defines the desired state of SandboxEnvironment
type SandboxEnvironmentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SandboxEnvironment. Edit sandboxenvironment_types.go to remove/update
	Runtime RuntimeSpec `json:"runtime"`
	Resources ResourcesSpec `json:"resources"`
	Storage StorageSpec `json:"storage,omitempty"`
	Network NetworkSpec `json:"network,omitempty"`
	Security SecuritySpec `json:"security,omitempty"`
	Timeout metav1.Duration `json:"timeout,omitempty"`
}

type RuntimeSpec struct {
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// +kubebuilder:validation:Enum=python;cpp;java
	Language string `json:"language"`
	Command []string `json:"command,omitempty"`
}

type ResourcesSpec struct {
	Requests RequestSpec `json:"requests,omitempty"`
	Limits LimitSpec `json:"limits,omitempty"`
}

type RequestSpec struct {
	// +kubebuilder:validation:Pattern=`^[0-9]+m?$`
	CPU string `json:"cpu"`

	// +kubebuilder:validation:Pattern=`^[0-9]+(Mi|Gi)$`
	Memory string `json:"memory"`
}

type LimitSpec struct {
	// +kubebuilder:validation:Pattern=`^[0-9]+m?$`
	CPU string `json:"cpu,omitempty"`

	// +kubebuilder:validation:Pattern=`^[0-9]+(Mi|Gi)$`
	Memory string `json:"memory,omitempty"`
}

type StorageSpec struct {
	// +kubebuilder:validation:Pattern=`^[0-9]+Gi$`
	Size string `json:"size"`

	// +kubebuilder:validation:MinLength=1
	MountPath string `json:"mountPath,omitempty"`
}

type NetworkSpec struct {
	Enabled bool `json:"enabled,omitempty"`
}

type SecuritySpec struct {
	RunAsNonRoot bool `json:"runAsNonRoot,omitempty"`
	ReadOnlyRootFilesystem bool `json:"readOnlyRootFilesystem,omitempty"`
	AllowPrivilegeEscalation bool `json:"allowPrivilegeEscalation,omitempty"`
}

// SandboxEnvironmentStatus defines the observed state of SandboxEnvironment
type SandboxEnvironmentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase    string `json:"phase,omitempty"`
	Ready    bool   `json:"ready,omitempty"`
	PodName  string `json:"podName,omitempty"`
	PodIP    string `json:"podIP,omitempty"`
	Message  string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SandboxEnvironment is the Schema for the sandboxenvironments API
type SandboxEnvironment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SandboxEnvironmentSpec   `json:"spec,omitempty"`
	Status SandboxEnvironmentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SandboxEnvironmentList contains a list of SandboxEnvironment
type SandboxEnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SandboxEnvironment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SandboxEnvironment{}, &SandboxEnvironmentList{})
}

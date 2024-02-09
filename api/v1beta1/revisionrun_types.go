/*
Copyright 2024 patrick.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RevisionRunSpec defines the desired state of RevisionRun
type RevisionRunSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of RevisionRun. Edit revisionrun_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// RevisionRunStatus defines the observed state of RevisionRun
type RevisionRunStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RevisionRun is the Schema for the revisionruns API
type RevisionRun struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RevisionRunSpec   `json:"spec,omitempty"`
	Status RevisionRunStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RevisionRunList contains a list of RevisionRun
type RevisionRunList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RevisionRun `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RevisionRun{}, &RevisionRunList{})
}

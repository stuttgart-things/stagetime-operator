/*
Copyright 2024 PATRICK HERMANN patrick.hermann@sva.de

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

// SchedulingConfig defines scheduling related properties.
type Technologies struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Path string `json:"path,omitempty"`
	// +kubebuilder:default=99
	Stage int `json:"stage,omitempty"`
	// +kubebuilder:default=false
	Canfail    bool   `json:"canfail,omitempty"`
	Resolver   string `json:"resolver,omitempty"`
	Params     string `json:"params,omitempty"`
	Listparams string `json:"listparams,omitempty"`
	Vclaims    string `json:"vclaims,omitempty"`
}

// RevisionRunSpec defines the desired state of RevisionRun
type RevisionRunSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Revision         string          `json:"revision,omitempty"`
	Repository       string          `json:"repository"`
	TechnologyConfig []*Technologies `json:"technologies"`
}

// RevisionRunStatus defines the observed state of RevisionRun
type RevisionRunStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`

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

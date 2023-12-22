/*
Copyright 2023.

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

// if move.Status.State == "duplicate" || move.Status.State == "notAllowed" || move.Status.State == "processed" {
const (
	NotAllowed = "NotAllowed"
	Duplicate  = "Duplicate"
	Processing = "Processing"
	Processed  = "Processed"
)

// MoveSpec defines the desired state of Move
type MoveSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	TicTacToeName string `json:"ticTacToeName"`

	Row    int `json:"row,omitempty"`
	Column int `json:"column,omitempty"`
}

// MoveStatus defines the observed state of Move
type MoveStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:default=Processing
	State string `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Move is the Schema for the moves API
type Move struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MoveSpec   `json:"spec,omitempty"`
	Status MoveStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MoveList contains a list of Move
type MoveList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Move `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Move{}, &MoveList{})
}

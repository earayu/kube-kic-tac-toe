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

const (
	NoPlayer = iota
	Human
	Bot
)

// TicTacToeSpec defines the desired state of TicTacToe
type TicTacToeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

}

// TicTacToeStatus defines the observed state of TicTacToe
type TicTacToeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	MoveHistory MoveList `json:"move,omitempty"`

	Chessboard string `json:"chessboard,omitempty"`

	// State: playing,draw,humanWins,botWins
	State string `json:"state,omitempty"`

	// default current time now
	Version metav1.Time `json:"timestamp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ttt

// TicTacToe is the Schema for the tictactoes API
type TicTacToe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TicTacToeSpec   `json:"spec,omitempty"`
	Status TicTacToeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TicTacToeList contains a list of TicTacToe
type TicTacToeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TicTacToe `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TicTacToe{}, &TicTacToeList{})
}

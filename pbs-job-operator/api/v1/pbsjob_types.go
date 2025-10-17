/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PBSJobSpec defines the desired state of PBSJob
type PBSJobSpec struct {
	// Script is the PBS job script content to be executed
	// +required
	Script string `json:"script"`

	// JobName is the name of the PBS job (used with qsub -N)
	// +optional
	JobName string `json:"jobName,omitempty"`

	// Queue specifies which PBS queue to submit the job to
	// +optional
	Queue string `json:"queue,omitempty"`

	// Resources specifies PBS resource requirements (nodes, ppn, walltime, etc.)
	// +optional
	Resources map[string]string `json:"resources,omitempty"`

	// Walltime specifies the maximum execution time for the job
	// +optional
	Walltime string `json:"walltime,omitempty"`

	// WorkingDirectory specifies where the job should execute
	// +optional
	WorkingDirectory string `json:"workingDirectory,omitempty"`

	// OutputPath specifies where to write job output
	// +optional
	OutputPath string `json:"outputPath,omitempty"`

	// ErrorPath specifies where to write job error output
	// +optional
	ErrorPath string `json:"errorPath,omitempty"`

	// Environment specifies environment variables for the job
	// +optional
	Environment map[string]string `json:"environment,omitempty"`
}

// PBSJobStatus defines the observed state of PBSJob.
type PBSJobStatus struct {
	// PBSJobID is the ID assigned by PBS when the job is submitted
	// +optional
	PBSJobID string `json:"pbsJobID,omitempty"`

	// Phase represents the current phase of the PBS job
	// Possible values: Pending, Queued, Running, Completed, Failed, Deleted
	// +optional
	Phase string `json:"phase,omitempty"`

	// State represents the detailed PBS job state from qstat
	// +optional
	State string `json:"state,omitempty"`

	// SubmissionTime is when the job was submitted to PBS
	// +optional
	SubmissionTime *metav1.Time `json:"submissionTime,omitempty"`

	// StartTime is when the job started running
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is when the job completed
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// ExitCode is the exit code of the completed job
	// +optional
	ExitCode *int32 `json:"exitCode,omitempty"`

	// Message provides additional information about the current state
	// +optional
	Message string `json:"message,omitempty"`

	// conditions represent the current state of the PBSJob resource.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PBSJob is the Schema for the pbsjobs API
type PBSJob struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of PBSJob
	// +required
	Spec PBSJobSpec `json:"spec"`

	// status defines the observed state of PBSJob
	// +optional
	Status PBSJobStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// PBSJobList contains a list of PBSJob
type PBSJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PBSJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PBSJob{}, &PBSJobList{})
}

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

package controller

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "github.com/subhas/pbs-job-operator/api/v1"
	"github.com/subhas/pbs-job-operator/internal/pbs"
)

const (
	PBSJobFinalizer = "batch.example.com/pbsjob-finalizer"
)

// PBSJobReconciler reconciles a PBSJob object
type PBSJobReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	PBSClient *pbs.PBSClient
}

// +kubebuilder:rbac:groups=batch.example.com,resources=pbsjobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch.example.com,resources=pbsjobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch.example.com,resources=pbsjobs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PBSJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// Fetch the PBSJob instance
	var pbsJob batchv1.PBSJob
	if err := r.Get(ctx, req.NamespacedName, &pbsJob); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get PBSJob")
		return ctrl.Result{}, err
	}

	// Check if the object is being deleted
	if !pbsJob.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, &pbsJob)
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(&pbsJob, PBSJobFinalizer) {
		controllerutil.AddFinalizer(&pbsJob, PBSJobFinalizer)
		if err := r.Update(ctx, &pbsJob); err != nil {
			logger.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Handle the job based on its current status
	switch pbsJob.Status.Phase {
	case "":
		// New job - submit to PBS
		return r.submitJob(ctx, &pbsJob)
	case "Pending", "Queued", "Running":
		// Job is active - check status
		return r.updateJobStatus(ctx, &pbsJob)
	case "Completed", "Failed":
		// Job is finished - no action needed
		return ctrl.Result{}, nil
	default:
		// Unknown phase - check status
		return r.updateJobStatus(ctx, &pbsJob)
	}
}

func (r *PBSJobReconciler) submitJob(ctx context.Context, pbsJob *batchv1.PBSJob) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// Submit job to PBS
	jobID, err := r.PBSClient.SubmitJob(pbsJob)
	if err != nil {
		// Update status to reflect failure
		pbsJob.Status.Phase = "Failed"
		pbsJob.Status.Message = fmt.Sprintf("Failed to submit job: %v", err)
		r.updateCondition(pbsJob, "Submitted", metav1.ConditionFalse, "SubmissionFailed", err.Error())

		if updateErr := r.Status().Update(ctx, pbsJob); updateErr != nil {
			logger.Error(updateErr, "Failed to update status after submission failure")
		}
		return ctrl.Result{}, err
	}

	// Update status with job ID and phase
	pbsJob.Status.PBSJobID = jobID
	pbsJob.Status.Phase = "Queued"
	pbsJob.Status.Message = "Job submitted to PBS"
	now := metav1.Now()
	pbsJob.Status.SubmissionTime = &now
	r.updateCondition(pbsJob, "Submitted", metav1.ConditionTrue, "SubmissionSuccessful", "Job submitted successfully")

	if err := r.Status().Update(ctx, pbsJob); err != nil {
		logger.Error(err, "Failed to update status after submission")
		return ctrl.Result{}, err
	}

	logger.Info("Job submitted successfully", "jobID", jobID)

	// Requeue to check status
	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (r *PBSJobReconciler) updateJobStatus(ctx context.Context, pbsJob *batchv1.PBSJob) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	if pbsJob.Status.PBSJobID == "" {
		logger.Info("No PBS job ID found, skipping status update")
		return ctrl.Result{}, nil
	}

	// Query PBS for job status
	jobInfo, err := r.PBSClient.GetJobInfo(pbsJob.Status.PBSJobID)
	if err != nil {
		logger.Error(err, "Failed to get job info from PBS", "jobID", pbsJob.Status.PBSJobID)

		// Job might have completed and been removed from queue
		if pbsJob.Status.Phase == "Running" {
			pbsJob.Status.Phase = "Completed"
			pbsJob.Status.Message = "Job completed (no longer in PBS queue)"
			now := metav1.Now()
			pbsJob.Status.CompletionTime = &now
			r.updateCondition(pbsJob, "Completed", metav1.ConditionTrue, "JobCompleted", "Job completed successfully")

			if updateErr := r.Status().Update(ctx, pbsJob); updateErr != nil {
				logger.Error(updateErr, "Failed to update status")
			}
		}

		return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
	}

	// Update status based on PBS job info
	oldPhase := pbsJob.Status.Phase
	newPhase := pbs.GetJobState(jobInfo.State)

	pbsJob.Status.State = jobInfo.State
	pbsJob.Status.Phase = newPhase

	// Update timestamps
	if !jobInfo.StartTime.IsZero() && pbsJob.Status.StartTime == nil {
		startTime := metav1.NewTime(jobInfo.StartTime)
		pbsJob.Status.StartTime = &startTime
	}

	if !jobInfo.CompleteTime.IsZero() && pbsJob.Status.CompletionTime == nil {
		completeTime := metav1.NewTime(jobInfo.CompleteTime)
		pbsJob.Status.CompletionTime = &completeTime
	}

	if jobInfo.ExitCode != 0 {
		pbsJob.Status.ExitCode = &jobInfo.ExitCode
	}

	// Update condition based on phase change
	if oldPhase != newPhase {
		switch newPhase {
		case "Running":
			r.updateCondition(pbsJob, "Running", metav1.ConditionTrue, "JobStarted", "Job started running")
		case "Completed":
			r.updateCondition(pbsJob, "Completed", metav1.ConditionTrue, "JobCompleted", "Job completed successfully")
		case "Failed":
			r.updateCondition(pbsJob, "Failed", metav1.ConditionTrue, "JobFailed", "Job failed")
		}
	}

	if err := r.Status().Update(ctx, pbsJob); err != nil {
		logger.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	// Requeue based on job phase
	switch newPhase {
	case "Queued", "Running":
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	default:
		return ctrl.Result{}, nil
	}
}

func (r *PBSJobReconciler) handleDeletion(ctx context.Context, pbsJob *batchv1.PBSJob) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// If PBS job ID exists, delete the job from PBS
	if pbsJob.Status.PBSJobID != "" {
		if err := r.PBSClient.DeleteJob(pbsJob.Status.PBSJobID); err != nil {
			logger.Error(err, "Failed to delete PBS job", "jobID", pbsJob.Status.PBSJobID)
			// Continue with finalizer removal even if PBS deletion fails
		}
	}

	// Remove finalizer
	controllerutil.RemoveFinalizer(pbsJob, PBSJobFinalizer)
	if err := r.Update(ctx, pbsJob); err != nil {
		logger.Error(err, "Failed to remove finalizer")
		return ctrl.Result{}, err
	}

	logger.Info("PBSJob deleted successfully")
	return ctrl.Result{}, nil
}

func (r *PBSJobReconciler) updateCondition(pbsJob *batchv1.PBSJob, conditionType string, status metav1.ConditionStatus, reason, message string) {
	condition := metav1.Condition{
		Type:               conditionType,
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             reason,
		Message:            message,
	}

	// Find existing condition or append new one
	for i, existingCondition := range pbsJob.Status.Conditions {
		if existingCondition.Type == conditionType {
			pbsJob.Status.Conditions[i] = condition
			return
		}
	}
	pbsJob.Status.Conditions = append(pbsJob.Status.Conditions, condition)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PBSJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.PBSJob{}).
		Named("pbsjob").
		Complete(r)
}

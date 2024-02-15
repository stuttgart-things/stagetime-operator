package controllers

import (
	"context"

	stagetimev1beta1 "github.com/stuttgart-things/stagetime-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ConditionStatus defines RevisionRun condition status.
type ConditionStatus string

const (
	TypeAvailable   ConditionStatus = "AVAILABLE"
	TypeProgressing ConditionStatus = "PROGRESSING"
	TypeDegraded    ConditionStatus = "DEGRADED"
)

func (r *RevisionRunReconciler) SetInitialCondition(ctx context.Context, req ctrl.Request, RevisionRun *stagetimev1beta1.RevisionRun) error {
	if RevisionRun.Status.Conditions != nil || len(RevisionRun.Status.Conditions) != 0 {
		return nil
	}

	err := r.SetCondition(ctx, req, RevisionRun, TypeAvailable, "STARTING RECONCILIATION")

	return err
}

// SetCondition sets the status condition of the RevisionRun.
func (r *RevisionRunReconciler) SetCondition(
	ctx context.Context, req ctrl.Request,
	revisionRun *stagetimev1beta1.RevisionRun, condition ConditionStatus,
	message string,
) error {
	log := log.FromContext(ctx)

	meta.SetStatusCondition(
		&revisionRun.Status.Conditions,
		metav1.Condition{
			Type:   string(condition),
			Status: metav1.ConditionUnknown, Reason: "RECONCILING",
			Message: message,
		},
	)

	if err := r.Status().Update(ctx, revisionRun); err != nil {
		log.Error(err, "FAILED TO UPDATE REVISIONRUN STATUS")

		return err
	}

	if err := r.Get(ctx, req.NamespacedName, revisionRun); err != nil {
		log.Error(err, "FAILED TO RE-FETCH REVISIONRUN")

		return err
	}

	return nil
}

func (r *RevisionRunReconciler) GetRevisionRun(ctx context.Context, req ctrl.Request, revisionRun *stagetimev1beta1.RevisionRun) error {
	err := r.Get(ctx, req.NamespacedName, revisionRun)
	if err != nil {
		return err
	}

	return nil
}

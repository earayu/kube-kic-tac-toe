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

package controller

import (
	"context"
	"fmt"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	earayugithubiov1alpha1 "earayu.github.io/kube-kic-tac-toe/api/v1alpha1"
)

// MoveReconciler reconciles a Move object
type MoveReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=earayu.github.io.earayu.github.io,resources=moves,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=earayu.github.io.earayu.github.io,resources=moves/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=earayu.github.io.earayu.github.io,resources=moves/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Move object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *MoveReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	move := earayugithubiov1alpha1.Move{}
	if err := r.Get(ctx, req.NamespacedName, &move); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	namespacedName := client.ObjectKey{
		Name:      move.Spec.TicTacToeName,
		Namespace: req.Namespace,
	}
	ticTacToe := earayugithubiov1alpha1.TicTacToe{}
	if err := r.Get(ctx, namespacedName, &ticTacToe); err != nil {
		//todo game not found, should clean up all the moves
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if err := controllerutil.SetControllerReference(&ticTacToe, &move, r.Scheme); err != nil {
		l.Error(err, "unable to set controller reference")
		return ctrl.Result{}, err
	}
	if err := r.Update(ctx, &move); err != nil {
		l.Error(err, "unable to update Move status")
		return ctrl.Result{}, err
	}

	// ignore resources that we've already processed
	if move.Status.State == earayugithubiov1alpha1.Duplicate || move.Status.State == earayugithubiov1alpha1.NotAllowed {
		return reconcile.Result{}, nil
	}

	// validate row & column
	if move.Spec.Row < 0 || move.Spec.Row >= 3 || move.Spec.Column < 0 || move.Spec.Column >= 3 {
		move.Status.State = earayugithubiov1alpha1.NotAllowed
		if err := r.Status().Update(ctx, &move); err != nil {
			l.Error(err, "unable to update Move status")
			return ctrl.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	duplicate := earayugithubiov1alpha1.AppendMoveRef(&ticTacToe.Status, &move)
	if duplicate {
		l.Info(fmt.Sprintf("current position has been taken. row:%d, col:%d", move.Spec.Row, move.Spec.Column))
		move.Status.State = earayugithubiov1alpha1.Duplicate
		if err := r.Status().Update(ctx, &move); err != nil {
			l.Error(err, "unable to update TicTacToe status")
			return ctrl.Result{}, err
		}
		return reconcile.Result{}, nil
	} else {
		move.Status.State = earayugithubiov1alpha1.Processed
		if err := r.Status().Update(ctx, &ticTacToe); err != nil {
			l.Error(err, "unable to update TicTacToe status")
			return ctrl.Result{}, err
		}
		if err := r.Status().Update(ctx, &move); err != nil {
			l.Error(err, "unable to update TicTacToe status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *MoveReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&earayugithubiov1alpha1.Move{}).
		Complete(r)
}

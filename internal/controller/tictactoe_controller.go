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
	earayugithubiov1alpha1 "earayu.github.io/kube-kic-tac-toe/api/v1alpha1"
	"earayu.github.io/kube-kic-tac-toe/internal/controller/portable"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// TicTacToeReconciler reconciles a TicTacToe object
type TicTacToeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=earayu.github.io.earayu.github.io,resources=tictactoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=earayu.github.io.earayu.github.io,resources=tictactoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=earayu.github.io.earayu.github.io,resources=tictactoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the TicTacToe object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *TicTacToeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var ticTacToe earayugithubiov1alpha1.TicTacToe
	if err := r.Get(ctx, req.NamespacedName, &ticTacToe); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	board, err := portable.GetBoard(&ticTacToe)
	if err != nil {
		l.Info(fmt.Sprintf("parse row data failed, err:%s", err))
		return ctrl.Result{}, nil
	}
	// bot try to make a move
	nextPlayer := portable.NextPlayer(board)
	if nextPlayer == earayugithubiov1alpha1.Bot {
		newBoard, hasMoved := portable.RandomMove(board, nextPlayer)
		if hasMoved {
			ticTacToe.Status.Row1, ticTacToe.Status.Row2, ticTacToe.Status.Row3 = portable.GetRow(newBoard)
		}
	}
	// check winner
	winner, finished := portable.CheckWinner(board)
	if winner == earayugithubiov1alpha1.Human {
		ticTacToe.Status.State = "human wins"
	} else if winner == earayugithubiov1alpha1.Bot {
		ticTacToe.Status.State = "bot wins"
	} else if winner == earayugithubiov1alpha1.NoPlayer && finished {
		ticTacToe.Status.State = "draw"
	} else {
		ticTacToe.Status.State = "playing"
	}

	if err = r.Status().Update(ctx, &ticTacToe); err != nil {
		return ctrl.Result{}, fmt.Errorf("update status err:%w", err)
	}
	return ctrl.Result{}, nil
}

func (r *TicTacToeReconciler) AnnounceWinner(player int) {
	//todo
}

func flipPlayer(player int) int {
	if player == 1 {
		return 2
	} else {
		return 1
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *TicTacToeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&earayugithubiov1alpha1.TicTacToe{}).
		Complete(r)
}

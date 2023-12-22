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
	"earayu.github.io/kube-kic-tac-toe/internal/controller/portable"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"

	earayugithubiov1alpha1 "earayu.github.io/kube-kic-tac-toe/api/v1alpha1"
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
	//l := log.FromContext(ctx)

	var ticTacToe earayugithubiov1alpha1.TicTacToe
	if err := r.Get(ctx, req.NamespacedName, &ticTacToe); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	board, err := portable.GetBoard(&ticTacToe)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("parse row data failed, err:%w", err)
	}
	// check if player1 has won
	if portable.CheckWinner(board, 1) {
		r.AnnounceWinner(1)
		return ctrl.Result{}, nil
	}
	// check if player2 has won
	if portable.CheckWinner(board, 2) {
		r.AnnounceWinner(2)
		return ctrl.Result{}, nil
	}
	//todo process player's move

	// no winner
	newBoard, hasMoved := portable.RandomMove(board, ticTacToe.Status.CurrentPlayer)
	if hasMoved {
		ticTacToe.Status.CurrentPlayer = flipPlayer(ticTacToe.Status.CurrentPlayer)
	}
	ticTacToe.Status.Row1, ticTacToe.Status.Row2, ticTacToe.Status.Row3 = getRow(newBoard)
	err = r.Status().Update(ctx, &ticTacToe)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("update status err:%w", err)
	}
	return ctrl.Result{}, nil
}

func getRow(board portable.Board) (row1 string, row2 string, row3 string) {
	sep := ""
	for i := range board[0] {
		row1 += sep + strconv.Itoa(board[0][i])
		sep = " "
	}

	sep = ""
	for i := range board[1] {
		row2 += sep + strconv.Itoa(board[1][i])
		sep = " "
	}

	sep = ""
	for i := range board[2] {
		row3 += sep + strconv.Itoa(board[2][i])
		sep = " "
	}
	return row1, row2, row3
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

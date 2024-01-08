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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sort"
	"sync"
)

var (
	ticTacToeOwnerKey = ".spec.ticTacToeName"
)

// TicTacToeReconciler reconciles a TicTacToe object
type TicTacToeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	mu     sync.Mutex
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
	r.mu.Lock()
	defer r.mu.Unlock()
	l := log.FromContext(ctx)

	var ticTacToe earayugithubiov1alpha1.TicTacToe
	if err := r.Get(ctx, req.NamespacedName, &ticTacToe); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// List all moves with '.spec.ticTacToeName' field matching the TicTacToe's name.
	var allMoveList earayugithubiov1alpha1.MoveList
	if err := r.List(ctx, &allMoveList, client.InNamespace(req.Namespace), client.MatchingFields{ticTacToeOwnerKey: req.Name}); err != nil {
		return reconcile.Result{}, fmt.Errorf("list moves err:%w", err)
	}
	var processingMoves []earayugithubiov1alpha1.Move
	for _, move := range allMoveList.Items {
		if move.Status.State == earayugithubiov1alpha1.Processing || move.Status.State == "" {
			processingMoves = append(processingMoves, move)
		}
	}
	if len(processingMoves) == 0 {
		return ctrl.Result{}, nil
	}
	// order allMoveList by creationTime
	sort.Slice(processingMoves, func(i, j int) bool {
		return processingMoves[i].CreationTimestamp.Before(&processingMoves[j].CreationTimestamp)
	})
	for _, move := range processingMoves {
		invalid, duplicate, err := portable.PutOnBoard(&ticTacToe.Status.MoveHistory, move)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("put on board err:%w", err)
		}
		if invalid {
			move.Status.State = earayugithubiov1alpha1.NotAllowed
		} else if duplicate {
			move.Status.State = earayugithubiov1alpha1.Duplicate
		} else {
			move.Status.State = earayugithubiov1alpha1.Processed
		}
		if err := r.Status().Update(ctx, &move); err != nil {
			l.Error(err, "unable to update Move status")
			return ctrl.Result{}, err
		}
	}

	board := portable.GetBoard(&ticTacToe.Status.MoveHistory)
	ticTacToe.Status.Chessboard1, ticTacToe.Status.Chessboard2, ticTacToe.Status.Chessboard3 = portable.GetChessBoard(board)

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
	if err := r.Status().Update(ctx, &ticTacToe); err != nil {
		return ctrl.Result{}, fmt.Errorf("update status err:%w", err)
	}

	if ticTacToe.Status.State == "playing" {
		// bot try to make a move
		nextPlayer := portable.NextPlayer(&ticTacToe.Status.MoveHistory)
		if nextPlayer == earayugithubiov1alpha1.Bot {
			row, col, hasMoved := portable.RandomMove(board)
			if hasMoved {
				botMove := &earayugithubiov1alpha1.Move{
					ObjectMeta: metav1.ObjectMeta{
						Name:      fmt.Sprintf("%s-bot-move-%d-%d", ticTacToe.Name, row, col),
						Namespace: ticTacToe.Namespace,
					},
					Spec: earayugithubiov1alpha1.MoveSpec{
						TicTacToeName: ticTacToe.Name,
						Player:        earayugithubiov1alpha1.Bot,
						Row:           row,
						Column:        col,
					},
				}
				controllerutil.SetControllerReference(&ticTacToe, botMove, r.Scheme)
				err := r.Create(ctx, botMove)
				if err != nil {
					return ctrl.Result{}, fmt.Errorf("create move err:%w", err)
				}
				return ctrl.Result{Requeue: true}, nil
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TicTacToeReconciler) SetupWithManager(mgr ctrl.Manager) error {

	// Index the Owner field so that we can efficiently look up all
	// TicTacToe objects that own a given Move object.
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &earayugithubiov1alpha1.Move{}, ticTacToeOwnerKey, func(rawObj client.Object) []string {
		// grab the Move object, extract the owner...
		move := rawObj.(*earayugithubiov1alpha1.Move)
		owner := metav1.GetControllerOf(move)
		if owner == nil {
			return nil
		}
		// ...make sure it's a TicTacToe...
		if owner.APIVersion != earayugithubiov1alpha1.GroupVersion.String() || owner.Kind != "TicTacToe" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&earayugithubiov1alpha1.TicTacToe{}).
		Owns(&earayugithubiov1alpha1.Move{}).
		Complete(r)
}

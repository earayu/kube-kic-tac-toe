package portable

import (
	earayugithubiov1alpha1 "earayu.github.io/kube-kic-tac-toe/api/v1alpha1"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Board = [][]int

func CheckWinner(board Board, player int) bool {
	for i := 0; i < 3; i++ {
		if board[i][0] == player && board[i][1] == player && board[i][2] == player {
			return true
		}
	}

	for i := 0; i < 3; i++ {
		if board[0][i] == player && board[1][i] == player && board[2][i] == player {
			return true
		}
	}

	if board[0][0] == player && board[1][1] == player && board[2][2] == player {
		return true
	}
	if board[0][2] == player && board[1][1] == player && board[2][0] == player {
		return true
	}

	return false
}

func RandomMove(board Board, player int) (newBoard Board, hasMoved bool) {
	var emptyPositions []int
	for i, row := range board {
		for j, val := range row {
			if val == 0 {
				emptyPositions = append(emptyPositions, i*3+j)
			}
		}
	}

	if len(emptyPositions) == 0 {
		return board, false
	}

	pos := emptyPositions[rand.Intn(len(emptyPositions))]
	row := pos / 3
	col := pos % 3

	board[row][col] = player
	return board, true
}

func Move(board Board, player int, row int, col int) (newBoard Board, boardIsFull bool, err error) {
	if board[row][col] != 0 {
		//todo split boardIsFull & duplicate
		return board, false, fmt.Errorf("not allowed, this position has beed taken")
	}

	board[row][col] = player
	return board, true, nil
}

func GetBoard(ticTacToe *earayugithubiov1alpha1.TicTacToe) (Board, error) {
	board := Board{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	for i, s := range strings.Split(ticTacToe.Status.Row1, " ") {
		if v, err := strconv.Atoi(s); err != nil {
			return nil, err
		} else {
			board[0][i] = v
		}
	}
	for i, s := range strings.Split(ticTacToe.Status.Row2, " ") {
		if v, err := strconv.Atoi(s); err != nil {
			return nil, err
		} else {
			board[0][i] = v
		}
	}
	for i, s := range strings.Split(ticTacToe.Status.Row3, " ") {
		if v, err := strconv.Atoi(s); err != nil {
			return nil, err
		} else {
			board[0][i] = v
		}
	}
	return board, nil
}

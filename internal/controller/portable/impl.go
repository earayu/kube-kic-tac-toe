package portable

import (
	earayugithubiov1alpha1 "earayu.github.io/kube-kic-tac-toe/api/v1alpha1"
	"math/rand"
	"strconv"
	"strings"
)

type Board = [][]int

func CheckWinner(board Board) (winner int, finished bool) {
	// Check rows and columns for a winner
	for i := 0; i < 3; i++ {
		if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][0] != earayugithubiov1alpha1.NoPlayer {
			return board[i][0], true
		}
		if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[0][i] != earayugithubiov1alpha1.NoPlayer {
			return board[0][i], true
		}
	}

	// Check diagonals for a winner
	if board[0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != earayugithubiov1alpha1.NoPlayer {
		return board[0][0], true
	}
	if board[0][2] == board[1][1] && board[1][1] == board[2][0] && board[0][2] != earayugithubiov1alpha1.NoPlayer {
		return board[0][2], true
	}

	// Check if the board is full (game finished with no winner)
	finished = true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == earayugithubiov1alpha1.NoPlayer {
				finished = false
				break
			}
		}
		if !finished {
			break
		}
	}

	return earayugithubiov1alpha1.NoPlayer, finished
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

func Move(board Board, player int, row int, col int) (newBoard Board, duplicate bool) {
	if board[row][col] != 0 {
		return board, true
	}

	board[row][col] = player
	return board, false
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

func GetRow(board Board) (row1 string, row2 string, row3 string) {
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

func NextPlayer(board Board) int {
	count1 := 0
	count2 := 0
	for _, row := range board {
		for _, v := range row {

			if v == earayugithubiov1alpha1.Human {
				count1++
			} else if v == earayugithubiov1alpha1.Bot {
				count2++
			}
		}
	}
	return 1 + count1 - count2
}

package portable

import (
	earayugithubiov1alpha1 "earayu.github.io/kube-kic-tac-toe/api/v1alpha1"
	"math/rand"
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

func RandomMove(board Board) (row int, col int, hasMoved bool) {
	var emptyPositions []int
	for i, row := range board {
		for j, val := range row {
			if val == 0 {
				emptyPositions = append(emptyPositions, i*3+j)
			}
		}
	}

	if len(emptyPositions) == 0 {
		return row, col, false
	}

	pos := emptyPositions[rand.Intn(len(emptyPositions))]
	row = pos / 3
	col = pos % 3

	return row, col, true
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

	for _, m := range ticTacToe.Status.MoveHistory {
		board[m.Spec.Row][m.Spec.Column] = m.Spec.Player
	}

	return board, nil
}

// GetChessBoard takes a Board and returns a string representation
func GetChessBoard(board Board) (chessBoard string) {
	var sb strings.Builder
	for i, row := range board {
		for j, cell := range row {
			switch cell {
			case 0:
				sb.WriteString(" - ") // Empty cell
			case 1:
				sb.WriteString(" O ") // Player 1
			case 2:
				sb.WriteString(" X ") // Player 2
			}
			if j < len(row)-1 {
				sb.WriteString("|") // Column separator
			}
		}
		if i < len(board)-1 {
			sb.WriteString("\n-----------\n") // Row separator
		}
	}
	return sb.String()
}

func NextPlayer(status *earayugithubiov1alpha1.TicTacToeStatus) int {
	moveCount := len(status.MoveHistory)
	if moveCount == 0 {
		return earayugithubiov1alpha1.Human
	}
	lastPlayer := status.MoveHistory[moveCount-1].Spec.Player
	if lastPlayer == earayugithubiov1alpha1.Human {
		return earayugithubiov1alpha1.Bot
	}
	return earayugithubiov1alpha1.Human
}

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

func GetBoard(moveHistory *earayugithubiov1alpha1.MoveList) Board {
	board := Board{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}

	for _, m := range moveHistory.Items {
		board[m.Spec.Row][m.Spec.Column] = m.Spec.Player
	}

	return board
}

func PutOnBoard(moveHistory *earayugithubiov1alpha1.MoveList, nextMove earayugithubiov1alpha1.Move) (invalid bool, duplicate bool, err error) {
	board := GetBoard(moveHistory)
	nextPlayer := NextPlayer(moveHistory)
	if nextMove.Spec.Player != nextPlayer {
		return true, false, nil
	}
	if nextMove.Spec.Row < 0 || nextMove.Spec.Row >= 3 || nextMove.Spec.Column < 0 || nextMove.Spec.Column >= 3 {
		return true, false, nil
	}
	_, duplicate = Move(board, nextMove.Spec.Player, nextMove.Spec.Row, nextMove.Spec.Column)
	if duplicate {
		return false, true, nil
	}

	moveHistory.Items = append(moveHistory.Items, nextMove)
	return false, false, nil
}

// GetChessBoard takes a Board and returns a string that represents the chessboard.
func GetChessBoard(board Board) (string, string, string) {
	chessBoard := make([]string, 3)
	for i, row := range board {
		var sb strings.Builder
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
		chessBoard[i] = sb.String() // Row separator
	}
	return chessBoard[0], chessBoard[1], chessBoard[2]
}

func NextPlayer(moveHistory *earayugithubiov1alpha1.MoveList) int {
	moveCount := len(moveHistory.Items)
	if moveCount == 0 {
		return earayugithubiov1alpha1.Human
	}
	lastPlayer := moveHistory.Items[moveCount-1].Spec.Player
	if lastPlayer == earayugithubiov1alpha1.Human {
		return earayugithubiov1alpha1.Bot
	}
	return earayugithubiov1alpha1.Human
}

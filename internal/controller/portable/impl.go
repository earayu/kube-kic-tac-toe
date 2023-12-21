package portable

import "math/rand"

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

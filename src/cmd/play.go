package main

import (
	"fmt"

	hexit "github.com/uyhcire/hexit/src"
)

func formatBoardSquare(boardSquareValue byte) string {
	if boardSquareValue == 1 {
		return "X"
	} else if boardSquareValue == 2 {
		return "O"
	} else {
		return "-"
	}
}

func PrintBoard(board *hexit.Board) {
	fmt.Printf(
		"%s %s %s %s %s\n"+
			" %s %s %s %s %s\n"+
			"  %s %s %s %s %s\n"+
			"   %s %s %s %s %s\n"+
			"    %s %s %s %s %s\n",
		formatBoardSquare(board[0][0]),
		formatBoardSquare(board[0][1]),
		formatBoardSquare(board[0][2]),
		formatBoardSquare(board[0][3]),
		formatBoardSquare(board[0][4]),
		formatBoardSquare(board[1][0]),
		formatBoardSquare(board[1][1]),
		formatBoardSquare(board[1][2]),
		formatBoardSquare(board[1][3]),
		formatBoardSquare(board[1][4]),
		formatBoardSquare(board[2][0]),
		formatBoardSquare(board[2][1]),
		formatBoardSquare(board[2][2]),
		formatBoardSquare(board[2][3]),
		formatBoardSquare(board[2][4]),
		formatBoardSquare(board[3][0]),
		formatBoardSquare(board[3][1]),
		formatBoardSquare(board[3][2]),
		formatBoardSquare(board[3][3]),
		formatBoardSquare(board[3][4]),
		formatBoardSquare(board[4][0]),
		formatBoardSquare(board[4][1]),
		formatBoardSquare(board[4][2]),
		formatBoardSquare(board[4][3]),
		formatBoardSquare(board[4][4]),
	)
}

func main() {
	board := hexit.NewBoard()
	player := byte(1)
	for hexit.GetWinner(board) == 0 {
		PrintBoard(&board)

		row := uint(0)
		col := uint(0)
		_, err := fmt.Scanf("%d,%d", &row, &col)
		if err != nil || row < 0 || row >= 5 || col < 0 || col >= 5 || board[row][col] != 0 {
			fmt.Println("Invalid move!")
			continue
		}

		board = hexit.PlayMove(board, player, row, col)

		if hexit.GetWinner(board) != 0 {
			break
		}

		otherPlayer := hexit.OtherPlayer(player)
		tree := hexit.NewSearchTree(board, otherPlayer)
		for i := 0; i < 1000; i++ {
			hexit.DoVisit(&tree)
		}
		bestMove := hexit.GetBestMove(&tree)
		board = hexit.PlayMove(board, otherPlayer, bestMove.Row, bestMove.Col)
	}

	fmt.Printf("Player %d wins!\n", hexit.GetWinner(board))
}

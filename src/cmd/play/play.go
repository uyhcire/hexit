package main

import (
	"fmt"

	hexit "github.com/uyhcire/hexit/src"
)

func main() {
	board := hexit.NewBoard()
	player := byte(1)
	for hexit.GetWinner(board) == 0 {
		hexit.PrintBoard(&board)

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

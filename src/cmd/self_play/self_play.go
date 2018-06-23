package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/uyhcire/hexit/src"
)

var numVisits = 800

func getSelfPlayMove(board hexit.Board, player byte) hexit.Move {
	tree := hexit.NewSearchTree(board, player)
	for i := 0; i < numVisits; i++ {
		hexit.DoVisit(&tree)
	}
	bestMove := hexit.GetBestMove(&tree)
	return bestMove
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	board := hexit.NewBoard()
	player := byte(1)

	for hexit.GetWinner(board) == 0 {
		fmt.Println("")
		move := getSelfPlayMove(board, player)
		board = hexit.PlayMove(board, player, move.Row, move.Col)
		player = hexit.OtherPlayer(player)
		hexit.PrintBoard(&board)
	}

	fmt.Printf("Player %d wins!\n", hexit.GetWinner(board))
}

package main

import (
	"fmt"
	"math/rand"
	"time"

	hexit "github.com/uyhcire/hexit/src"
)

func playMatchGame() byte {
	board := hexit.NewBoard()
	player := byte(1)
	for hexit.GetWinner(board) == 0 {
		hexit.PrintBoard(&board)
		fmt.Println("")
		var evaluatePosition hexit.Evaluator
		if player == 1 {
			evaluatePosition = hexit.EvaluatePositionRandomly
		} else {
			evaluatePosition = hexit.EvaluatePositionWithNN
		}

		tree := hexit.NewSearchTree(evaluatePosition, board, player)
		for i := 0; i < 1000; i++ {
			hexit.DoVisit(&tree, evaluatePosition)
		}
		bestMove := hexit.GetBestMove(&tree)
		board = hexit.PlayMove(board, player, bestMove.Row, bestMove.Col)

		player = hexit.OtherPlayer(player)

		time.Sleep(time.Second)
	}

	winner := hexit.GetWinner(board)
	fmt.Printf("Player %d wins!\n", winner)
	return winner
}

func main() {
	hexit.InitializeModel()

	rand.Seed(time.Now().UTC().UnixNano())

	playerTwoWinCount := 0
	for i := 0; i < 100; i++ {
		winner := playMatchGame()
		if winner == 2 {
			playerTwoWinCount++
		}
		fmt.Printf("Player 2 won %d/%d games.\n", playerTwoWinCount, i+1)
	}
}

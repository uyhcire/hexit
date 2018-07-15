package main

import (
	"fmt"
	"math/rand"
	"time"

	hexit "github.com/uyhcire/hexit/src"
)

func playMatchGame() byte {
	var err error
	game := hexit.NewGame()
	for hexit.GetWinner(game.Board) == 0 {
		if game.MoveNum == 2 {
			err, game = hexit.DoNotSwitchSides(game)
			if err != nil {
				panic(err)
			}
			continue
		}

		hexit.PrintBoard(&game.Board)
		fmt.Println("")
		var evaluatePosition hexit.Evaluator
		if hexit.GetOriginalPlayer(game) == 1 {
			evaluatePosition = hexit.EvaluatePositionRandomly
		} else {
			evaluatePosition = hexit.EvaluatePositionWithNN
		}

		tree := hexit.NewSearchTree(evaluatePosition, game)
		for i := 0; i < 1000; i++ {
			hexit.DoVisit(&tree, evaluatePosition)
		}
		bestMove := hexit.GetBestMove(&tree)
		err, game = hexit.PlayGameMove(game, bestMove.Row, bestMove.Col)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Second)
	}

	winner := hexit.GetWinner(game.Board)
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

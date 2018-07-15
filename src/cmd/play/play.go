package main

import (
	"errors"
	"fmt"

	hexit "github.com/uyhcire/hexit/src"
)

func getHumanMove(board hexit.Board) (error, hexit.Move) {
	hexit.PrintBoard(&board)
	row := uint(0)
	col := uint(0)
	_, err := fmt.Scanf("%d,%d", &row, &col)
	if err != nil || row < 0 || row >= 5 || col < 0 || col >= 5 || board[row][col] != 0 {
		return errors.New("Invalid move"), hexit.Move{Row: 0, Col: 0}
	}
	return nil, hexit.Move{Row: row, Col: col}
}

func main() {
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

		var move hexit.Move
		if hexit.GetOriginalPlayer(game) == 1 {
			err, move = getHumanMove(game.Board)
			if err != nil {
				fmt.Println("Invalid move!")
				continue
			}
		} else {
			tree := hexit.NewSearchTree(hexit.EvaluatePositionRandomly, game)
			for i := 0; i < 1000; i++ {
				hexit.DoVisit(&tree, hexit.EvaluatePositionRandomly)
			}
			move = hexit.GetBestMove(&tree)
		}

		err, game = hexit.PlayGameMove(game, move.Row, move.Col)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Player %d wins!\n", hexit.GetWinner(game.Board))
}

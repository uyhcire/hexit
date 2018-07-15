package hexit

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"
)

var numVisits = 800

func PlaySelfPlayGame(outputFilename string) {
	rand.Seed(time.Now().UTC().UnixNano())

	var err error
	game := NewGame()

	moveSnapshots := make([]*TrainingGame_MoveSnapshot, 0)

	for GetWinner(game.Board) == 0 {
		if game.MoveNum == 2 {
			err, game = DoNotSwitchSides(game)
			if err != nil {
				panic(err)
			}
			continue
		}

		squaresOccupiedByMyself, squaresOccupiedByOtherPlayer := GetOccupiedSquaresForNN(game.Board, game.CurrentPlayer)
		moveSnapshot := TrainingGame_MoveSnapshot{
			NormalizedVisitCounts:        nil,
			Winner:                       TrainingGame_MYSELF,
			SquaresOccupiedByMyself:      squaresOccupiedByMyself,
			SquaresOccupiedByOtherPlayer: squaresOccupiedByOtherPlayer,
		}

		tree := NewSearchTree(EvaluatePositionRandomly, game)
		ApplyDirichletNoise(&tree)
		for i := 0; i < numVisits; i++ {
			DoVisit(&tree, EvaluatePositionRandomly)
		}

		normalizedVisitCounts := make([]float32, 5*5)
		for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
			row, col := childNode.move.Row, childNode.move.Col
			if game.Board[row][col] != 0 {
				panic("Illegal move")
			}
			normalizedVisitCounts[row*5+col] = float32(childNode.n) / float32(numVisits)
		}
		moveSnapshot.NormalizedVisitCounts = normalizedVisitCounts

		move := GetMoveWithTemperatureOne(&tree)
		err, game = PlayGameMove(game, move.Row, move.Col)
		if err != nil {
			panic(err)
		}

		moveSnapshots = append(moveSnapshots, &moveSnapshot)
	}

	winner := GetWinner(game.Board)

	// Fill in move snapshots with the winner
	player := byte(1)
	for _, moveSnapshot := range moveSnapshots {
		if player == winner {
			moveSnapshot.Winner = TrainingGame_MYSELF
		} else {
			moveSnapshot.Winner = TrainingGame_OTHER_PLAYER
		}
		player = OtherPlayer(player)
	}

	trainingGame := TrainingGame{MoveSnapshots: moveSnapshots}
	trainingGameBytes, err := proto.Marshal(&trainingGame)
	if err != nil {
		panic(err)
	}

	trainingDataPath := filepath.Join(".", "training_games")
	err = os.MkdirAll(trainingDataPath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(".", "training_games", outputFilename))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(trainingGameBytes)
	if err != nil {
		panic(err)
	}
}

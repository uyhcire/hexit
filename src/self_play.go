package hexit

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"
)

type moveSnapshotWithoutWinner struct {
	normalizedVisitCounts        []float32
	squaresOccupiedByMyself      []float32
	squaresOccupiedByOtherPlayer []float32
}

type trainingGameBuilder struct {
	moveSnapshots  []*moveSnapshotWithoutWinner
	didSwitchSides bool
}

func newTrainingGameBuilder() trainingGameBuilder {
	return trainingGameBuilder{
		moveSnapshots:  make([]*moveSnapshotWithoutWinner, 0),
		didSwitchSides: false,
	}
}

func recordTrainingGameMove(builder *trainingGameBuilder, game Game, normalizedVisitCounts []float32) {
	squaresOccupiedByMyself, squaresOccupiedByOtherPlayer := GetOccupiedSquaresForNN(game.Board, game.CurrentPlayer)
	builder.moveSnapshots = append(builder.moveSnapshots, &moveSnapshotWithoutWinner{
		normalizedVisitCounts:        normalizedVisitCounts,
		squaresOccupiedByMyself:      squaresOccupiedByMyself,
		squaresOccupiedByOtherPlayer: squaresOccupiedByOtherPlayer,
	})
}

func recordTrainingGameSwitchedSides(builder *trainingGameBuilder) {
	builder.didSwitchSides = true
}

func buildTrainingGame(builder *trainingGameBuilder, winner byte) TrainingGame {
	moveSnapshots := make([]*TrainingGame_MoveSnapshot, 0)
	player := byte(1)

	for i, moveSnapshotWithoutWinner := range builder.moveSnapshots {
		var trainingGameWinner TrainingGame_Player
		originalPlayer := player
		if i == 0 && builder.didSwitchSides {
			originalPlayer = OtherPlayer(originalPlayer)
		}
		if originalPlayer == winner {
			trainingGameWinner = TrainingGame_MYSELF
		} else {
			trainingGameWinner = TrainingGame_OTHER_PLAYER
		}

		moveSnapshots = append(moveSnapshots, &TrainingGame_MoveSnapshot{
			NormalizedVisitCounts:        moveSnapshotWithoutWinner.normalizedVisitCounts,
			Winner:                       trainingGameWinner,
			SquaresOccupiedByMyself:      moveSnapshotWithoutWinner.squaresOccupiedByMyself,
			SquaresOccupiedByOtherPlayer: moveSnapshotWithoutWinner.squaresOccupiedByOtherPlayer,
		})

		player = OtherPlayer(player)
	}

	return TrainingGame{MoveSnapshots: moveSnapshots}
}

var numVisits = 800

func playTrainingGame() TrainingGame {
	rand.Seed(time.Now().UTC().UnixNano())

	var err error
	game := NewGame()
	trainingGameBuilder := newTrainingGameBuilder()

	for GetWinner(game.Board) == 0 {
		tree := NewSearchTree(EvaluatePositionRandomly, game)
		ApplyDirichletNoise(&tree)
		for i := 0; i < numVisits; i++ {
			DoVisit(&tree, EvaluatePositionRandomly)
		}

		if game.MoveNum == 2 {
			if GetExpectedValueOfGame(&tree) > 0 {
				err, game = SwitchSides(game)
				recordTrainingGameSwitchedSides(&trainingGameBuilder)
			} else {
				err, game = DoNotSwitchSides(game)
			}
			if err != nil {
				panic(err)
			}
			continue
		}

		normalizedVisitCounts := make([]float32, 5*5)
		for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
			row, col := childNode.move.Row, childNode.move.Col
			if game.Board[row][col] != 0 {
				panic("Illegal move")
			}
			normalizedVisitCounts[row*5+col] = float32(childNode.n) / float32(numVisits)
		}
		recordTrainingGameMove(&trainingGameBuilder, game, normalizedVisitCounts)

		move := GetMoveWithTemperatureOne(&tree)
		err, game = PlayGameMove(game, move.Row, move.Col)
		if err != nil {
			panic(err)
		}
	}

	winner := GetWinner(game.Board)
	return buildTrainingGame(&trainingGameBuilder, winner)
}

func GenerateTrainingGame(outputFilename string) {
	trainingGame := playTrainingGame()

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

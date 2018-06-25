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

	board := NewBoard()
	player := byte(1)

	moveSnapshots := make([]*TrainingGame_MoveSnapshot, 0)

	for GetWinner(board) == 0 {
		moveSnapshot := TrainingGame_MoveSnapshot{
			//TODO:serialize visit counts
			NormalizedVisitCounts:        nil,
			Winner:                       TrainingGame_MYSELF,
			SquaresOccupiedByMyself:      nil,
			SquaresOccupiedByOtherPlayer: nil,
		}
		WriteBoardToMoveSnapshot(&moveSnapshot, board, player)

		tree := NewSearchTree(board, player)
		ApplyDirichletNoise(&tree)
		for i := 0; i < numVisits; i++ {
			DoVisit(&tree)
		}

		normalizedVisitCounts := make([]float32, 5*5)
		for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
			row, col := childNode.move.Row, childNode.move.Col
			if board[row][col] != 0 {
				panic("Illegal move")
			}
			normalizedVisitCounts[row*5+col] = float32(childNode.n) / float32(numVisits)
		}
		moveSnapshot.NormalizedVisitCounts = normalizedVisitCounts

		move := GetMoveWithTemperatureOne(&tree)
		board = PlayMove(board, player, move.Row, move.Col)

		moveSnapshots = append(moveSnapshots, &moveSnapshot)

		player = OtherPlayer(player)
	}

	winner := GetWinner(board)

	// Fill in move snapshots with the winner
	player = 1
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

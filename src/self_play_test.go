package hexit

import "testing"

func TestNewTrainingGameBuilder(t *testing.T) {
	builder := newTrainingGameBuilder()
	if len(builder.moveSnapshots) != 0 {
		t.Error("New builder should have no moves")
	}
}

func makeUniformVisitCounts(board Board) []float32 {
	numLegalMoves := 0
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if board[i][j] == 0 {
				numLegalMoves++
			}
		}
	}

	normalizedVisitCounts := make([]float32, 5*5)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if board[i][j] == 0 {
				normalizedVisitCounts[5*i+j] = 1 / float32(numLegalMoves)
			}
		}
	}

	return normalizedVisitCounts
}

func TestBuildSimpleTrainingGame(t *testing.T) {
	game := NewGame()
	builder := newTrainingGameBuilder()

	recordTrainingGameMove(&builder, game, makeUniformVisitCounts(game.Board))
	trainingGame := buildTrainingGame(&builder, 1)

	if trainingGame.MoveSnapshots[0].Winner != TrainingGame_MYSELF {
		t.Errorf("Expected Player 1 to win")
	}
}

func TestBuildTrainingGameWithSideSwitching(t *testing.T) {
	game := NewGame()
	builder := newTrainingGameBuilder()

	recordTrainingGameMove(&builder, game, makeUniformVisitCounts(game.Board))
	err, game := PlayGameMove(game, 0, 0)
	if err != nil {
		t.Error(err.Error())
	}

	recordTrainingGameSwitchedSides(&builder)
	err, game = SwitchSides(game)
	if err != nil {
		t.Error(err.Error())
	}

	recordTrainingGameMove(&builder, game, makeUniformVisitCounts(game.Board))
	// Player 2's color (played by Player 1) wins the game
	trainingGame := buildTrainingGame(&builder, 2)

	if trainingGame.MoveSnapshots[0].Winner != TrainingGame_MYSELF {
		t.Errorf("Expected Player 1 to win")
	}
	if trainingGame.MoveSnapshots[1].Winner != TrainingGame_MYSELF {
		t.Errorf("Expected Player 1 to win (playing as Player 2's color)")
	}
}

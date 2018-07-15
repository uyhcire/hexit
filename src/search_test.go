package hexit

import (
	"testing"
)

func TestNewSearchTree(t *testing.T) {
	game := NewGame()
	tree := NewSearchTree(EvaluatePositionUniformly, game)
	if tree.rootNode.firstChild == nil {
		t.Error("Root node should have children attached!")
	}
}

func TestNumLegalMoves(t *testing.T) {
	game := NewGame()
	err, game := PlayGameMove(game, 0, 0)
	if err != nil {
		t.Error(err.Error())
	}
	tree := NewSearchTree(EvaluatePositionUniformly, game)

	numLegalMoves := 0
	for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
		numLegalMoves++
	}

	if numLegalMoves != 24 {
		t.Errorf("Expected 24 legal moves, but got %d (25 locations, 1 already filled)", numLegalMoves)
	}
}

func TestDoVisit(t *testing.T) {
	game := NewGame()
	// Player 1 can win by playing at (4, 0)
	game.Board = [5][5]byte{
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	winningMoveNode := (*SearchNode)(nil)
	for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
		if childNode.move.Row == 4 && childNode.move.Col == 0 {
			winningMoveNode = childNode
		}
	}
	if winningMoveNode == nil {
		t.Error("Winning move should be in search tree!")
	}
	if winningMoveNode.n <= 800 {
		t.Error("Expected winning move to have the vast majority of visits")
	}
}

func TestGetBestMovePlayerOne(t *testing.T) {
	game := NewGame()
	// Player 1 can win by playing at (4, 0)
	game.Board = [5][5]byte{
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	bestMove := GetBestMove(&tree)
	if bestMove.Row != 4 || bestMove.Col != 0 {
		t.Error("Failed to find the winning move")
	}
}

func TestGetBestMovePlayerTwo(t *testing.T) {
	game := NewGame()
	// Player 2 can win by playing at (0, 4)
	game.CurrentPlayer = 2
	game.Board = [5][5]byte{
		[5]byte{2, 2, 2, 2, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	bestMove := GetBestMove(&tree)
	if bestMove.Row != 0 || bestMove.Col != 4 {
		t.Error("Failed to find the winning move")
	}
}

// newGameWithSideSwitching creates a game where Player 2 can win by switching sides.
func newGameWithSideSwitching() Game {
	game := NewGame()
	game.MoveNum = 1
	game.CurrentPlayer = 1
	// If Player 1 plays anywhere other than (3, 0), Player 2 can switch sides and immediately win.
	game.Board = [5][5]byte{
		[5]byte{1, 1, 1, 1, 1},
		[5]byte{1, 1, 1, 1, 1},
		[5]byte{1, 1, 1, 1, 1},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{2, 0, 0, 0, 0},
	}
	return game
}

func TestEvalWithSideSwitching(t *testing.T) {
	game := newGameWithSideSwitching()
	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	bestMove := GetBestMove(&tree)
	if bestMove.Row != 3 || bestMove.Col != 0 {
		PrintVisitDistribution(tree.rootNode)
		t.Error("Failed to find the only non-losing move")
	}
}

func TestGetExpectedValueOfGame(t *testing.T) {
	game := NewGame()
	tree := NewSearchTree(EvaluatePositionUniformly, game)
	DoVisit(&tree, EvaluatePositionUniformly)

	expectedValue := GetExpectedValueOfGame(&tree)
	if expectedValue != 0 {
		t.Error("Expected value of game should be 0 initially")
	}
}

func TestGetExpectedValueOfGameWithSideSwitching(t *testing.T) {
	game := newGameWithSideSwitching()
	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	expectedValue := GetExpectedValueOfGame(&tree)
	if expectedValue > -0.1 {
		t.Error("Player 2 should be expected to win")
	}
}

func TestGetExpectedValueOfGameWithSideSwitchingAsPlayerTwo(t *testing.T) {
	game := newGameWithSideSwitching()
	err, game := PlayGameMove(game, 3, 1)
	if err != nil {
		t.Error(err.Error())
	}

	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	expectedValue := GetExpectedValueOfGame(&tree)
	if expectedValue > -0.8 {
		t.Error("Player 2 can switch sides and immediately win!")
	}
}

func TestGetExpectedValueOfGamePlayerOneHasWinningFirstMove(t *testing.T) {
	game := NewGame()
	// Player 1 can win by playing at (4, 0)
	game.MoveNum = 1
	game.Board = [5][5]byte{
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	expectedValue := GetExpectedValueOfGame(&tree)
	if expectedValue < 0.8 {
		t.Error("Player 1 has a winning move!")
	}
}

func TestGetExpectedValueOfGamePlayerOneHasWinningMove(t *testing.T) {
	game := NewGame()
	// Player 1 can win by playing at (4, 0)
	game.MoveNum = 5
	game.Board = [5][5]byte{
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	expectedValue := GetExpectedValueOfGame(&tree)
	if expectedValue < 0.8 {
		t.Error("Player 1 has a winning move!")
	}
}

func TestGetExpectedValueOfGamePlayerTwoHasWinningMove(t *testing.T) {
	game := NewGame()
	// Player 2 can win by playing at (0, 4)
	game.MoveNum = 5
	game.CurrentPlayer = 2
	game.Board = [5][5]byte{
		[5]byte{2, 2, 2, 2, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	tree := NewSearchTree(EvaluatePositionUniformly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionUniformly)
	}

	expectedValue := GetExpectedValueOfGame(&tree)
	if expectedValue > -0.8 {
		t.Error("Player 2 has a winning move!")
	}
}

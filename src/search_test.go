package hexit

import (
	"math/rand"
	"testing"
)

func TestNewSearchTree(t *testing.T) {
	game := NewGame()
	tree := NewSearchTree(EvaluatePositionRandomly, game)
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
	tree := NewSearchTree(EvaluatePositionRandomly, game)

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

	rand.Seed(1)

	tree := NewSearchTree(EvaluatePositionRandomly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionRandomly)
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

	rand.Seed(1)

	tree := NewSearchTree(EvaluatePositionRandomly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionRandomly)
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

	rand.Seed(1)

	tree := NewSearchTree(EvaluatePositionRandomly, game)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree, EvaluatePositionRandomly)
	}

	bestMove := GetBestMove(&tree)
	if bestMove.Row != 0 || bestMove.Col != 4 {
		t.Error("Failed to find the winning move")
	}
}

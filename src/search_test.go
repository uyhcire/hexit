package hexit

import (
	"math/rand"
	"testing"
)

func TestNewSearchTree(t *testing.T) {
	board := NewBoard()
	tree := NewSearchTree(board, 1)
	if tree.rootNode.firstChild == nil {
		t.Error("Root node should have children attached!")
	}
}

func TestNumLegalMoves(t *testing.T) {
	board := NewBoard()
	board = PlayMove(board, 1, 0, 0)
	tree := NewSearchTree(board, 2)

	numLegalMoves := 0
	for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
		numLegalMoves++
	}

	if numLegalMoves != 24 {
		t.Errorf("Expected 24 legal moves, but got %d (25 locations, 1 already filled)", numLegalMoves)
	}
}

func TestDoVisit(t *testing.T) {
	almostWonBoard := [5][5]byte{
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	rand.Seed(1)

	tree := NewSearchTree(almostWonBoard, 1)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree)
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
	almostWonBoard := [5][5]byte{
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{1, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	rand.Seed(1)

	tree := NewSearchTree(almostWonBoard, 1)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree)
	}

	bestMove := GetBestMove(&tree)
	if bestMove.Row != 4 || bestMove.Col != 0 {
		t.Error("Failed to find the winning move")
	}
}

func TestGetBestMovePlayerTwo(t *testing.T) {
	almostWonBoard := [5][5]byte{
		[5]byte{2, 2, 2, 2, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
		[5]byte{0, 0, 0, 0, 0},
	}

	rand.Seed(1)

	tree := NewSearchTree(almostWonBoard, 2)
	for i := 0; i < 1000; i++ {
		DoVisit(&tree)
	}

	bestMove := GetBestMove(&tree)
	if bestMove.Row != 0 || bestMove.Col != 4 {
		t.Error("Failed to find the winning move")
	}
}

package hex

import "testing"

func TestNewBoard(t *testing.T) {
	board := NewBoard()
	if board[0][0] != 0 {
		t.Error("Board should be empty")
	}
}

func TestPlayMove(t *testing.T) {
	board := NewBoard()
	board = PlayMove(board, 1, 0, 0)
	if board[0][0] != 1 {
		t.Error("Expected Player 1 to fill a grid location")
	}
}

/*
 X - - - -
  X - - - -
   X - - - -
    X - - - -
     X - - - -
*/
func TestPlayerOneWins(t *testing.T) {
	board := NewBoard()
	board = PlayMove(board, 1, 0, 0)
	board = PlayMove(board, 1, 1, 0)
	board = PlayMove(board, 1, 2, 0)
	board = PlayMove(board, 1, 3, 0)
	if PlayerOneWins(board) {
		t.Error("Player 1 hasn't won yet!")
	}
	board = PlayMove(board, 1, 4, 0)
	if !PlayerOneWins(board) {
		t.Error("Expected Player 1 to be the winner")
	}
}

/*
 X - - - -
  X - X X X
   X X - - X
    - - - X -
     - - - X -
*/
func TestPlayerOneWinsWithWindingPath(t *testing.T) {
	board := NewBoard()
	board = PlayMove(board, 1, 0, 0)
	board = PlayMove(board, 1, 1, 0)
	board = PlayMove(board, 1, 2, 0)
	board = PlayMove(board, 1, 2, 1)
	board = PlayMove(board, 1, 1, 2)
	board = PlayMove(board, 1, 1, 3)
	board = PlayMove(board, 1, 1, 4)
	board = PlayMove(board, 1, 2, 4)
	board = PlayMove(board, 1, 3, 3)
	if PlayerOneWins(board) {
		t.Error("Player 1 hasn't won yet!")
	}
	board = PlayMove(board, 1, 4, 3)
	if !PlayerOneWins(board) {
		t.Error("Expected Player 1 to be the winner")
	}
}

/*
 O O O O O
  - - - - -
   - - - - -
    - - - - -
     - - - - -
*/
func TestPlayerTwoWins(t *testing.T) {
	board := NewBoard()
	board = PlayMove(board, 2, 0, 0)
	board = PlayMove(board, 2, 0, 1)
	board = PlayMove(board, 2, 0, 2)
	board = PlayMove(board, 2, 0, 3)
	if PlayerTwoWins(board) {
		t.Error("Player 2 hasn't won yet!")
	}
	board = PlayMove(board, 2, 0, 4)
	if !PlayerTwoWins(board) {
		t.Error("Expected Player 2 to be the winner")
	}
}

func TestGetWinnerEmptyBoard(t *testing.T) {
	board := NewBoard()
	if GetWinner(board) != 0 {
		t.Error("The game just started, there's no winner yet!")
	}
}

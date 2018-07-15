package hexit

import "testing"

func TestNewGame(t *testing.T) {
	game := NewGame()
	if game.MoveNum != 1 {
		t.Error("Expected game to start on move 1")
	}
}

func TestPlayGameMove(t *testing.T) {
	game := NewGame()
	err, game := PlayGameMove(game, 0, 0)
	if err != nil {
		t.Error(err.Error())
	}
	if game.Board[0][0] != 1 {
		t.Error("Expected square to be occupied by Player 1")
	}
}

func TestSwitchSides(t *testing.T) {
	game := NewGame()
	err, game := PlayGameMove(game, 0, 0)
	if err != nil {
		t.Error(err.Error())
	}

	err, game = SwitchSides(game)
	if err != nil {
		t.Error(err.Error())
	}

	if game.CurrentPlayer != 1 {
		t.Error("If Player 2 switched sides, they should now be playing as Player 1")
	}
	if game.MoveNum != 3 {
		t.Error("Switching sides should count as a move")
	}
	if game.SwitchedSides != true {
		t.Error("Expected to switch sides!")
	}
}

func TestDoNotSwitchSides(t *testing.T) {
	game := NewGame()
	err, game := PlayGameMove(game, 0, 0)
	if err != nil {
		t.Error(err.Error())
	}

	err, game = DoNotSwitchSides(game)
	if err != nil {
		t.Error(err.Error())
	}

	if game.CurrentPlayer != 2 {
		t.Error("If Player 2 did not switch sides, they should still be playing as Player 2")
	}
	if game.MoveNum != 3 {
		t.Error("Deciding not to switch sides should count as a move")
	}
	if game.SwitchedSides != false {
		t.Error("Expected not to switch sides!")
	}
}

func TestGetOriginalPlayer(t *testing.T) {
	game := NewGame()
	err, game := PlayGameMove(game, 0, 0)
	if err != nil {
		t.Error(err.Error())
	}

	err, game = SwitchSides(game)
	if err != nil {
		t.Error(err.Error())
	}

	if GetOriginalPlayer(game) != 2 {
		t.Error("Original player should be Player 2 even if they decided to switch sides")
	}
}

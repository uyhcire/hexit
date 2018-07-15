package hexit

import "errors"

// Game represents the state of a game
type Game struct {
	CurrentPlayer byte
	// The game starts on move 1, and each player's turn is a new move
	MoveNum int
	// Player 2 has the option to switch sides on their first move (the "pie rule").
	// They switch sides on move 2, and their next move is considered move 3.
	SwitchedSides bool
	Board         Board
}

// NewGame creates a new game state
func NewGame() Game {
	return Game{
		CurrentPlayer: 1,
		MoveNum:       1,
		SwitchedSides: false,
		Board:         NewBoard(),
	}
}

// PlayGameMove plays a regular move
func PlayGameMove(game Game, row uint, col uint) (error, Game) {
	if game.MoveNum == 2 {
		return errors.New("Move 2 is when Player 2 decides whether to switch sides"), game
	}

	return nil, Game{
		CurrentPlayer: OtherPlayer(game.CurrentPlayer),
		MoveNum:       game.MoveNum + 1,
		SwitchedSides: game.SwitchedSides,
		Board:         PlayMove(game.Board, game.CurrentPlayer, row, col),
	}
}

// SwitchSides switches sides on Player 2's first move
func SwitchSides(game Game) (error, Game) {
	if game.MoveNum != 2 {
		return errors.New("Can only switch sides on move 2"), game
	}

	return nil, Game{
		CurrentPlayer: OtherPlayer(game.CurrentPlayer),
		MoveNum:       game.MoveNum + 1,
		SwitchedSides: true,
		Board:         game.Board,
	}
}

// DoNotSwitchSides: Player 2 decided not to switch sides
func DoNotSwitchSides(game Game) (error, Game) {

	if game.MoveNum != 2 {
		return errors.New("Can only switch sides on move 2"), game
	}

	return nil, Game{
		CurrentPlayer: game.CurrentPlayer,
		MoveNum:       game.MoveNum + 1,
		SwitchedSides: false,
		Board:         game.Board,
	}
}

// GetOriginalPlayer: get the "original" player, ignoring side-switching.
// For example, if Player 2 switched sides, they play their first move as Player 1, but GetOriginalPlayer still returns Player 2.
func GetOriginalPlayer(game Game) byte {
	if game.SwitchedSides {
		return OtherPlayer(game.CurrentPlayer)
	}
	return game.CurrentPlayer
}

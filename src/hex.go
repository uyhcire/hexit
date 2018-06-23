package hexit

import (
	"fmt"
)

// Board encoding: 0 = blank, 1 = Player 1, 2 = Player 2
// First index is the number of rows from the top
// Second index is the number of columns from the left
type Board = [5][5]byte

// BoardLocation is a location on a board
type BoardLocation struct {
	Row uint
	Col uint
}

// Move that a player can make
type Move = BoardLocation

func formatBoardSquare(boardSquareValue byte) string {
	if boardSquareValue == 1 {
		return "X"
	} else if boardSquareValue == 2 {
		return "O"
	} else {
		return "-"
	}
}

// PrintBoard prints the board to the console
func PrintBoard(board *Board) {
	fmt.Printf(
		"%s %s %s %s %s\n"+
			" %s %s %s %s %s\n"+
			"  %s %s %s %s %s\n"+
			"   %s %s %s %s %s\n"+
			"    %s %s %s %s %s\n",
		formatBoardSquare(board[0][0]),
		formatBoardSquare(board[0][1]),
		formatBoardSquare(board[0][2]),
		formatBoardSquare(board[0][3]),
		formatBoardSquare(board[0][4]),
		formatBoardSquare(board[1][0]),
		formatBoardSquare(board[1][1]),
		formatBoardSquare(board[1][2]),
		formatBoardSquare(board[1][3]),
		formatBoardSquare(board[1][4]),
		formatBoardSquare(board[2][0]),
		formatBoardSquare(board[2][1]),
		formatBoardSquare(board[2][2]),
		formatBoardSquare(board[2][3]),
		formatBoardSquare(board[2][4]),
		formatBoardSquare(board[3][0]),
		formatBoardSquare(board[3][1]),
		formatBoardSquare(board[3][2]),
		formatBoardSquare(board[3][3]),
		formatBoardSquare(board[3][4]),
		formatBoardSquare(board[4][0]),
		formatBoardSquare(board[4][1]),
		formatBoardSquare(board[4][2]),
		formatBoardSquare(board[4][3]),
		formatBoardSquare(board[4][4]),
	)
}

// OtherPlayer returns the other player
func OtherPlayer(player byte) byte {
	if player == 1 {
		return 2
	} else if player == 2 {
		return 1
	} else {
		panic("Invalid player")
	}
}

// NewBoard creates an empty board
func NewBoard() Board {
	return [5][5]byte{}
}

// CopyBoard returns a copy of an existing board
func CopyBoard(board Board) Board {
	newBoard := NewBoard()
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			newBoard[i][j] = board[i][j]
		}
	}
	return newBoard
}

// PlayMove returns a copy of the board after playing a move.
func PlayMove(board Board, player byte, row uint, col uint) Board {
	newBoard := CopyBoard(board)
	if newBoard[row][col] != 0 {
		panic("Location is already occupied")
	}
	newBoard[row][col] = player
	return newBoard
}

func getAdjacentLocations(location BoardLocation) []BoardLocation {
	adjacent := make([]BoardLocation, 0)

	// Row above
	adjacent = append(adjacent, BoardLocation{
		Row: location.Row - 1,
		Col: location.Col,
	})
	adjacent = append(adjacent, BoardLocation{
		Row: location.Row - 1,
		Col: location.Col + 1,
	})

	// Same row
	adjacent = append(adjacent, BoardLocation{
		Row: location.Row,
		Col: location.Col - 1,
	})
	adjacent = append(adjacent, BoardLocation{
		Row: location.Row,
		Col: location.Col + 1,
	})

	// Next row
	adjacent = append(adjacent, BoardLocation{
		Row: location.Row + 1,
		Col: location.Col - 1,
	})
	adjacent = append(adjacent, BoardLocation{
		Row: location.Row + 1,
		Col: location.Col,
	})

	validAdjacent := make([]BoardLocation, 0)
	for _, adjacentLocation := range adjacent {
		if adjacentLocation.Row >= 0 && adjacentLocation.Row < 5 && adjacentLocation.Col >= 0 && adjacentLocation.Col < 5 {
			validAdjacent = append(validAdjacent, adjacentLocation)
		}
	}

	return validAdjacent
}

// Do breadth-first search to determine if the player has a connected path from one side of the board to the other.
func didPlayerWin(
	board Board,
	player byte,
	// Side of the board to start looking from
	startingLocations []BoardLocation,
	// Is this location on the other side of the board?
	isLocationWinning func(BoardLocation) bool,
) bool {
	visited := [5][5]bool{}
	locationQueue := make([]BoardLocation, 0)
	for _, location := range startingLocations {
		row := location.Row
		col := location.Col
		if board[row][col] == player {
			locationQueue = append(locationQueue, location)
			visited[row][col] = true
		}
	}

	for len(locationQueue) != 0 {
		location := locationQueue[0]
		locationQueue = locationQueue[1:]
		for _, adjacentLocation := range getAdjacentLocations(location) {
			if board[adjacentLocation.Row][adjacentLocation.Col] != player {
				continue
			}
			if isLocationWinning(adjacentLocation) {
				return true
			}
			if !visited[adjacentLocation.Row][adjacentLocation.Col] {
				visited[adjacentLocation.Row][adjacentLocation.Col] = true
				locationQueue = append(locationQueue, adjacentLocation)
			}
		}
	}

	return false
}

// PlayerOneWins returns true if Player 1 has connected the top to the bottom
func PlayerOneWins(board Board) bool {
	return didPlayerWin(
		board,
		1,
		[]BoardLocation{
			BoardLocation{Row: 0, Col: 0},
			BoardLocation{Row: 0, Col: 1},
			BoardLocation{Row: 0, Col: 2},
			BoardLocation{Row: 0, Col: 3},
			BoardLocation{Row: 0, Col: 4},
		},
		func(location BoardLocation) bool {
			return location.Row == 4
		})
}

// PlayerTwoWins returns true if Player 2 has connected the left to the right
func PlayerTwoWins(board Board) bool {
	return didPlayerWin(
		board,
		2,
		[]BoardLocation{
			BoardLocation{Row: 0, Col: 0},
			BoardLocation{Row: 1, Col: 0},
			BoardLocation{Row: 2, Col: 0},
			BoardLocation{Row: 3, Col: 0},
			BoardLocation{Row: 4, Col: 0},
		},
		func(location BoardLocation) bool {
			return location.Col == 4
		})
}

// GetWinner returns the player that won the game, or 0 if the game is still in progress.
func GetWinner(board Board) byte {
	if PlayerOneWins(board) {
		return 1
	} else if PlayerTwoWins(board) {
		return 2
	} else {
		return 0
	}
}

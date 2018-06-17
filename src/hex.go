package hex

// Board encoding: 0 = blank, 1 = Player 1, 2 = Player 2
// First index is the number of rows from the top
// Second index is the number of columns from the left
type Board = [5][5]byte

type boardLocation struct {
	row int
	col int
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

func getAdjacentLocations(location boardLocation) []boardLocation {
	adjacent := make([]boardLocation, 0)

	// Row above
	adjacent = append(adjacent, boardLocation{
		row: location.row - 1,
		col: location.col,
	})
	adjacent = append(adjacent, boardLocation{
		row: location.row - 1,
		col: location.col + 1,
	})

	// Same row
	adjacent = append(adjacent, boardLocation{
		row: location.row,
		col: location.col - 1,
	})
	adjacent = append(adjacent, boardLocation{
		row: location.row,
		col: location.col + 1,
	})

	// Next row
	adjacent = append(adjacent, boardLocation{
		row: location.row + 1,
		col: location.col - 1,
	})
	adjacent = append(adjacent, boardLocation{
		row: location.row + 1,
		col: location.col,
	})

	validAdjacent := make([]boardLocation, 0)
	for _, adjacentLocation := range adjacent {
		if adjacentLocation.row >= 0 && adjacentLocation.row < 5 && adjacentLocation.col >= 0 && adjacentLocation.col < 5 {
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
	startingLocations []boardLocation,
	// Is this location on the other side of the board?
	isLocationWinning func(boardLocation) bool,
) bool {
	visited := [5][5]bool{}
	locationQueue := make([]boardLocation, 0)
	for _, location := range startingLocations {
		row := location.row
		col := location.col
		if board[row][col] == player {
			locationQueue = append(locationQueue, location)
			visited[row][col] = true
		}
	}

	for len(locationQueue) != 0 {
		location := locationQueue[0]
		locationQueue = locationQueue[1:]
		for _, adjacentLocation := range getAdjacentLocations(location) {
			if board[adjacentLocation.row][adjacentLocation.col] != player {
				continue
			}
			if isLocationWinning(adjacentLocation) {
				return true
			}
			if !visited[adjacentLocation.row][adjacentLocation.col] {
				visited[adjacentLocation.row][adjacentLocation.col] = true
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
		[]boardLocation{
			boardLocation{row: 0, col: 0},
			boardLocation{row: 0, col: 1},
			boardLocation{row: 0, col: 2},
			boardLocation{row: 0, col: 3},
			boardLocation{row: 0, col: 4},
		},
		func(location boardLocation) bool {
			return location.row == 4
		})
}

// PlayerTwoWins returns true if Player 2 has connected the left to the right
func PlayerTwoWins(board Board) bool {
	return didPlayerWin(
		board,
		2,
		[]boardLocation{
			boardLocation{row: 0, col: 0},
			boardLocation{row: 1, col: 0},
			boardLocation{row: 2, col: 0},
			boardLocation{row: 3, col: 0},
			boardLocation{row: 4, col: 0},
		},
		func(location boardLocation) bool {
			return location.col == 4
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

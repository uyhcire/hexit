package hexit

import (
	"math"
	"math/rand"
)

// SearchNode is a node in a search tree
type SearchNode struct {
	// Move that led to this node
	move Move
	// Player that moves next
	player byte

	// Total number of visits
	n uint32
	// Value estimate from NN
	v float32
	// Average value
	q float32
	// Total value across all visits
	w float32
	// Policy estimate from NN
	p float32
	// Does this move end the game?
	isTerminal bool

	// Other nodes
	parent      *SearchNode
	firstChild  *SearchNode
	nextSibling *SearchNode
}

// SearchTree is an MCTS search tree
type SearchTree struct {
	board    Board
	rootNode *SearchNode
}

// NewSearchNode creates a new SearchNode
func NewSearchNode(parent *SearchNode, move Move, player byte) SearchNode {
	nan := float32(math.NaN())
	return SearchNode{
		move:        move,
		player:      player,
		n:           0,
		v:           nan,
		q:           0,
		w:           0,
		p:           nan,
		isTerminal:  false,
		parent:      parent,
		firstChild:  nil,
		nextSibling: nil,
	}
}

// NewSearchTree creates a new SearchTree
func NewSearchTree(board Board, player byte) SearchTree {
	if GetWinner(board) != 0 {
		panic("Can't search from a terminal node")
	}

	rootNode := NewSearchNode(nil, Move{Row: 1000, Col: 1000}, player)
	searchTree := SearchTree{
		board:    board,
		rootNode: &rootNode,
	}
	EvaluateAtNode(searchTree.rootNode, board)
	return searchTree
}

// EvaluatePosition returns the NN's value and policy estimates for a position.
func EvaluatePosition(board Board) (float32, [5][5]float32) {
	//TODO:should be deterministic
	valueEstimate := rand.Float32()*2 - 1
	policyEstimates := [5][5]float32{}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			policyEstimates[i][j] = rand.Float32()
		}
	}
	return valueEstimate, policyEstimates
}

// EvaluateAtNode evaluates the NN at a single node.
func EvaluateAtNode(node *SearchNode, board Board) {
	if node.isTerminal {
		panic("Should not evaluate the NN at a terminal node")
	}

	valueEstimate, policyEstimates := EvaluatePosition(board)
	node.v = valueEstimate

	firstChildNode := (*SearchNode)(nil)
	totalLegalPolicy := float32(0.0)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if board[i][j] != 0 {
				// Illegal move
				continue
			}
			totalLegalPolicy += policyEstimates[i][j]
			childNode := NewSearchNode(node, Move{Row: uint(i), Col: uint(j)}, OtherPlayer(node.player))
			childNode.p = policyEstimates[i][j]
			childNode.nextSibling = firstChildNode
			firstChildNode = &childNode
		}
	}
	if firstChildNode == nil {
		panic("There should be at least 1 legal move")
	}

	// Normalize policy
	for childNode := firstChildNode; childNode != nil; childNode = childNode.nextSibling {
		childNode.p /= totalLegalPolicy
	}

	node.firstChild = firstChildNode
}

// CalculateUctValue computes the priority of a node for exploration.
// Nodes with higher values should be explored first.
func CalculateUctValue(node *SearchNode, numParentVisits uint) float32 {
	cpuct := float32(1.2)
	return node.q +
		cpuct*node.p*float32(
			math.Sqrt(float64(numParentVisits)/
				float64(1.0+node.n)))
}

// DoVisit performs one iteration of tree search.
func DoVisit(tree *SearchTree) {
	// Select a leaf node to visit
	currentNode := tree.rootNode
	currentBoard := tree.board
	for currentNode.firstChild != nil {
		// While we're not at a leaf node:
		bestCandidateNode := (*SearchNode)(nil)
		bestUctValue := float32(math.Inf(-1))
		candidateNode := currentNode.firstChild
		for candidateNode != nil {
			uctValue := CalculateUctValue(candidateNode, uint(currentNode.n))
			if math.IsNaN(float64(uctValue)) {
				panic("UCT value should not be NaN")
			}
			if uctValue > bestUctValue {
				bestCandidateNode = candidateNode
				bestUctValue = uctValue
			}
			candidateNode = candidateNode.nextSibling
		}
		currentNode = bestCandidateNode
		currentBoard = PlayMove(
			currentBoard,
			// Previous player
			OtherPlayer(currentNode.player),
			bestCandidateNode.move.Row,
			bestCandidateNode.move.Col)
	}

	// Expand the selected leaf node
	winner := GetWinner(currentBoard)
	if winner != 0 {
		currentNode.isTerminal = true
		currentNode.v = 1
	} else {
		EvaluateAtNode(currentNode, currentBoard)
	}

	// Back up the evaluated value
	visitValue := currentNode.v
	nodeToUpdate := currentNode
	for nodeToUpdate != nil {
		nodeToUpdate.w += visitValue
		nodeToUpdate.n++
		nodeToUpdate.q = nodeToUpdate.w / float32(nodeToUpdate.n)

		if nodeToUpdate.parent == nil {
			break
		}
		nodeToUpdate = nodeToUpdate.parent
		// Flip value for opponent
		visitValue = -visitValue
	}
}

// GetBestMove gets the estimated best move at the root of a search tree
func GetBestMove(tree *SearchTree) Move {
	maxVisits := -1
	bestMove := (*Move)(nil)
	for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
		if int(childNode.n) > maxVisits {
			maxVisits = int(childNode.n)
			bestMove = &childNode.move
		}
	}
	return *bestMove
}

// GetMoveWithTemperatureOne picks a move using a "temperature" of 1.0
// The probability a move is picked is proportional to the number of visits.
func GetMoveWithTemperatureOne(tree *SearchTree) Move {
	totalVisits := 0
	for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
		totalVisits += int(childNode.n)
	}

	randInt := rand.Intn(totalVisits)
	cumulativeVisits := 0
	for childNode := tree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
		cumulativeVisits += int(childNode.n)
		if randInt <= cumulativeVisits {
			return childNode.move
		}
	}
	panic("Expected to pick a move")
}

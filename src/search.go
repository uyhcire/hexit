package hexit

import (
	"math"
	"math/rand"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"gonum.org/v1/gonum/stat/distuv"
)

// SearchNode is a node in a search tree
type SearchNode struct {
	// Move that led to this node
	move Move
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
	game     Game
	rootNode *SearchNode
}

// NewSearchNode creates a new SearchNode
func NewSearchNode(parent *SearchNode, move Move) SearchNode {
	nan := float32(math.NaN())
	return SearchNode{
		move:        move,
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
func NewSearchTree(evaluatePosition Evaluator, game Game) SearchTree {
	if GetWinner(game.Board) != 0 {
		panic("Can't search from a terminal node")
	}

	rootNode := NewSearchNode(nil, Move{Row: 1000, Col: 1000})
	searchTree := SearchTree{
		game:     game,
		rootNode: &rootNode,
	}
	EvaluateAtNode(evaluatePosition, searchTree.rootNode, game)
	return searchTree
}

// ApplyDirichletNoise applies noise to the root node's policy estimates
func ApplyDirichletNoise(newSearchTree *SearchTree) {
	epsilon := float32(0.25)
	alpha := 0.3
	gammaDistribution := distuv.Gamma{Alpha: alpha, Beta: 1.0}

	totalNoise := float32(0)
	noiseVector := make([]float32, 0)
	for childNode := newSearchTree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
		noise := float32(gammaDistribution.Rand())
		noiseVector = append(noiseVector, noise)
		totalNoise += noise
	}
	for i := range noiseVector {
		noiseVector[i] /= totalNoise
	}

	i := 0
	for childNode := newSearchTree.rootNode.firstChild; childNode != nil; childNode = childNode.nextSibling {
		childNode.p = (1-epsilon)*childNode.p + epsilon*noiseVector[i]
		i++
	}
}

type Evaluator = func(Board, byte) (float32, [5][5]float32)

// EvaluatePositionRandomly returns random value and policy estimates for a position.
func EvaluatePositionRandomly(board Board, player byte) (float32, [5][5]float32) {
	valueEstimate := rand.Float32()*2 - 1
	policyEstimates := [5][5]float32{}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			policyEstimates[i][j] = rand.Float32()
		}
	}
	return valueEstimate, policyEstimates
}

var model *tf.SavedModel

func InitializeModel() {
	if model == nil {
		savedModel, err := tf.LoadSavedModel("hexit_saved_model", []string{"serve"}, nil)
		if err != nil {
			panic(err)
		}
		model = savedModel
	}
}

func EvaluatePositionWithNN(board Board, player byte) (float32, [5][5]float32) {
	if model == nil {
		panic("Model not initialized")
	}

	squaresOccupiedByMyself, squaresOccupiedByOtherPlayer := GetOccupiedSquaresForNN(board, player)
	boardInput := [][]float32{
		append(squaresOccupiedByMyself, squaresOccupiedByOtherPlayer...),
	}
	boardInputTensor, err := tf.NewTensor(boardInput)
	if err != nil {
		panic(err)
	}

	boardInputOperation := model.Graph.Operation("boardInput")
	policyOutputOperation := model.Graph.Operation("policyOutput/Softmax")
	valueOutputOperation := model.Graph.Operation("valueOutput/Tanh")
	if boardInputOperation == nil {
		panic("boardInput operation not found")
	}
	if policyOutputOperation == nil {
		panic("policyOutput operation not found")
	}
	if valueOutputOperation == nil {
		panic("valueOutput operation not found")
	}
	result, err := model.Session.Run(
		map[tf.Output]*tf.Tensor{
			boardInputOperation.Output(0): boardInputTensor,
		},
		[]tf.Output{
			policyOutputOperation.Output(0),
			valueOutputOperation.Output(0),
		},
		nil,
	)
	if err != nil {
		panic(err)
	}

	policyOutputs := result[0].Value().([][]float32)
	valueOutputs := result[1].Value().([][]float32)

	valueEstimate := float32(0.0)
	if player == 1 {
		valueEstimate = valueOutputs[0][0]
	} else {
		valueEstimate = -valueOutputs[0][0]
	}

	policyEstimates := [5][5]float32{}
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if player == 1 {
				policyEstimates[i][j] = policyOutputs[0][i*5+j]
			} else {
				policyEstimates[j][i] = policyOutputs[0][i*5+j]
			}
		}
	}

	return valueEstimate, policyEstimates
}

// EvaluateAtNode evaluates the NN at a single node.
func EvaluateAtNode(evaluatePosition Evaluator, node *SearchNode, game Game) {
	if node.isTerminal {
		panic("Should not evaluate the NN at a terminal node")
	}

	valueEstimate, policyEstimates := evaluatePosition(game.Board, game.CurrentPlayer)
	node.v = valueEstimate

	firstChildNode := (*SearchNode)(nil)
	totalLegalPolicy := float32(0.0)
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if game.Board[i][j] != 0 {
				// Illegal move
				continue
			}
			totalLegalPolicy += policyEstimates[i][j]
			childNode := NewSearchNode(node, Move{Row: uint(i), Col: uint(j)})
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
func DoVisit(tree *SearchTree, evaluatePosition Evaluator) {
	// Select a leaf node to visit
	currentNode := tree.rootNode
	currentGame := tree.game
	var err error
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
		// Skip the side-switching move
		if currentGame.MoveNum == 2 {
			err, currentGame = DoNotSwitchSides(currentGame)
			if err != nil {
				panic(err)
			}
		}
		err, currentGame = PlayGameMove(currentGame, bestCandidateNode.move.Row, bestCandidateNode.move.Col)
		if err != nil {
			panic(err)
		}
	}

	// Expand the selected leaf node
	winner := GetWinner(currentGame.Board)
	if winner != 0 {
		currentNode.isTerminal = true
		currentNode.v = 1
	} else {
		EvaluateAtNode(evaluatePosition, currentNode, currentGame)
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

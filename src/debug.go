package hexit

import (
	fmt "fmt"
	"sort"
)

type By func(node1, node2 *SearchNode) bool

func (by By) Sort(nodes []*SearchNode) {
	sorter := &nodeSorter{
		nodes: nodes,
		by:    by,
	}
	sort.Sort(sorter)
}

type nodeSorter struct {
	nodes []*SearchNode
	by    func(node1, node2 *SearchNode) bool
}

func (s *nodeSorter) Len() int {
	return len(s.nodes)
}

func (s *nodeSorter) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

func (s *nodeSorter) Less(i, j int) bool {
	return s.by(s.nodes[i], s.nodes[j])
}

func PrintVisitDistribution(node *SearchNode) {
	childNodes := make([]*SearchNode, 0)
	for childNode := node.firstChild; childNode != nil; childNode = childNode.nextSibling {
		childNodes = append(childNodes, childNode)
	}

	descendingQValues := func(node1, node2 *SearchNode) bool {
		return node1.q > node2.q
	}
	By(descendingQValues).Sort(childNodes)

	for _, childNode := range childNodes {
		fmt.Printf(
			"(%d, %d) (N: %d) (Q: %2f) (U: %2f)\n",
			childNode.move.Row,
			childNode.move.Col,
			childNode.n,
			childNode.q,
			calculateUctU(childNode, uint(node.n)),
		)
	}
}

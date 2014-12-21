package main
//
// SIMPLE FAST ALGORITHMS FOR THE EDITING DISTANCE BETWEEN TREES AND RELATED PROBLEMS
// http://www.grantjenks.com/wiki/_media/ideas/simple_fast_algorithms_for_the_editing_distance_between_tree_and_related_problems.pdf
//

import (
	"log"
)

type Node struct {
	Id       int
	ParentId int
	Label    string
}

func (n *Node) IsRoot() bool {
	if n.ParentId == -1 {
		return true
	} else {
		return false
	}
}

type Tree struct {
	Root        Node
	Nodes       []Node
	ChildlenIds map[int][]int
	PostOrder   []int
}

func NewTree(parents []int, labels []string) Tree {
	var root Node
	nodes := make([]Node, 0, len(parents))
	childlenIds := make(map[int][]int, 0)
	for nid, pid := range parents {
		n := Node{Id: nid, ParentId: pid, Label: labels[nid]}
		nodes = append(nodes, n)
		childlenIds[pid] = append(childlenIds[pid], nid)
	}

	for nid, pid := range parents {
		if pid == -1 {
			root = nodes[nid]
			break
		}
	}
	t := Tree{Root: root, Nodes: nodes,
		ChildlenIds: childlenIds, PostOrder: make([]int, 0, len(nodes))}
	t.SetPostOrder(t.Root)
	return t
}

func (t *Tree) SetPostOrder(n Node) {
	// left-to-right postorder numbering
	for _, chid := range t.ChildlenIds[n.Id] {
		t.SetPostOrder(t.Nodes[chid])
	}
	t.PostOrder = append(t.PostOrder, n.Id)
}

func (t *Tree) LeftMostLeaf(node Node) Node {
	var childNode Node
	for {
		if len(t.ChildlenIds[node.Id]) == 0 {
			return node // left most leaf is this node
		}
		leftChildId := t.ChildlenIds[node.Id][0]
		childNode = t.Nodes[leftChildId]
		if len(t.ChildlenIds[childNode.Id]) == 0 {
			break
		}
		node = childNode
	}
	return childNode
}

func (t *Tree) LRKeyRoots() []Node {
	keynodes := make([]Node, 0, len(t.Nodes))
	for _, node := range t.Nodes {
		if node.ParentId == -1 {
			keynodes = append(keynodes, node)
		} else {
			p := t.Nodes[node.ParentId]
			if t.LeftMostLeaf(node).Id != t.LeftMostLeaf(p).Id {
				keynodes = append(keynodes, node)
			}
		}
	}
	return keynodes
}

func (t *Tree) Size() int {
	return len(t.Nodes)
}

type TreeEditDistance struct {
	T1       Tree
	T2       Tree
	TreeDist [][]float64
}

func NewTreeEditDistance(t1, t2 Tree) TreeEditDistance {
	td := make([][]float64, 0, t1.Size())
	for i := 0; i < t1.Size(); i++ {
		td = append(td, make([]float64, 0, t2.Size()))
		for j := 0; j < t2.Size(); j++ {
			td[i] = append(td[i], 0)
		}
	}
	ted := TreeEditDistance{T1: t1, T2: t2, TreeDist: td}
	return ted
}

func (ted *TreeEditDistance) Calc() float64 {
	for _, n1 := range ted.T1.LRKeyRoots() {
		for _, n2 := range ted.T2.LRKeyRoots() {
			ted.CalcTreeDist(n1, n2)
		}
	}
	log.Println("TREE DIST(FINAL)")
	for _, v := range ted.TreeDist {
		log.Println(v)
	}
	return ted.TreeDist[ted.T1.Size()-1][ted.T2.Size()-1]
}

func (ted *TreeEditDistance) CalcTreeDist(ni, nj Node) {
	m := ted.T1.PostOrder[ni.Id] - ted.T1.PostOrder[ted.T1.LeftMostLeaf(ni).Id] + 1
	n := ted.T2.PostOrder[nj.Id] - ted.T2.PostOrder[ted.T2.LeftMostLeaf(nj).Id] + 1
	//     log.Println("for l(i) to i", ted.T1.PostOrder[ted.T1.LeftMostLeaf(ni).Id], ted.T1.PostOrder[ni.Id])
	//     log.Println("for l(j) to j", ted.T2.PostOrder[ted.T2.LeftMostLeaf(nj).Id], ted.T2.PostOrder[nj.Id])
	list1 := make([]int, 0, m)
	for mm := ted.T1.PostOrder[ted.T1.LeftMostLeaf(ni).Id]; mm <= ted.T1.PostOrder[ni.Id]; mm++ {
		list1 = append(list1, mm)
	}
	list2 := make([]int, 0, n)
	for nn := ted.T2.PostOrder[ted.T2.LeftMostLeaf(nj).Id]; nn <= ted.T2.PostOrder[nj.Id]; nn++ {
		list2 = append(list2, nn)
	}

	forestdist := make([][]float64, 0, m)
	for i1 := 0; i1 <= m; i1++ {
		forestdist = append(forestdist, make([]float64, 0, n))
		for j1 := 0; j1 <= n; j1++ {
			forestdist[i1] = append(forestdist[i1], 0.)
		}
	}
	for i1 := 1; i1 <= m; i1++ {
		forestdist[i1][0] = forestdist[i1-1][0] + 1.
	}
	for j1 := 1; j1 <= n; j1++ {
		forestdist[0][j1] = forestdist[0][j1-1] + 1.
	}

	var repCost float64
	for i1id, i1 := range list1 {
		i1id += 1
		for j1id, j1 := range list2 {
			//             log.Println("i1", "j1", i1, j1, "m", "n", m, n)
			j1id += 1
			li1 := ted.T1.LeftMostLeaf(ted.T1.Nodes[i1]).Id
			li := ted.T1.LeftMostLeaf(ni).Id
			lj1 := ted.T2.LeftMostLeaf(ted.T2.Nodes[j1]).Id
			lj := ted.T2.LeftMostLeaf(nj).Id
			//             log.Println("l(i1), l(i), l(j1), l(j)", li1, li, lj1, lj)
			if li1 == li && lj1 == lj {
				x := forestdist[i1id-1][j1id] + 1.
				y := forestdist[i1id][j1id-1] + 1.
				if ted.T1.Nodes[ted.T1.PostOrder[i1]].Label == ted.T2.Nodes[ted.T2.PostOrder[j1]].Label {
					repCost = 0.
				} else {
					repCost = 1.
				}
				z := forestdist[i1id-1][j1id-1] + repCost
				forestdist[i1id][j1id] = Min(x, y, z)
				ted.TreeDist[i1][j1] = forestdist[i1id][j1id]
			} else {
				x := forestdist[i1id-1][j1id] + 1.
				y := forestdist[i1id][j1id-1] + 1.
				z := forestdist[i1id-1][j1id-1] + ted.TreeDist[i1][j1]
				forestdist[i1id][j1id] = Min(x, y, z)
			}
		}
	}
	log.Println("TREE DIST(i, j)", ni.Id+1, nj.Id+1)
	for _, v := range forestdist {
		log.Println(v)
	}
}

func Min(x, y, z float64) float64 {
	var min float64
	if x <= y && x <= z {
		min = x
	} else if y <= x && y <= z {
		min = y
	} else if z <= x && z <= y {
		min = z
	}
	return min
}

func main() {
	parents1 := []int{3, 2, 3, 5, 5, -1}
	labels1 := []string{"a", "b", "c", "d", "e", "f"}

	t1 := NewTree(parents1, labels1)
	//     log.Println("POSTORDER", t1.PostOrder)
	//     log.Println("LRKeyRoots", t1.LRKeyRoots()) // 3, 5, 6

	parents2 := []int{3, 3, 5, 2, 5, -1}
	labels2 := []string{"a", "b", "c", "d", "e", "f"}

	t2 := NewTree(parents2, labels2)
	//     log.Println("POSTORDER", t2.PostOrder)
	//     log.Println("LRKeyRoots", t2.LRKeyRoots()) // 2, 5, 6

	ted := NewTreeEditDistance(t1, t2)
	d := ted.Calc()
	log.Println("TED:", d)
}

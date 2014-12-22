package main

//
// SIMPLE FAST ALGORITHMS FOR THE EDITING DISTANCE BETWEEN TREES AND RELATED PROBLEMS
// http://www.grantjenks.com/wiki/_media/ideas/simple_fast_algorithms_for_the_editing_distance_between_tree_and_related_problems.pdf
//

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	Orig2Po     map[int]int
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
		ChildlenIds: childlenIds, PostOrder: make([]int, 0, len(nodes)),
		Orig2Po: make(map[int]int)}
	t.SetPostOrder(t.Root)
	return t
}

func (t *Tree) SetPostOrder(n Node) {
	// left-to-right postorder numbering
	for _, chid := range t.ChildlenIds[n.Id] {
		//         log.Println(t.Nodes[chid].Label)
		t.SetPostOrder(t.Nodes[chid])
	}
	t.Orig2Po[n.Id] = len(t.PostOrder)
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
	for _, nId := range t.PostOrder {
		node := t.Nodes[nId]
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
	SeqOpe   []string
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
	var (
		list1   []int
		list2   []int
		fd      [][]float64
		op      string
		repCost float64
	)

	for _, n1 := range ted.T1.LRKeyRoots() {
		for _, n2 := range ted.T2.LRKeyRoots() {
			//             log.Println(n1.Label, n2.Label)
			list1, list2, fd = ted.CalcTreeDist(n1, n2)
		}
	}
	//     log.Println("TREE DIST(FINAL)")
	//     for _, v := range ted.TreeDist {
	//         log.Println(v)
	//     }

	a := ted.T1.Size()
	b := ted.T2.Size()
	ni := ted.T1.Nodes[ted.T1.PostOrder[a-1]]
	nj := ted.T2.Nodes[ted.T2.PostOrder[b-1]]
	//     log.Println("TREEDIST")
	//     for _, v := range ted.TreeDist {
	//         log.Println(v)
	//     }
	//     log.Println("")
	//     log.Println("FORESTDIST")
	//     for _, v := range fd {
	//         log.Println(v)
	//     }
	//     log.Println("")
	//     os.Exit(1)

	var largerSize int
	if a > b {
		largerSize = a
	} else {
		largerSize = b
	}
	seqOpe := make([]string, 0, largerSize)
	for {
		//         log.Println("(a, b)=", a, b)
		nodeI1 := ted.T1.Nodes[ted.T1.PostOrder[a-1]]
		LnodeI1 := ted.T1.LeftMostLeaf(nodeI1)
		Lni := ted.T1.LeftMostLeaf(ni)

		nodeJ1 := ted.T2.Nodes[ted.T2.PostOrder[b-1]]
		LnodeJ1 := ted.T2.LeftMostLeaf(nodeJ1)
		Lnj := ted.T2.LeftMostLeaf(nj)

		li1 := ted.T1.PostOrder[LnodeI1.Id]
		li := ted.T1.PostOrder[Lni.Id]
		lj1 := ted.T2.PostOrder[LnodeJ1.Id]
		lj := ted.T2.PostOrder[Lnj.Id]
		if li1 == li && lj1 == lj {
			//             log.Println("USE ONLY TreeDist", li1 == li, lj1 == lj, a, b)
			x := fd[a-1][b] + 1.
			y := fd[a][b-1] + 1.
			if nodeI1.Label == nodeJ1.Label {
				repCost = 0.
			} else {
				repCost = 1.
			}
			z := fd[a-1][b-1] + repCost
			if x <= y && x <= z {
				op = "INSERT"
				a -= 1
			} else if y <= x && y <= z {
				op = "DELETE"
				b -= 1
			} else {
				op = "REPLACE"
				a -= 1
				b -= 1
			}
			switch op {
			case "INSERT":
				//                 log.Println(op, ted.T1.Nodes[ted.T1.PostOrder[a]])
				seqOpe = append(seqOpe, fmt.Sprintf("%s\t%s\tNULL", op, ted.T1.Nodes[ted.T1.PostOrder[a]].Label))
			case "DELETE":
				//                 log.Println(op, ted.T2.Nodes[ted.T2.PostOrder[b]])
				seqOpe = append(seqOpe, fmt.Sprintf("%s\tNULL\t%s", op, ted.T2.Nodes[ted.T2.PostOrder[b]].Label))
			case "REPLACE":
				if repCost > 0 {
					//                     log.Println(op, ted.T1.Nodes[ted.T1.PostOrder[a]], ted.T2.Nodes[ted.T2.PostOrder[b]])
					seqOpe = append(seqOpe, fmt.Sprintf("%s\t%s\t%s", op,
						ted.T1.Nodes[ted.T1.PostOrder[a]].Label, ted.T2.Nodes[ted.T2.PostOrder[b]].Label))
				}
			}
		} else {
			//             log.Println("USE TreeDist and forestdist", li1 == li, lj1 == lj, a, b)
			x := fd[a-1][b] + 1.
			y := fd[a][b-1] + 1.
			z := fd[a-1][b-1] + ted.TreeDist[list1[a-1]][list2[b-1]]
			if x <= y && x <= z {
				op = "INSERT"
				a -= 1
			} else if y <= x && y <= z {
				op = "DELETE"
				b -= 1
			} else {
				op = "REPLACE"
				a -= 1
				b -= 1
			}
			switch op {
			case "INSERT":
				//                 log.Println(op, ted.T1.Nodes[ted.T1.PostOrder[a]])
				//                 seqOpe = append(seqOpe, op)
				seqOpe = append(seqOpe, fmt.Sprintf("%s\t%s\tNULL", op, ted.T1.Nodes[ted.T1.PostOrder[a]].Label))
			case "DELETE":
				//                 log.Println(op, ted.T2.Nodes[ted.T2.PostOrder[b]])
				seqOpe = append(seqOpe, fmt.Sprintf("%s\tNULL\t%s", op, ted.T2.Nodes[ted.T2.PostOrder[b]].Label))
			case "REPLACE":
				if ted.TreeDist[list1[a]][list2[b]] > 0 {
					//                     log.Println(op, ted.T1.Nodes[ted.T1.PostOrder[a]], ted.T2.Nodes[ted.T2.PostOrder[b]])
					seqOpe = append(seqOpe, fmt.Sprintf("%s\t%s\t%s", op,
						ted.T1.Nodes[ted.T1.PostOrder[a]].Label, ted.T2.Nodes[ted.T2.PostOrder[b]].Label))
				}
			}
		}
		if a == 0 && b == 0 {
			break
		}
	}
	//     log.Println(seqOpe)
	ted.SeqOpe = seqOpe
	return ted.TreeDist[ted.T1.Size()-1][ted.T2.Size()-1]
}

func (ted *TreeEditDistance) CalcTreeDist(ni, nj Node) ([]int, []int, [][]float64) {
	m := ted.T1.Orig2Po[ni.Id] - ted.T1.Orig2Po[ted.T1.LeftMostLeaf(ni).Id] + 1
	n := ted.T2.Orig2Po[nj.Id] - ted.T2.Orig2Po[ted.T2.LeftMostLeaf(nj).Id] + 1
	list1 := make([]int, 0, m)
	for mm := ted.T1.Orig2Po[ted.T1.LeftMostLeaf(ni).Id]; mm <= ted.T1.Orig2Po[ni.Id]; mm++ {
		list1 = append(list1, mm) // for i1 = l(i) to i
	}
	list2 := make([]int, 0, n)
	for nn := ted.T2.Orig2Po[ted.T2.LeftMostLeaf(nj).Id]; nn <= ted.T2.Orig2Po[nj.Id]; nn++ {
		list2 = append(list2, nn) // for j1 = l(j) to j
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
	for i1id, i1 := range list1 { // i1 is position in postorder
		i1id += 1
		for j1id, j1 := range list2 { // j1 is position in postorder
			j1id += 1
			nodeI1 := ted.T1.Nodes[ted.T1.PostOrder[i1]]
			LnodeI1 := ted.T1.LeftMostLeaf(nodeI1)
			Lni := ted.T1.LeftMostLeaf(ni)

			nodeJ1 := ted.T2.Nodes[ted.T2.PostOrder[j1]]
			LnodeJ1 := ted.T2.LeftMostLeaf(nodeJ1)
			Lnj := ted.T2.LeftMostLeaf(nj)

			li1 := ted.T1.PostOrder[LnodeI1.Id]
			li := ted.T1.PostOrder[Lni.Id]
			lj1 := ted.T2.PostOrder[LnodeJ1.Id]
			lj := ted.T2.PostOrder[Lnj.Id]

			if li1 == li && lj1 == lj {
				x := forestdist[i1id-1][j1id] + 1.
				y := forestdist[i1id][j1id-1] + 1.
				if nodeI1.Label == nodeJ1.Label {
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
	return list1, list2, forestdist
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

func ReadCaboCha() ([][]int, [][]string) {
	scanner := bufio.NewScanner(os.Stdin)
	parents := make([][]int, 0, 100)
	labels := make([][]string, 0, 100)
	numSent := 0
	for i := 0; i < 2; i++ {
		parents = append(parents, make([]int, 0, 100))
		labels = append(labels, make([]string, 0, 100))
	}
	var (
		chunkText string
		modId     int
	)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "EOS" {
			labels[numSent] = append(labels[numSent], chunkText)
			parents[numSent] = append(parents[numSent], modId)
			numSent += 1
			if numSent > 2 {
				log.Fatal("The number of sentences must be 2.")
			}
			chunkText = ""
		} else if strings.HasPrefix(line, "* ") {
			if len(chunkText) > 0 {
				labels[numSent] = append(labels[numSent], chunkText)
				parents[numSent] = append(parents[numSent], modId)
			}
			id0 := strings.Index(line, " ")
			line = line[id0+1:]
			id1 := strings.Index(line, " ")
			//             headId, _ := strconv.Atoi(line[:id1])
			line = line[id1+1:]
			id2 := strings.Index(line, " ")
			modId, _ = strconv.Atoi(line[:id2-1])
			chunkText = ""
		} else {
			id0 := strings.Index(line, "\t") // end position of surface
			surface := line[0:id0]
			chunkText += surface
		}
	}
	return parents, labels
}

func main() {
	//     parents1 := []int{3, 2, 3, 5, 5, -1}
	//     labels1 := []string{"a", "b", "c", "d", "e", "f"}

	//     parents1 := []int{2, 2, -1}
	//     labels1 := []string{"太郎は", "花子に", "あげた"}

	//     parents1 := []int{1, -1, 3, 1, 1}
	//     labels1 := []string{"次郎は", "あげた", "泣いている", "三郎に", "飴を"}
	//     parents1 := []int{5, 2, 5, 4, 5, -1}
	//     labels1 := []string{"次郎は", "泣いている", "三郎に", "甘い", "飴を", "あげた"}

	//     t1 := NewTree(parents1, labels1)
	//     log.Println("T1", t1)
	//     log.Println("POSTORDER", t1.PostOrder)
	//     log.Println("LRKeyRoots", t1.LRKeyRoots()) // 3, 5, 6

	//     parents2 := []int{3, 3, 5, 2, 5, -1}
	//     labels2 := []string{"a", "b", "c", "d", "e", "f"}
	//     parents2 := []int{3, 3, 3, -1}
	//     labels2 := []string{"太郎は", "花子に", "プレゼントを", "あげた"}

	//     t2 := NewTree(parents2, labels2)
	//     log.Println("T2", t2)
	//     log.Println("POSTORDER", t2.PostOrder)
	//     log.Println("LRKeyRoots", t2.LRKeyRoots()) // 2, 5, 6

	//     ted := NewTreeEditDistance(t1, t2)
	//     d := ted.Calc()
	//     log.Println("TED:", d)
	parents, labels := ReadCaboCha()
	t1 := NewTree(parents[0], labels[0])
	t2 := NewTree(parents[1], labels[1])
	ted := NewTreeEditDistance(t1, t2)
	d := ted.Calc()
	fmt.Println(fmt.Sprintf("TED\t%f", d))
	for _, op := range ted.SeqOpe {
		fmt.Println(fmt.Sprintf("OPERATION\t%s", op))
	}
}

# Tree Edit Distance

## Usage

```
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
```


## Author
Takuya Makino

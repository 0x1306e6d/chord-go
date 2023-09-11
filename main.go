package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

var m = 6

type Node struct {
	Key         int
	FingerTable []*Node
	Successor   *Node
	Predecessor *Node
}

func (n Node) String() string {
	f := transform(n.FingerTable, func(n *Node) string { return fmt.Sprintf("%d", n.Key) })
	return fmt.Sprintf("{Key: %d, Finger: %s}", n.Key, strings.Join(f, ", "))
}

type nodeSort []Node

func (a nodeSort) Len() int {
	return len(a)
}

func (a nodeSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a nodeSort) Less(i, j int) bool {
	return a[i].Key < a[j].Key
}

// n ∈ (x, y)
func element(n, x, y int) bool {
	if x == y {
		panic("x == y")
	}
	if x < y {
		return n > x && n < y
	} else {
		return n > x || n < y
	}
}

func transform(nodes []*Node, f func(*Node) string) []string {
	var result []string
	for _, n := range nodes {
		result = append(result, f(n))
	}
	return result
}

func (n Node) findStart(id int) *Node {
	fmt.Println("findSuccessor", id, "N:", n.Key, "Successor:", n.Successor, "Predecessor:", n.Predecessor)
	if element(id, n.Key, n.Successor.Key+1) {
		return n.Successor
	}
	return n.Successor.findStart(id)
}

func pow(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func build(m int, nodes []Node) []Node {
	sort.Sort(nodeSort(nodes))

	nodes[0].Predecessor = &nodes[len(nodes)-1]
	nodes[len(nodes)-1].Successor = &nodes[0]
	for i := 1; i < len(nodes); i++ {
		nodes[i].Predecessor = &nodes[i-1]
		nodes[i-1].Successor = &nodes[i]
	}

	for i := range nodes {
		f := []*Node{}
		for k := 1; k <= m; k++ {
			s := (nodes[i].Key + pow(2, k-1)) % pow(2, m)
			succ := nodes[i].findStart(s)
			f = append(f, succ)
		}
		nodes[i].FingerTable = f
	}
	return nodes
}

func (n Node) FindSuccessor(id int) *Node {
	fmt.Println("FindSuccessor", id, "node:", n)
	if element(id, n.Key, n.Successor.Key+1) {
		return n.Successor
	} else {
		p := n.closestPrecedingFinger(id)
		return p.FindSuccessor(id)
	}
}

func (n Node) closestPrecedingFinger(id int) Node {
	fmt.Println("closestPrecedingFinger", id, "node:", n)
	for i := m - 1; i > 0; i-- {
		if element(n.FingerTable[i].Key, n.Key, id) {
			fmt.Println("return closestPrecedingFinger:", n.FingerTable[i])
			return *n.FingerTable[i]
		}
	}
	return n
}

func main() {
	rand.Seed(time.Now().Unix())

	nodes := build(m, []Node{
		{Key: 38},
		{Key: 21},
		{Key: 32},
		{Key: 8},
		{Key: 51},
		{Key: 14},
		{Key: 42},
		{Key: 48},
		{Key: 1},
		{Key: 56},
	})
	for _, n := range nodes {
		fmt.Println(n.Key, ":", n)
	}
	fmt.Println("=====================================")

	tests := []struct {
		key  int
		node int
	}{
		{10, 14},
		{24, 32},
		{30, 32},
		{38, 38},
		{54, 56},
		{0, 1},
		{56, 56},
	}
	for _, t := range tests {
		n := nodes[rand.Intn(len(nodes))]
		fmt.Println(n.Key, ":", n)
		succ := n.FindSuccessor(t.key)
		if succ.Key != t.node {
			panic(fmt.Sprintf("expected %d, got %d", t.node, succ.Key))
		} else {
			fmt.Println("Node", succ.Key, "has key", t.key)
		}
	}
}

package backend_utils

import (
	"fmt"
	"os"
)

type State struct {
	id        int
	epsilon   []*State          // ε-transitions
	trans     map[rune][]*State // normal transitions on a rune
	accepting bool
}

type NFA struct {
	start  *State
	accept *State
}

var stateID int // counter

func newState() *State {
	s := &State{
		id:      stateID,
		epsilon: []*State{},
		trans:   make(map[rune][]*State),
	}
	stateID++
	return s
}

func BuildNFA(node *RegexTreeNode) *NFA {
	if node == nil {
		return nil
	}
	start, accept := buildState(node)
	accept.accepting = true // last return state us accepting
	return &NFA{start: start, accept: accept}
}

func buildState(n *RegexTreeNode) (*State, *State) {
	switch n.operation {
	case "atom":
		s1 := newState()
		s2 := newState()
		s1.trans[n.value] = append(s1.trans[n.value], s2)
		return s1, s2

	case "charset":
		s1 := newState()
		s2 := newState()
		// Create transitions for all characters in the character set
		for _, char := range n.charSet {
			s1.trans[char] = append(s1.trans[char], s2)
		}
		return s1, s2

	case "concat":
		leftStart, leftAccept := buildState(n.left)
		rightStart, rightAccept := buildState(n.right)
		leftAccept.epsilon = append(leftAccept.epsilon, rightStart)
		return leftStart, rightAccept

	case "or":
		s := newState()
		e := newState()
		lStart, lAccept := buildState(n.left)
		rStart, rAccept := buildState(n.right)
		s.epsilon = append(s.epsilon, lStart, rStart)
		lAccept.epsilon = append(lAccept.epsilon, e)
		rAccept.epsilon = append(rAccept.epsilon, e)
		return s, e

	case "star":
		s := newState()
		e := newState()
		subStart, subAccept := buildState(n.left)
		s.epsilon = append(s.epsilon, subStart, e)
		subAccept.epsilon = append(subAccept.epsilon, subStart, e)
		return s, e

	case "plus":
		// one or more: X+ = X X*
		xStart, xAccept := buildState(n.left)
		loopStart, loopAccept := buildState(&RegexTreeNode{operation: "star", left: n.left})
		xAccept.epsilon = append(xAccept.epsilon, loopStart)
		return xStart, loopAccept

	case "optional":
		s := newState()
		e := newState()
		subStart, subAccept := buildState(n.left)
		s.epsilon = append(s.epsilon, subStart, e)
		subAccept.epsilon = append(subAccept.epsilon, e)
		return s, e
	}
	return nil, nil
}

//to DOT

func (nfa *NFA) ToDOT(filename string) error {
	visited := map[int]bool{}
	out := "digraph NFA {\n"
	out += "  rankdir=LR;\n"
	out += "  node [shape=circle];\n"
	out += fmt.Sprintf("  start -> %d;\n", nfa.start.id)
	out += writeStates(nfa.start, visited)
	out += fmt.Sprintf("  %d [shape=doublecircle];\n", nfa.accept.id)
	out += "}\n"

	return os.WriteFile(filename, []byte(out), 0644)
}

func writeStates(s *State, visited map[int]bool) string {
	if visited[s.id] {
		return ""
	}
	visited[s.id] = true
	str := ""
	for r, targets := range s.trans {
		for _, t := range targets {
			str += fmt.Sprintf("  %d -> %d [label=\"%c\"];\n", s.id, t.id, r)
			str += writeStates(t, visited)
		}
	}
	for _, t := range s.epsilon {
		str += fmt.Sprintf("  %d -> %d [label=\"ε\"];\n", s.id, t.id)
		str += writeStates(t, visited)
	}
	return str
}

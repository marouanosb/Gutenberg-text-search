package backend_utils

import (
	"fmt"
	"os"
	"sort"
)

// DFAState represents one DFA state, which is really a set of NFA states.
type DFAState struct {
	id     int
	nfaSet map[*State]struct{} // The NFA states this DFA state represents
	trans  map[rune]*DFAState  // Deterministic transitions
	final  bool                // Accepting if any NFA state is accepting
}

// DFA is the deterministic automaton.
type DFA struct {
	Start  *DFAState // Changed from start to Start (exported)
	states []*DFAState
}

var dfaID int

func newDFAState(set map[*State]struct{}, final bool) *DFAState {
	s := &DFAState{
		id:     dfaID,
		nfaSet: set,
		trans:  make(map[rune]*DFAState),
		final:  final,
	}
	dfaID++
	return s
}

// epsilonClosure returns all NFA states reachable from 'set' via ε-transitions.
func epsilonClosure(set map[*State]struct{}) map[*State]struct{} {
	closure := make(map[*State]struct{})
	stack := []*State{}
	for s := range set {
		closure[s] = struct{}{}
		stack = append(stack, s)
	}
	for len(stack) > 0 {
		s := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, e := range s.epsilon {
			if _, ok := closure[e]; !ok {
				closure[e] = struct{}{}
				stack = append(stack, e)
			}
		}
	}
	return closure
}

// move returns the set of NFA states reachable from 'set' on symbol r.
func move(set map[*State]struct{}, r rune) map[*State]struct{} {
	out := make(map[*State]struct{})
	for s := range set {
		for _, t := range s.trans[r] {
			out[t] = struct{}{}
		}
	}
	return out
}

// nfaToDFA determinizes the given NFA into a DFA.
func NFAToDFA(nfa *NFA) *DFA {
	dfaID = 0
	startSet := epsilonClosure(map[*State]struct{}{nfa.start: {}})
	startFinal := containsAccepting(startSet, nfa.accept)
	startDFA := newDFAState(startSet, startFinal)

	dfa := &DFA{Start: startDFA, states: []*DFAState{startDFA}}

	unmarked := []*DFAState{startDFA}
	seen := map[string]*DFAState{keyForSet(startSet): startDFA}

	for len(unmarked) > 0 {
		cur := unmarked[0]
		unmarked = unmarked[1:]

		// collect all symbols from NFA states in cur
		symbols := make(map[rune]struct{})
		for s := range cur.nfaSet {
			for r := range s.trans {
				symbols[r] = struct{}{}
			}
		}

		for r := range symbols {
			m := move(cur.nfaSet, r)
			if len(m) == 0 {
				continue
			}
			closure := epsilonClosure(m)
			k := keyForSet(closure)
			next, ok := seen[k]
			if !ok {
				next = newDFAState(closure, containsAccepting(closure, nfa.accept))
				seen[k] = next
				dfa.states = append(dfa.states, next)
				unmarked = append(unmarked, next)
			}
			cur.trans[r] = next
		}
	}

	return dfa
}

// helper: returns true if the NFA accept state is in set
func containsAccepting(set map[*State]struct{}, accept *State) bool {
	for s := range set {
		if s == accept {
			return true
		}
	}
	return false
}

// stable key for a set of *State by sorted ids
func keyForSet(set map[*State]struct{}) string {
	ids := []int{}
	for s := range set {
		ids = append(ids, s.id)
	}
	sort.Ints(ids)
	key := ""
	for i, id := range ids {
		if i > 0 {
			key += ","
		}
		key += fmt.Sprint(id)
	}
	return key
}

// ToDOT writes a Graphviz DOT file for the DFA.
func (d *DFA) ToDOT(filename string) error {
	out := "digraph DFA {\n"
	out += "  rankdir=LR;\n"
	out += "  node [shape=circle];\n"
	//out += fmt.Sprintf("  start -> %d;\n", d.start.id)

	for _, st := range d.states {
		shape := "circle"
		if st.final {
			shape = "doublecircle"
		}
		out += fmt.Sprintf("  %d [shape=%s];\n", st.id, shape)
		for r, t := range st.trans {
			out += fmt.Sprintf("  %d -> %d [label=\"%c\"];\n", st.id, t.id, r)
		}
	}
	out += "}\n"
	return os.WriteFile(filename, []byte(out), 0644)
}

// Accept checks whether the DFA accepts a string.
func (d *DFA) Accept(input string) bool {
	cur := d.Start
	for _, r := range input {
		nxt, ok := cur.trans[r]
		if !ok {
			return false
		}
		cur = nxt
	}
	return cur.final
}

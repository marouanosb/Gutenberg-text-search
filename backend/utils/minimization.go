package backend_utils

// Minimize returns a new DFA that is equivalent but with the minimal number of states.
func (d *DFA) Minimize() *DFA {
	// 1. collect alphabet
	alphabet := map[rune]struct{}{}
	for _, st := range d.states {
		for r := range st.trans {
			alphabet[r] = struct{}{}
		}
	}
	symbols := []rune{}
	for r := range alphabet {
		symbols = append(symbols, r)
	}

	// 2. initial partition: final vs non-final
	final, nonfinal := map[*DFAState]struct{}{}, map[*DFAState]struct{}{}
	for _, s := range d.states {
		if s.final {
			final[s] = struct{}{}
		} else {
			nonfinal[s] = struct{}{}
		}
	}
	P := []map[*DFAState]struct{}{}
	if len(final) > 0 {
		P = append(P, final)
	}
	if len(nonfinal) > 0 {
		P = append(P, nonfinal)
	}

	// Hopcroft refinement
	W := append([]map[*DFAState]struct{}{}, P...)
	for len(W) > 0 {
		A := W[len(W)-1]
		W = W[:len(W)-1]
		for _, c := range symbols {
			// X = states whose c-transition goes to A
			X := map[*DFAState]struct{}{}
			for _, s := range d.states {
				if t, ok := s.trans[c]; ok {
					if _, in := A[t]; in {
						X[s] = struct{}{}
					}
				}
			}
			newP := []map[*DFAState]struct{}{}
			for _, Y := range P {
				// split Y into Y∩X and Y\X
				iPart := map[*DFAState]struct{}{}
				dPart := map[*DFAState]struct{}{}
				for s := range Y {
					if _, in := X[s]; in {
						iPart[s] = struct{}{}
					} else {
						dPart[s] = struct{}{}
					}
				}
				if len(iPart) > 0 && len(dPart) > 0 {
					newP = append(newP, iPart, dPart)
					// keep sets to refine
					for idx, Z := range W {
						if sameSet(Z, Y) {
							W[idx] = iPart
							W = append(W, dPart)
							break
						}
					}
					if !containsSet(W, Y) {
						if len(iPart) < len(dPart) {
							W = append(W, iPart)
						} else {
							W = append(W, dPart)
						}
					}
				} else {
					newP = append(newP, Y)
				}
			}
			P = newP
		}
	}

	// 3. Build new DFA states for each partition
	partMap := map[*DFAState]int{}
	for i, part := range P {
		for s := range part {
			partMap[s] = i
		}
	}
	newStates := make([]*DFAState, len(P))
	for i, part := range P {
		// pick a rep to decide final
		final := false
		for s := range part {
			if s.final {
				final = true
				break
			}
		}
		newStates[i] = &DFAState{id: i, trans: make(map[rune]*DFAState), final: final}
	}
	for i, part := range P {
		rep := first(part)
		for r, t := range rep.trans {
			newStates[i].trans[r] = newStates[partMap[t]]
		}
	}

	return &DFA{Start: newStates[partMap[d.Start]], states: newStates}
}

func sameSet(a, b map[*DFAState]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for s := range a {
		if _, ok := b[s]; !ok {
			return false
		}
	}
	return true
}
func containsSet(list []map[*DFAState]struct{}, target map[*DFAState]struct{}) bool {
	for _, s := range list {
		if sameSet(s, target) {
			return true
		}
	}
	return false
}
func first(m map[*DFAState]struct{}) *DFAState {
	for s := range m {
		return s
	}
	return nil
}

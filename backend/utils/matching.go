package utils

import (
	"bufio"
)

// MatchInText returns true if the DFA accepts any substring of `text`.
func matchInText(dfaStart *DFAState, text string) bool {
	runes := []rune(text)

	for i := 0; i < len(runes); i++ { // start position
		state := dfaStart
		for j := i; j < len(runes); j++ { // extend the substring
			next, ok := state.trans[runes[j]]
			if !ok {
				break // no transition, stop this substring
			}
			state = next
			if state.final {
				return true // found a substring that matches
			}
		}
	}
	return false
}

func MatchAllText(DFAStart *DFAState, scanner *bufio.Scanner) (matched bool, number_matches int, matches map[int]string) {
	number_matches = 0
	matches = map[int]string{}
	line_number := 0

	for scanner.Scan() {
		line_number++
		line := scanner.Text()

		if matchInText(DFAStart, line) {
			matches[line_number] = line
			number_matches++
		}
	}
	return number_matches > 0, number_matches, matches
}

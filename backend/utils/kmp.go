package utils

import "bufio"

func CreateCarryOverTable(f string) []int {
	n := len(f)
	co := make([]int, n+1)

	// regles initailes
	co[0] = -1
	co[1] = 0
	co[n] = 0

	// regle 1
	for i := 2; i < n; i++ {
		lps := detectLongestPrefixSuffix(f[:i]) // :i] not inclusive
		co[i] = lps
	}

	// regle 2 + 3
	for i := 2; i < n-1; i++ {
		if f[i] == f[co[i]] {
			if co[co[i]] == -1 { // regele 2
				co[i] = -1
			} else { // regele 3
				co[i] = co[co[i]]
			}
		}
	}

	return co
}

func detectLongestPrefixSuffix(pattern string) int {
	n := len(pattern)
	if n == 0 {
		return 0
	}
	pi := make([]int, n)
	pi[0] = 0
	for i := 1; i < n; i++ {
		j := pi[i-1]
		for j > 0 && pattern[i] != pattern[j] {
			j = pi[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		pi[i] = j
	}
	return pi[n-1]
}

func kmpSearchSingleLine(pattern string, text string, co []int) bool {
	runePattern := []rune(pattern)
	runeText := []rune(text)

	i := 0 // index for text
	j := 0 // index for pattern

	for i < len(runeText) {
		if runeText[i] == runePattern[j] {
			i++
			j++
			if j == len(runePattern) {
				return true // found a match
			}
		} else {
			if co[j] == -1 { // -1 on avance dans le texte mais on ressaye avec le pattern de 0
				i++
				j = 0
			} else {
				j = co[j]
			}
		}
	}
	return false
}

func KMPSearch(pattern string, scanner *bufio.Scanner, co []int) (matched bool, number_matches int, matches map[int]string) {
	number_matches = 0
	matches = map[int]string{}
	line_number := 0

	for scanner.Scan() {
		line_number++
		line := scanner.Text()
		if kmpSearchSingleLine(pattern, line, co) {
			matches[line_number] = line
			number_matches++
		}
	}
	return number_matches > 0, number_matches, matches
}

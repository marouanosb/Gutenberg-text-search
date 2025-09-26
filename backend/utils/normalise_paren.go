package backend_utils

func AddParentheses(pattern string) string {
	result := pattern
	i := 0

	for i < len(pattern) {
		c := pattern[i]
		if c == '*' || c == '+' || c == '?' {
			if i > 0 {
				// Check what comes before the quantifier
				if pattern[i-1] == ')' {
					// Already has parentheses, skip
					i++
					continue
				} else if pattern[i-1] == ']' {
					// Character class, find the opening bracket
					start := findMatchingOpenBracket(pattern, i-1)
					if start != -1 {
						// Wrap the entire character class
						result = result[:start] + "(" + result[start:i] + ")" + result[i:]
						i += 2 // account for added parentheses
						continue
					}
				} else {
					// Single character, wrap it
					result = result[:i-1] + "(" + string(pattern[i-1]) + ")" + result[i:]
					i += 2 // account for added parentheses
					continue
				}
			}
		}
		i++
	}
	return result
}

// Find the matching opening bracket for a closing bracket at position pos
func findMatchingOpenBracket(pattern string, pos int) int {
	if pos >= len(pattern) || pattern[pos] != ']' {
		return -1
	}

	for i := pos - 1; i >= 0; i-- {
		if pattern[i] == '[' {
			return i
		}
	}
	return -1
}

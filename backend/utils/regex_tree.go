package backend_utils

// regex tree node structure
type RegexTreeNode struct {
	operation string
	value     rune
	charSet   []rune // for character class [a-z], [abc], etc.
	left      *RegexTreeNode
	right     *RegexTreeNode
}

// check if the node is an atom
func (n *RegexTreeNode) isAtom() bool {
	return n.left == nil && n.right == nil
}

// calculate the size of atoms or parentheses with quantifier
func atomSize(pattern string) int {
	if len(pattern) == 0 {
		return 0
	}

	// Handle character class [...]
	if pattern[0] == '[' {
		end := matchingBracket(pattern)
		if end == -1 {
			// unbalanced brackets, treat as single char
			end = 0
		}
		end++ // include the closing ']'
		// check for trailing quantifier
		if end < len(pattern) && (pattern[end] == '*' || pattern[end] == '+' || pattern[end] == '?') {
			end++
		}
		return end
	}

	if pattern[0] == '(' {
		end := matchingParen(pattern)
		if end == -1 {
			// unbalanced parentheses, treat as single char
			end = 0
		}
		end++ // include the closing ')'
		// check for trailing quantifier
		if end < len(pattern) && (pattern[end] == '*' || pattern[end] == '+' || pattern[end] == '?') {
			end++
		}
		return end
	}

	// single character atom
	size := 1
	if size < len(pattern) && (pattern[size] == '*' || pattern[size] == '+' || pattern[size] == '?') {
		size++
	}
	return size
}

// return index of closing matching parenthese
func matchingParen(s string) int {
	depth := 0
	for i, c := range s {
		if c == '(' {
			depth++
		} else if c == ')' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// return index of closing matching bracket
func matchingBracket(s string) int {
	if len(s) == 0 || s[0] != '[' {
		return -1
	}

	for i := 1; i < len(s); i++ {
		if s[i] == ']' {
			return i
		}
		// Handle escaped characters inside brackets
		if s[i] == '\\' && i+1 < len(s) {
			i++ // skip next character
		}
	}
	return -1
}

// parse character class content like "a-z", "abc", "0-9A-F"
func parseCharacterClass(content string) []rune {
	var charSet []rune
	i := 0

	for i < len(content) {
		if i+2 < len(content) && content[i+1] == '-' {
			// Handle range like a-z, 0-9
			start := rune(content[i])
			end := rune(content[i+2])
			for r := start; r <= end; r++ {
				charSet = append(charSet, r)
			}
			i += 3
		} else {
			// Handle single character
			charSet = append(charSet, rune(content[i]))
			i++
		}
	}

	return charSet
}

// parse regex pattern into a tree
func (n *RegexTreeNode) ParseRegex(pattern string) *RegexTreeNode {
	if len(pattern) == 0 {
		return nil
	}
	if len(pattern) == 1 {
		if pattern[0] != '(' && pattern[0] != ')' && pattern[0] != '[' && pattern[0] != ']' {
			return &RegexTreeNode{operation: "atom", value: rune(pattern[0])}
		}
		return nil
	}

	// First, handle | operators at top level (lowest precedence)
	depth := 0
	bracketDepth := 0
	for i, char := range pattern {
		switch char {
		case '(':
			depth++
		case ')':
			depth--
		case '[':
			bracketDepth++
		case ']':
			bracketDepth--
		case '|':
			if depth == 0 && bracketDepth == 0 { // found a | at top level
				return &RegexTreeNode{
					operation: "or",
					left:      n.ParseRegex(pattern[:i]),
					right:     n.ParseRegex(pattern[i+1:]),
				}
			}
		}
	}

	// Handle outer parentheses
	if pattern[0] == '(' && matchingParen(pattern) == len(pattern)-1 {
		return n.ParseRegex(pattern[1 : len(pattern)-1])
	}

	// Handle quantifiers (* + ?) at the end
	last := pattern[len(pattern)-1]
	if last == '*' || last == '+' || last == '?' {
		op := map[byte]string{'*': "star", '+': "plus", '?': "optional"}[last]
		return &RegexTreeNode{
			operation: op,
			left:      n.ParseRegex(pattern[:len(pattern)-1]),
		}
	}

	// Handle concatenation by splitting at the first atom
	split := atomSize(pattern)
	if split >= len(pattern) {
		// Handle single character class or atom
		if pattern[0] == '[' {
			end := matchingBracket(pattern)
			if end != -1 && end == len(pattern)-1 {
				charSet := parseCharacterClass(pattern[1:end])
				return &RegexTreeNode{operation: "charset", charSet: charSet}
			}
		}
		// Single character
		return &RegexTreeNode{operation: "atom", value: rune(pattern[0])}
	}

	return &RegexTreeNode{
		operation: "concat",
		left:      n.ParseRegex(pattern[:split]),
		right:     n.ParseRegex(pattern[split:]),
	}
}

// print the tree
func (n *RegexTreeNode) PrintTree() {
	if n == nil {
		print("nil")
		return
	}
	if n.isAtom() {
		if n.operation == "charset" {
			print("[")
			for _, r := range n.charSet {
				print(string(r))
			}
			print("]")
		} else {
			print(string(n.value))
		}
		return
	}
	print(n.operation + " ( ")
	if n.left != nil {
		n.left.PrintTree()
	} else {
		print("nil")
	}
	if n.right != nil {
		print(" , ")
		n.right.PrintTree()
		print(" ) ")
	} else {
		print(" ) ")
	}
}

package utils

// regex tree node structure
type RegexTreeNode struct {
	operation string
	value     rune
	charSet   []rune
	left      *RegexTreeNode
	right     *RegexTreeNode
}

func (n *RegexTreeNode) isAtom() bool {
	return n.left == nil && n.right == nil
}

func atomSize(pattern string) int {
	if len(pattern) == 0 {
		return 0
	}
	if pattern[0] == '[' {
		end := matchingBracket(pattern)
		if end == -1 {
			end = 0
		}
		end++
		if end < len(pattern) && (pattern[end] == '*' || pattern[end] == '+' || pattern[end] == '?') {
			end++
		}
		return end
	}
	if pattern[0] == '(' {
		end := matchingParen(pattern)
		if end == -1 {
			end = 0
		}
		end++
		if end < len(pattern) && (pattern[end] == '*' || pattern[end] == '+' || pattern[end] == '?') {
			end++
		}
		return end
	}
	size := 1
	if size < len(pattern) && (pattern[size] == '*' || pattern[size] == '+' || pattern[size] == '?') {
		size++
	}
	return size
}

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

func matchingBracket(s string) int {
	if len(s) == 0 || s[0] != '[' {
		return -1
	}
	for i := 1; i < len(s); i++ {
		if s[i] == ']' {
			return i
		}
		if s[i] == '\\' && i+1 < len(s) {
			i++
		}
	}
	return -1
}

func parseCharacterClass(content string) []rune {
	var charSet []rune
	i := 0
	for i < len(content) {
		if i+2 < len(content) && content[i+1] == '-' {
			start := rune(content[i])
			end := rune(content[i+2])
			for r := start; r <= end; r++ {
				charSet = append(charSet, r)
			}
			i += 3
		} else {
			charSet = append(charSet, rune(content[i]))
			i++
		}
	}
	return charSet
}

func splitAtLastAtom(s string) (string, string) {
	if len(s) == 0 {
		return "", ""
	}
	i := 0
	lastStart := 0
	lastSize := 0
	for i < len(s) {
		size := atomSize(s[i:])
		if size <= 0 {
			lastStart = i
			lastSize = len(s) - i
			break
		}
		lastStart = i
		lastSize = size
		i += size
	}
	if lastSize == 0 {
		return s, ""
	}
	prefix := s[:lastStart]
	lastAtom := s[lastStart : lastStart+lastSize]
	return prefix, lastAtom
}

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
			if depth == 0 && bracketDepth == 0 {
				return &RegexTreeNode{
					operation: "or",
					left:      n.ParseRegex(pattern[:i]),
					right:     n.ParseRegex(pattern[i+1:]),
				}
			}
		}
	}
	if pattern[0] == '(' && matchingParen(pattern) == len(pattern)-1 {
		return n.ParseRegex(pattern[1 : len(pattern)-1])
	}
	last := pattern[len(pattern)-1]
	if last == '*' || last == '+' || last == '?' {
		op := map[byte]string{'*': "star", '+': "plus", '?': "optional"}[last]
		inner := pattern[:len(pattern)-1]
		prefix, atom := splitAtLastAtom(inner)
		if atom == "" {
			return &RegexTreeNode{
				operation: op,
				left:      n.ParseRegex(inner),
			}
		}
		quantNode := &RegexTreeNode{
			operation: op,
			left:      n.ParseRegex(atom),
		}
		if prefix == "" {
			return quantNode
		}
		return &RegexTreeNode{
			operation: "concat",
			left:      n.ParseRegex(prefix),
			right:     quantNode,
		}
	}
	split := atomSize(pattern)
	if split >= len(pattern) {
		if pattern[0] == '[' {
			end := matchingBracket(pattern)
			if end != -1 && end == len(pattern)-1 {
				charSet := parseCharacterClass(pattern[1:end])
				return &RegexTreeNode{operation: "charset", charSet: charSet}
			}
		}
		return &RegexTreeNode{operation: "atom", value: rune(pattern[0])}
	}
	return &RegexTreeNode{
		operation: "concat",
		left:      n.ParseRegex(pattern[:split]),
		right:     n.ParseRegex(pattern[split:]),
	}
}

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

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	utils "backend_utils"
)

func loadTitles() map[int]string {
	books := make(map[int]string)

	data, err := os.ReadFile("./resources/books.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &books)
	if err != nil {
		log.Fatal(err)
	}

	return books
}

func main() {
	// Load args
	pattern := os.Args[1]
	file_path := os.Args[2]

	// Read file arg
	file, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	println("-----")
	println("Pattern :", pattern)

	// Fixing parentheses
	pattern_par := utils.AddParentheses(pattern)
	if pattern != pattern_par {
		println("Pattern after fixing parentheses : ", pattern_par)
	}

	// Regex Tree
	tree := &utils.RegexTreeNode{}
	tree = tree.ParseRegex(pattern)
	print("Regex tree : ")
	tree.PrintTree()
	println("")
	println("-----")

	// NDFA
	nfa := utils.BuildNFA(tree)
	err = nfa.ToDOT("../outputs/nfa.dot")
	if err != nil {
		panic(err)
	}

	// DFA
	dfa := utils.NFAToDFA(nfa)
	dfa_min := dfa.Minimize()
	err = dfa_min.ToDOT("../outputs/min_dfa.dot")
	if err != nil {
		panic(err)
	}

	// Matching
	scanner := bufio.NewScanner(file)
	matched, number_matches, matches := utils.MatchAllText(dfa_min.Start, scanner)
	if matched {
		println("Matches found :", number_matches)
		matches_showed := 0

		for line, match := range matches {
			matches_showed++
			println("#", line, ":", strings.TrimSpace(match))
			// Show max 10 matches
			if matches_showed >= 10 {
				fmt.Printf("... %v more matches", number_matches-matches_showed)
				break
			}
		}
	} else {
		println("No matches found.")
	}

	// Serve
	//e := echo.New()
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(200, "Hello, World!")
	//})
	//e.Logger.Fatal(e.Start(":9111"))
}

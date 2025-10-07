package main

import (
	"backend_main/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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

type application struct {
	pattern string
	file    string
	algo    string
}

func main() {
	// Load args
	var app application
	if len(os.Args) == 4 { // algo mentioned
		app = application{
			pattern: os.Args[1],
			file:    os.Args[2],
			algo:    os.Args[3],
		}
	}
	if len(os.Args) == 3 { // algo not mentioned
		pattern := os.Args[1]
		regex_special_chars := []rune{'|', '.', '*', '+', '?', '(', ')', '[', ']', '\\', '^', '$'}
		is_regex := false
		for _, c := range pattern {
			if strings.ContainsRune(string(regex_special_chars), c) {
				is_regex = true
				break
			}
		}

		var algo string
		if is_regex {
			algo = "regex"
		} else {
			algo = "kmp"
		}

		app = application{
			pattern: pattern,
			file:    os.Args[2],
			algo:    algo,
		}
	}
	if len(os.Args) < 3 {
		app = application{
			pattern: "exampSargonle",
			file:    "../resources/livre_sur_babylone.txt",
			algo:    "kmp",
		}
	}

	// Read file arg
	file, err := os.Open(app.file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	println("-----")
	println("Pattern :", app.pattern)
	println("-----")
	fmt.Printf("Used < %s >  algo.\n", app.algo)

	// Choosing algo
	if app.algo == "regex" {
		// Regex
		// Fixing parentheses
		pattern_par := utils.AddParentheses(app.pattern)
		if app.pattern != pattern_par {
			println("Pattern after fixing parentheses : ", pattern_par)
		}

		// Regex Tree
		tree := &utils.RegexTreeNode{}
		tree = tree.ParseRegex(app.pattern)
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
		dfa_min := dfa.Minimize() // minimisation
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
	} else if app.algo == "kmp" {
		// KMP
		co := utils.CreateCarryOnTable(app.pattern)

		// Matching
		scanner := bufio.NewScanner(file)
		matched, number_matches, matches := utils.KMPSearch(app.pattern, scanner, co)
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
	}

	// Serve
	//e := echo.New()
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(200, "Hello, World!")
	//})
	//e.Logger.Fatal(e.Start(":9111"))
}

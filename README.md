# Moteur de recherche sur le texte des livres (Gutenberg)

Moteur de recherche textuelle sur les livres mis à la disposition par [The Gutenberg Project](http://gutenberg.org/). En utilisant la recherche de motifs textuels par RegEx avec la méthode Aho-Ullman.

## 1. How to run
    
### 1.1. Pre-requisites
- Installing [Go](https://go.dev/).

### 1.2. Running
- Navigate to `/backend` folder.
```shell
cd backend
```
- Run the main `server.go` file, with the expected arguments.
```shell
go run . <REGEX_PATTERN> <TEXT_FILE_PATH>
```
For example :
```shell
go run . "S((a|r|g)*)on" "../resources/livre_sur_babylone.txt"
```
Output :
```
-----
Pattern : S((a|r|g)*)on
Regex tree : concat ( S , concat ( star ( or ( a , or ( r , g )  )  )  , concat ( o , n )  )  ) 
-----
Matches found : 30
# 10930 : Sargon and the Assyrian army before its walls. Merodach-baladan was
# 432 : state--Sargon and Merodach-baladan--Sennacherib's attempt
# 1833 : to the Ishtar Gate, precisely the two points mentioned in Sargon's
# 10937 : After the defeat of Shabaka and the Egyptians at Raphia, Sargon was
# 10941 : their appearance from the north and east. In fact, Sargon's conquest of
# 1844 : A: Sargon's quay-wall. B: Older moat-wall. O: Later moat-wall of
# 1853 : quay-walls, which succeeded that of Sargon. The three narrow walls
# 1870 : in view of Sargon's earlier reference.
# 10920 : Sargon's army had secured the capture of Samaria, he was obliged to
# 1016 : to Sargon of Akkad; but that marked the extreme limit of Babylonian
# 1827 : upon it."[45] The two walls of Sargon, which he here definitely names
# 1914 : excavations. The discovery of Sargon's inscriptions proved that in
# 10954 : On Sargon's death in 705 B.C. the subject provinces of the empire
# 12453 : fifteen hundred years before the birth of Sargon I., who is supposed
# 1832 : the quay of Sargon,[46] which run from the old bank of the Euphrates
# 1832 : the quay of Sargon,[46] which run from the old bank of the Euphrates
... 20 more matches
```
Also, **.DOT** files corresponding to the NFA and DFA automatons are generated and can be found in `/outputs` folder. Which can be visualised using [Graphviz Online](https://dreampuf.github.io/GraphvizOnline).
Here's the example's DFA (note that the final DFA is minimized) :

!["S((a|r|g)*)on" regex pattern DFA](/resources/example_dfa.png)

## 2. Codebase

### 2.1. Backend
The backend file structure is as follows :
```
backend/
  ├─ utils/
  │   ├─ dfa_automat.go
  │   ├─ extract_books.py
  │   ├─ matching.go
  │   ├─ minimization.go
  │   ├─ ndfa_automat.go
  │   ├─ normalise_paren.go
  │   └─ regex_tree.go
  └─ server.go
```
#### Workflow
- `server.go`  
The main running file of the project, which regroups all the steps of the process.
    - Reads the given command-line arguments.
    - Normalizes the regex pattern's parenthesis to follow UNIX standard.
    - Generates the pattern's regex tree.
    - Generates the NFA from the given tree.
    - Generates the DFA from the given NFA.
    - Minimizes the DFA.
    - Reads the given line line by line and checks for matching patterns.

### 2.2. Frontend

#### Workflow

### 2.3. Outputs

### 2.4. Resources
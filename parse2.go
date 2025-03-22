package main

import (
	"fmt"
	"regexp"
	"strings"
	"os"
	"bufio"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func NoTrailingNewline(s string) string {
	if len(s) == 0 {
		return s
	}

	if s[len(s) - 1] == '\n' {
		s = s[:len(s) - 1]
	}

	return s
}

// Snippet is made of a list of source locations, which are
// either plain strings or a pasteme command. A non-inlined
// command appears on its own line
type SourceLocation string

// SourceLocation always just emits its string contents
func (s SourceLocation) Emit(t *SnippetTable) string {
	return NoTrailingNewline(string(s))
}

func (s SourceLocation) IsInlined() bool {
	return false
}

// A pasteme command is an empty string unless matched with
// a corresponding snippet one or more times
type PastemeCommand struct {
	name string
	inlined bool
	multiple bool
	zero bool
}

func NewPastemeEmitter(indentation int, name string, multiple bool) SourceLine {
	return SourceLine{
		indentation,
		[]Emitter{Emitter(PastemeCommand{name, false, multiple, false})},
	}
}

func (p PastemeCommand) IsInlined() bool {
	return p.inlined
}

// PastemeCommand emits a snippet if it matches one
func (p PastemeCommand) Emit(t *SnippetTable) string {
	i := -1
	var match *Snippet
	str := ""

	for {
		i, match = t.FindSnippet(p.name, i + 1)

		if match == nil {
			break;
		}

		// Snippet Match
		str += match.Emit(t)

		if !p.multiple {
			break
		}
	}

	return NoTrailingNewline(str)
}

// An emitter must return a string that does not end with a newline
type Emitter interface {
	Emit(*SnippetTable) string
	IsInlined() bool
}

// A source line remembers its indentation level and all locations
// it contains. Locations are separated by spaces.
type SourceLine struct {
	indentation int
	locations []Emitter
}

func NewSourceLine(indentation int, source string) SourceLine {
	return SourceLine{
		indentation, 
		[]Emitter{Emitter(SourceLocation(source))},
	}
}

func (s SourceLine) Emit(t *SnippetTable) string {
	str := ""

	// Concatenate locations
	for j, loc := range s.locations {
		// Check if next location is inlined
		if j != 0 && !loc.IsInlined() && !s.locations[j-1].IsInlined() {
			str += " "
		}

		str += loc.Emit(t)
	}

	// Split
	q, _ := regexp.Compile("\\n")
	lines := q.Split(str, -1)

	// Apply indentation to each line
	i := s.indentation
	indent := ""

	for i > 0 {
		if i >= 4 {
			indent += "\t"
			i -= 4
		} else {
			indent += " "
			i -= 1
		}
	}

	str = ""

	for _, line := range lines {
		str += indent + line + "\n"
	}

	return NoTrailingNewline(str)
}

// A snippet is a list of lines. A parsed file is also a snippet.
// name is used to match valid paste locations
type Snippet struct {
	name string
	lines []SourceLine
}

func (s Snippet) Emit(t *SnippetTable) string {
	str := ""

	for _, line := range s.lines {
		nextline := line.Emit(t)

		// Don't emit superfluous newlines after non-matched pastemes
		if len(strings.TrimSpace(nextline)) != 0 {
			str += nextline
			str += "\n"
		} else {
			// Don't emit extra whitespace on blank lines
			str += "\n"
		}
	}

	return str
}

// Dummy method for printing with empty snippet table
func (s Snippet) String() string {
	var t SnippetTable

	return s.Emit(&t)
}

// List of snippets used for executing pasteme commands
// Nextscope points to the next symbol table to check if a symbol
// is not contained in this table. Used for scoping
type SnippetTable struct {
	snippets []Snippet
	nextscope *SnippetTable
}

func (t *SnippetTable) FindSnippet(name string, startindex int) (int, *Snippet) {
	i := startindex

	for i < len(t.snippets) {
		snip := t.snippets[i]
		if snip.name == name {
			return i, &snip
		}

		i += 1
	}

	return 0, nil
}

func (t *SnippetTable) FirstSnippet(name string) *Snippet {
	_, snip := t.FindSnippet(name, 0)

	return snip
}

func (t *SnippetTable) AddSnippet(snip Snippet) {
	t.snippets = append(t.snippets, snip)
}

func (t *SnippetTable) String() string {
	str := ""

	for _, snip := range t.snippets {
		str += fmt.Sprint(" ", snip.name, " ::\n")
		str += snip.String()
	}

	return str + "\n"
}

func (t *SnippetTable) SnippetNames() []string {
	r := make([]string, len(t.snippets))

	for i, snip := range t.snippets {
		r[i] = snip.name
	}

	return r
}

func ParsePasteme(command string) *PastemeCommand {
	if len(command) == 0 {
		return nil
	}

	// Check for inlining modifier
	inlined := command[len(command) - 1] == '\\'

	if inlined {
		command = command[:len(command) - 1]
	}

	// Check for repetition modifier
	rep_mod := command[len(command) - 1]
	rep_zero := false
	rep_multiple := false

	switch {
		case rep_mod == '?' || rep_mod == '*': rep_zero = true
		case rep_mod == '+' || rep_mod == '*': rep_multiple = true
	}

	if rep_zero || rep_multiple {
		command = command[:len(command) - 1]
	}

	q, _ := regexp.Compile("\\s+")
	tokens := q.Split(command, -1)

	c := PastemeCommand{
		name: tokens[1],
		inlined: inlined,
		multiple: rep_multiple,
		zero: rep_zero,
	}

	return &c
}

// Separate indentation from actual text
func SeparateIndentation(line string) (int, string) {
	r, _ := regexp.Compile("^\\s*")
	match := r.FindString(line)

	indentation := 0

	for _, char := range match {
		switch {
			case char == '\t': indentation += 4
			case char == ' ': indentation += 1
			default: break
		} 
	}

	return indentation, line[len(match):]
}

// Return list of command tokens, if this string is a command
func GetCommandTokens(command string) []string {
	r, _ := regexp.Compile("^--\\s*")
	q, _ := regexp.Compile("\\s+")
	match := r.FindStringIndex(command)

	if len(match) > 0 {
		return q.Split(command[match[1]:], -1)
	}

	return nil
}

// Parse a snippet, returning any internal snippets in a SnippetTable
func ParseSnippet(name string, subsnippets *SnippetTable, scanner *bufio.Scanner) {
	var snippet Snippet
	snippet.name = name
	prev_inlined := false

	for scanner.Scan() {
		// Read line
		line := scanner.Text()
		indentation, source := SeparateIndentation(line)
		
		tokens := GetCommandTokens(source)

		var sline SourceLine

		if len(tokens) == 0 {
			// Source line, not command
			if len(source) != 0 || !prev_inlined {
				if prev_inlined && len(snippet.lines) > 0 {
					// Add to current snippet
					prevline := &snippet.lines[len(snippet.lines) - 1]
					prevline.locations = append(prevline.locations, Emitter(SourceLocation(source)))
				} else {
					sline = NewSourceLine(indentation, source)
					snippet.lines = append(snippet.lines, sline)
				}

				prev_inlined = false
			}
		} else {
			switch tokens[0] {
				case "pasteme":
					pasteme := ParsePasteme(source)

					if (pasteme.inlined && len(snippet.lines) > 0) {
						// Combine with previous line
						prevline := &snippet.lines[len(snippet.lines) - 1]
						prevline.locations = append(prevline.locations, Emitter(pasteme))
						prev_inlined = true
					} else {
						sline = SourceLine{indentation, []Emitter{Emitter(pasteme)}}	
						snippet.lines = append(snippet.lines, sline)
					}

				case "snippet":
					// Foreach?
					ParseSnippet(tokens[1], subsnippets, scanner)

				case "endsnippet":
					if len(tokens) >= 3 && tokens[1] == "foreach" {
						// Generate this snippet for each instance of the variable
						loopvar := tokens[2]

						i := -1
						var match *Snippet

						for {
							i, match = subsnippets.FindSnippet(loopvar, i + 1)

							if match == nil {
								break;
							}

							var temp SnippetTable
							temp.AddSnippet(*match)

							generated := snippet.Emit(&temp)
							ParseSnippet(snippet.name, subsnippets, bufio.NewScanner(strings.NewReader(generated)))
						}
					} else {
						subsnippets.AddSnippet(snippet)
					}

					return

				default: 
					// Ignore unknown commands
			}
		}
	}

	// Add current snippet to snippet table
	subsnippets.AddSnippet(snippet)
}

// Parse a file, add snippets to subsnippets
func ParseFile(filename string, subsnippets *SnippetTable) {
	f, err := os.Open(filename)
	check(err)

	// Create scanner for reading lines
	scanner := bufio.NewScanner(f)

	ParseSnippet(filename, subsnippets, scanner)
}

func main () {
	var params SnippetTable
	params.AddSnippet(Snippet{"test_condition", []SourceLine{NewSourceLine(0, "UnitTestCondition1")}})
	params.AddSnippet(Snippet{"test_condition", []SourceLine{NewSourceLine(0, "UnitTestCondition2")}})

	// Get snippets from file 
	ParseFile("unit_tests_template.vhd", &params)

	fmt.Println((&params).FirstSnippet("unit_tests_template.vhd").Emit(&params))
	// fmt.Println((&params).FirstSnippet("unit_tests_template.vhd").lines)
}
package main

import (
	"fmt"
	"os"
	"bufio"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Indentation int

func (n Indentation) PrintIndentation() string {
	str := ""

	i := n
	for i > 0 {
		if i >= 4 {
			str += "\t"
			i -= 4
		} else {
			str += " "
			i -= 1
		}
	}

	return str
}

type SourceLine struct {
	indentation Indentation
	source string
}

func (s SourceLine) Emit(indentation int) string {
	return s.indentation.PrintIndentation() + s.source
}

// Object which can emit indented lines
type Emitter interface {
	Emit(int) string
}

type SourceCommand struct {
	indentation Indentation
	tokens []string
}

func (s SourceCommand) Emit(indentation int) string {
	return s.indentation.PrintIndentation() + "Some command..."
}

// PastemeCommand
type PastemeCommand struct {
	SourceCommand
	inlined bool
	zero bool
	multiple bool
}

// Snippet, list of emitters under a name
type Snippet struct {
	name string
	lines []Emitter
}

func (s Snippet) Emit(indentation int) string {
	// Use a string builder to generate the output file
	var builder strings.Builder

	for _, line := range s.lines {
		builder.WriteString(line.Emit(indentation) + "\n")
	}

	return builder.String()
}

func (s *Snippet) AddLine(line string) {
	// Create source line
	sline := NewSourceLine(line)

	// Try and find a command
	command := ParseCommand(sline)

	var e Emitter = sline
	if command != nil {
		e = command
	} else {
		// Append source line
		s.lines = append(s.lines, e)	
	}
}
// Snippet

func NewSourceLine(line string) SourceLine {
	// Extract whitespace
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

	sline := SourceLine{}
	line = strings.TrimSpace(line)

	if len(line) != 0 {
		sline.indentation = Indentation(indentation)
	}
	
	sline.source = line

	return sline
}

func ParseCommandTokens(command string) Emitter {
	q, _ := regexp.Compile("\\s+")

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

	c := new(PastemeCommand)

	c.tokens = q.Split(command, -1)
	c.zero = rep_zero
	c.multiple = rep_multiple

	return c
}

func ParseCommand(s SourceLine) Emitter {
	r, _ := regexp.Compile("^--\\s*")

	match := r.FindStringIndex(s.source)
	
	if len(match) > 0 {
		// Sanitize command
		command := strings.ToLower(s.source[match[1]:])

		sourcecommand := ParseCommandTokens(command)

		if sourcecommand != nil {
			return sourcecommand
		}
	}

	return nil
}

func main () {
	filename := "unit_tests_template.vhd"
	f, err := os.Open(filename)
	check(err)

	var snip Snippet
	snip.name = filename

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// Read line
		line := scanner.Text()
		
		snip.AddLine(line)
	}

	fmt.Println(snip.Emit(0))
}
package main

import (
    "fmt"
    "os"
)

var VerboseOption bool

func Help() {
    fmt.Println(`    VHDL Dumb Parser
Usage: {Option} (top | -o)

    Option :=
-f      Filename 
-d      Name "snippet string" 

    Options are evaluated left to right

top     Name of snippet to print to stdout
-o      Equivalent to re-stating the most recent -f value
-f      Parse a file, executing all commands line by line
-d      Define a snippet from a given name and string literal.

    Parser Commands

--pasteme Name [Modifier] [Inline]

        If a snippet with a matching Name has been defined, it is
        pasted at the current location. The indentation of the
        snippet is adjusted to match that of the command.

    Name := 
        Name of a snippet

    Modifier := * | + | ?

*       Match all snippets with the correct name, possibly zero
+       Match all snippets with the correct name, at least once
?       Match zero or one snippet with the correct name
        No modifier means accept exactly one match.

        An error is emitted if a command which requires at least
        one match does not match any snippets.

    Inline := \
        Remove newlines before and after paste location

--snippet Name(new)
*
--endsnippet [foreach Name(existing)]

        Define a snippet consisting of the lines between snippet and
        endsnippet.

        If foreach is used, a snippet is defined for each existing
        matching snippet. The generated snippets all have Name(new)

--snippets Name(new)
*
--endsnippets

        Define a new snippet for each line between snippets and
        endsnippets. The generated snippets all have Name(new).`)
}

func Perror(a ...any) {
    fmt.Fprintln(os.Stderr, a...)
}

func Verbose(a ...any) {
    if VerboseOption {
       fmt.Fprintln(os.Stderr, a...)
    }
}
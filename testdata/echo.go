// Simple echo program used for testing processes
// When executed with arguments it echos them back and exits
// When executed without arguments it echoes standard in to standard out.
//
// It always exits with a 0 exit code.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// we never want a non 0 exit code
	defer func() { recover() }()

	// if they called with args echo them back
	if len(os.Args) > 1 {
		s := strings.Join(os.Args[1:], " ")
		fmt.Print(s)
		return
	}

	// otherwise echo standard input
	io.Copy(os.Stdout, os.Stdin)
}

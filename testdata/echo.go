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
	defer func() {
		recover()
	}()

	if len(os.Args) > 1 {
		s := strings.Join(os.Args[1:], " ")
		fmt.Print(s)
		return
	}

	io.Copy(os.Stdout, os.Stdin)
}

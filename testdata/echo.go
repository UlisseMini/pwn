// Simple echo program used for testing processes
// When executed with arguments it echos them back and exits
// When executed without arguments it echoes standard in to standard out.
//
// It always exits with a 0 exit code.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("[*] ")

	// initalize logger
	f, err := os.Create("log.txt")
	if err != nil {
		os.Exit(132) // strange exit code to notify me
	}
	defer f.Close()
	log.SetOutput(f)

	// we never want a non 0 exit code
	defer func() {
		e := recover()
		if e != nil {
			log.Printf("recovered: %v", e)
		}
	}()

	if len(os.Args) > 1 {
		log.Printf("len(os.Args) > 1")
		s := strings.Join(os.Args[1:], " ")
		fmt.Print(s)
		return
	}

	log.Printf("io.Copy(os.Stdout, os.Stdin)")
	_, err = ioCopy(os.Stdout, os.Stdin)
	log.Printf("io.Copy done, err: %v", err)
}

// because i dont trust that sneaky beaky io.Copy (that uses writeTo)
func ioCopy(dst io.Writer, src io.Reader) (int, error) {
	defer log.Println("ioCopy done")

	buf := make([]byte, 1024)
	for {
		nr, err := src.Read(buf[:])
		if err != nil {
			return 0, err
		}
		log.Printf("Read %q", buf[:nr])

		nw, err := dst.Write(buf[:nr])
		if err != nil {
			return 0, err
		}

		log.Printf("Wrote %d bytes", nw)
	}

	return 0, nil
}

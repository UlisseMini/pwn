package pwn

import (
	"io"
	"os"
	"os/exec"
	"time"
)

// Start starts cmd and returns a Process for it
func Start(cmd *exec.Cmd) (Process, error) {
	// setup file descriptors
	stdin, err1 := cmd.StdinPipe()
	stdout, err2 := cmd.StdoutPipe()
	stderr, err3 := cmd.StderrPipe()
	for _, err := range []error{err1, err2, err3} {
		if err != nil {
			return Process{}, err
		}
	}

	err := cmd.Start()
	if err != nil {
		return Process{}, err
	}

	return Process{
		Cmd:    cmd,
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		// the maximum line length to be used with ReadLine
		maxLen: MaxLenDefault,
	}, nil
}

// Spawn spawns a new process and returns it
func Spawn(path string, args ...string) (Process, error) {
	cmd := exec.Command(path, args...)
	return Start(cmd)
}

// Process represents a spawned process
// It has the methods of a os.Process and os.Cmd
type Process struct {
	// the embedded cmd
	*exec.Cmd

	// file descriptors we can manipulate
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser

	// the max length to be used with ReadLine
	maxLen int
}

// WriteLine writes a line to the standard input of the running process
// t can be anything convertable to []byte (see ToBytes function)
// ToBytes will panic if it fails to convert to bytes
func (p Process) WriteLine(t interface{}) error {
	// write the data to the processes standard input
	return WriteLine(p.Stdin, t)
}

// ReadLine reads until newline or timeout expires
// TODO: implement timeout
func (p Process) ReadLine(timeout time.Duration) ([]byte, error) {
	return ReadTill(p.Stdout, p.maxLen, '\n')
}

// Interactive sets the file descriptors to os.Stderr os.Stdout and os.Stdin
func (p Process) Interactive() error {
	return interactive(p, os.Stdin, os.Stdout, os.Stderr)
}

// the actual implementation of Process.Interactive
func interactive(p Process, in io.Reader, out, err io.Writer) error {
	// Copy file descriptors
	go io.Copy(out, p.Stdout)
	go io.Copy(err, p.Stderr)

	// when copying stdin we need to close it when done
	go func() {
		defer p.Stdin.Close()
		io.Copy(p.Stdin, in)
	}()

	// Wait for the process to exit
	return p.Wait()
}

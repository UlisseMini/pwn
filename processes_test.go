// Tests for processes.go
// TODO:
// Instead of relying on external programs,
// call the tests with an envoriment variable set and have the test
// act as the process it is testing.

package pwn

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

// shell to be executed in tests.
var shell string

// Try and make this work on windows, even though its SHIET
func init() {
	switch runtime.GOOS {
	case "windows":
		shell = "cmd.exe"
	default:
		shell = "sh"
	}
}

func TestEcho(t *testing.T) {
	t.Parallel()
	expected := []byte("Hello, world!")
	p, err := Spawn("echo", "Hello, world!")
	if err != nil {
		t.Fatal(err)
	}

	output, err := p.ReadLine(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	// now make sure we got what we expected
	if !bytes.Equal(output, expected) {
		t.Fatalf("wanted %q got %q", expected, output)
	}
}

func TestSh(t *testing.T) {
	expected := []byte("Hello, world")
	p, err := Spawn(shell)
	if err != nil {
		t.Fatal(err)
	}

	err = p.WriteLine("echo Hello, world")
	if err != nil {
		t.Fatal(err)
	}

	out, err := p.ReadLine(time.Second)
	if err != nil {
		t.Fatal(err)
	}

	// now check that we got the expected output
	if !bytes.Equal(out, expected) {
		t.Fatalf("wanted %q got %q", expected, out)
	}
}

// test the interactive function in processes.go
func TestInteractive(t *testing.T) {
	var testcases = []struct {
		// where the process reads standard in from
		stdin io.Reader

		// expected file descriptor outputs
		wantStdout string
		wantStderr string

		// expected error return value
		wantErr error
	}{
		{
			stdin: bytes.NewBufferString("Hello, world\n"),

			wantStdout: "Hello, world\n",
			wantStderr: "",
			wantErr:    nil,
		},
	}

	for _, tc := range testcases {
		// Buffers to use
		outBuf := &bytes.Buffer{}
		errBuf := &bytes.Buffer{}

		// Start the process, it will echo back what is given to it
		cmd := echo()
		p, err := Start(cmd)
		if err != nil {
			t.Fatalf("spawn child process: %v", err)
		}

		// terminate the process after some time, prevent the test from
		// blocking forever if something goes wrong.
		time.AfterFunc(time.Second, func() {
			if err := p.Signal(os.Interrupt); err == nil {
				return
			}

			if err := p.Kill(); err != nil {
				t.Fatalf("failed to kill process: %v", err)
			}
		})

		// wantErr will usually be nil, so this is effectively `if err != nil`
		err = interactive(p, tc.stdin, outBuf, errBuf)
		if err != tc.wantErr {
			t.Fatalf("got error = %v, want error %v", err, tc.wantErr)
		}

		if gotOut := outBuf.String(); gotOut != tc.wantStdout {
			t.Fatalf("got = %q, want %q", gotOut, tc.wantStdout)
		}

		if gotErr := errBuf.String(); gotErr != tc.wantStderr {
			t.Fatalf("got = %q, want %q", gotErr, tc.wantStderr)
		}
	}
}

// wrapper around 'go run testdata/echo.go'
func echo(a ...string) *exec.Cmd {
	args := []string{"run", "testdata/echo.go"}
	args = append(args, a...)

	return exec.Command("go", args...)
}

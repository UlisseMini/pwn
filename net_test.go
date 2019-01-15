package pwn

import (
	"bytes"
	"io"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/golang/net/nettest"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// TestNormalConn uses the nettest package to test that Conn works
// testing readtill is pointless since io_test already covers it as long
// as it has a valid reader.
func TestNormalConn(t *testing.T) {
	t.Parallel()
	nettest.TestConn(t, mp)
}

// TestMaxLen makes sure that Conn.MaxLen works correctly.
//
// Copied from io_test.go -- possibily find a way to reuse code?
func TestMaxLen(t *testing.T) {
	server, client, cleanup, err := mpPwn()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	var testcases = []struct {
		// expected input and expecteds
		input    []byte
		expected []byte

		delim byte
		// do we expect ErrMaxLen?
		overMaxLen bool
		maxLen     int
	}{
		{
			input:    []byte("AAAAAAAAAABBBBBBBBBB"),
			expected: []byte("AAAAAAAAAA"),

			delim: 'B',
		},
		{
			input:    []byte("Hello\n World"),
			expected: []byte("Hello"),
			delim:    '\n',
		},
		{
			// What happens with a nil delim?
			input:    []byte("Hello\n World"),
			expected: []byte("Hello\n World"),
		},
		{
			// test max len
			input:      []byte("Hello, World!"),
			expected:   []byte("Hello,"),
			maxLen:     6,
			overMaxLen: true,
		},
	}

	for _, tc := range testcases {

		// send the client some data
		go func(data []byte) { server.Write(data) }(tc.input)

		// set the client maxLen
		client.MaxLen(tc.maxLen)

		output, err := client.ReadTill(tc.delim)
		if err != nil && err != io.EOF {
			if !tc.overMaxLen && err != ErrMaxLen {
				t.Fatal(err)
			}
		}

		if !bytes.Equal(output, tc.expected) {
			t.Fatalf("wanted %q got %q", tc.expected, output)
		}
	}

	// test that readtill returns correctly on a nil reader
	t.Run("test nil", func(t *testing.T) {
		_, err := ReadTill(nil, 0, '\n')
		if err != ErrNilReader {
			t.Fatalf("expected ErrNilReader, got: %v", err)
		}
	})
}

// mp connects a Listener with a Dialer
// c2 is the client c1 is the server
//
// instead of returning Conn it returns net.Conn
// so nettest can test it.
func mp() (c1, c2 net.Conn, stop func(), err error) {
	return mpPwn()
}

// mp connects a Listener with a Dialer
// c2 is the client c1 is the server
func mpPwn() (c1, c2 Conn, stop func(), err error) {
	addr := "127.0.0.1:" + randPort()
	connChan := make(chan Conn)
	errChan := make(chan error)
	go func() {
		l, err := Listen("tcp", addr)
		if err != nil {
			errChan <- err
			return
		}
		conn, err := l.Accept()
		if err != nil {
			errChan <- err
			return
		}
		connChan <- conn
	}()

	time.Sleep(20 * time.Millisecond)
	c2, err = Dial("tcp", addr)
	if err != nil {
		return
	}

	// check possible error from the server goroutine
	select {
	case err = <-errChan:
		if err != nil {
			return
		}
	case c1 = <-connChan:
		break
	}

	stop = func() {
		c1.Close()
		c2.Close()
	}
	return
}

// get a random port from min 1024 to max 65535
func randPort() string {
	var port int
	for port <= 1024 {
		port = rand.Intn(65535)
	}
	return strconv.Itoa(port)
}

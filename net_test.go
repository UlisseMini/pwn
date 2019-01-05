package pwn

import (
	"bytes"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Test the client reading from the connection
// using my custom ReadTill method.
func TestReadTill(t *testing.T) {
	type testcase struct {
		send     []byte
		expected []byte
		delim    byte
	}

	// readline testcases
	testcases := []testcase{
		testcase{
			send:     []byte("Hello\nThere!"),
			expected: []byte("Hello"),
			delim:    '\n',
		},
		testcase{
			send:     []byte("Hey there"),
			expected: []byte("Hey"),
			delim:    ' ',
		},
		testcase{
			send:     []byte("AAAAAAAABBBBBBBBB"),
			expected: []byte("AAAAAAAA"),
			delim:    'B',
		},
	}

	var serverConn net.Conn
	for _, tc := range testcases {
		port := randPort()

		// get the client connection
		go func() {
			l, err := net.Listen("tcp", "127.0.0.1:"+port)
			if err != nil {
				t.Fatal(err)
			}

			defer l.Close()
			serverConn, err = l.Accept()
			_, err = serverConn.Write(tc.send)
			if err != nil {
				t.Fatal(err)
			}
		}()

		// add a delay to give the server time to start up
		time.Sleep(100 * time.Millisecond)

		// connect to the server, in a function so defer works
		func() {
			c, err := Dial("tcp", "127.0.0.1:"+port)
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()

			// call ReadTill
			output, err := c.ReadTill(tc.delim)
			if err != nil {
				t.Fatal(err)
			}

			// check that output is equal to the expected output
			if !bytes.Equal(output, tc.expected) {
				// if it fails print both, since i'm using text i can do %q
				// but if i was using bytes %X should be used to print the hex.
				t.Fatalf("%q != %q", output, tc.expected)
			}
		}()
	}
}

// get a random port from min 1024 to max 65535
func randPort() string {
	portInt := rand.Intn(65535 - 1024)
	return strconv.Itoa(portInt)
}

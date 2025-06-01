// Run the same tests, but from the command line
package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type mdTest struct {
	out string
	in  string
}

var vectors = []mdTest{
	{"9c1185a5c5e9fc54612808977ee8f548b2258d31", ""},
	{"0bdc9d2d256b3ee9daae347be6f4dc835a467ffe", "a"},
	{"8eb208f7e05d987a9b044a8e98c6b087f15a0bfc", "abc"},
	{"5d0689ef49d2fae572b881b123a85ffa21595f36", "message digest"},
	{"f71c27109c692c1b56bbdceb5b9d2865b3708dbc", "abcdefghijklmnopqrstuvwxyz"},
	{"12a053384a9c0c88e405a06c27dcf49ada62eb2b", "abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"},
	{"b0e20b6e3116640286ed3a87a5713079b21f5189", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"},
	{"9b752e45573d4b39f4dbd3323cab82bf63326bfb", "12345678901234567890123456789012345678901234567890123456789012345678901234567890"},
}

func TestVectors(t *testing.T) {
	// see how our command was called/compiled
	cmdFilename := "./rmd160"
	if _, err := os.Open(cmdFilename); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// try now with ./cmd:
			cmdFilename = "./cmd"
			if _, err := os.Open(cmdFilename); err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					// abort test, cannot find executable command
					cmdFilename = ""
					t.Fatal("cannot find executable in path")
				}
			}
		}
	}

	// add a million 'A', using strings.Builder (memory-efficient allocations)
	var millionA strings.Builder
	millionA.Grow(1000000)
	for range 1000000 {
		millionA.WriteByte('a')
	}

	vectors = append(vectors,
		mdTest{
			out: millionA.String(),
			in:  "52783243c1697bdbe16d37f97f68f08325dc1528",
		})
	// run test
	for i := range len(vectors) {
		tv := vectors[i]

		cmd := exec.Command(cmdFilename, "--quiet")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, tv.in)
		}()

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		cmdOut := padExactly(string(out), 40)

		if cmdOut != tv.out {
			// trim tv.in to the first 100 chars, or else we'll blow up everything
			t.Fatalf("RIPEMD-160(%100s) = %s, expected %s", tv.in, out, tv.out)
		}
	}
}

func padExactly(s string, max int) string {
	if len(s) < max {
		return fmt.Sprintf("%-*s", max, s)
	}
	return s[:max]
}

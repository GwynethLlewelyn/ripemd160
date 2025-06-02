// Run the same tests, but from the command line
package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
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
			out: "52783243c1697bdbe16d37f97f68f08325dc1528",
			in:  millionA.String(),
		})
	// run test
	for i := range len(vectors) {
		tv := vectors[i]

		t.Logf("Vector %02d expected result: %s.\n", i, tv.out)

		cmd := exec.Command(cmdFilename)

		cmdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Logf("connecting to STDIN aborted with %v\n", err)
			t.Fatal(err)
		}
		cmdOut, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatalf("connecting to STDOUT failed with: %v\n", err)
		}
		if err := cmd.Start(); err != nil {
			t.Logf("failed executing %s, error was: %v\n", cmdFilename, err)
			t.Fatal(err)
		}
		if n, err := cmdIn.Write([]byte(tv.in)); err != nil {
			t.Logf("writing to cmd's STDIN failed with %v\n", err)
			t.Fatal(err)
		} else {
			t.Logf("wrote %d bytes to STDIN\n", n)
		}
		if err := cmdIn.Close(); err != nil {
			t.Logf("failed terminating %s, error was: %v\n", cmdFilename, err)
			t.Fatal(err)
		}
		cmdBytes, err := io.ReadAll(cmdOut)
		if err != nil {
			t.Logf("did not read anything from STDOUT: %v\n", err)
			t.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			t.Logf("waiting for command %q aborted with %v\n", cmdFilename, err)
			t.Fatal(err)
		}

		t.Logf("output received from command was: %q\n", cmdBytes)

		output := padExactly(string(cmdBytes), 40)

		if output != tv.out {
			// trim tv.in to the first 100 chars, or else we'll blow up everything
			t.Fatalf("RIPEMD-160(%s) = %s, expected %s ❌", strings.TrimSpace(padExactly(tv.in, 100)), output, tv.out)
		} else {
			t.Logf("%s ✅\n", tv.out)
		}
	}
}

// Truncates (slices) or pads string (with spaces) so that its size is exactly `max` bytes.
// Note: it's not Unicode-safe!
func padExactly(s string, max int) string {
	if len(s) < max {
		return fmt.Sprintf("%-*s", max, s)
	}
	return s[:max]
}

func padExactlyUnicode(s string, max int) string {
	res := make([]rune, max, max)
	var pos, i int
	// truncate
	for pos, char := range s {
		if pos < max {
			res[pos] = char
		}
	}
	// pad
	for i = pos; i < max; i++ {
		res[pos] = ' '
	}
	return string(res)
}

// Simple & fast random string generator by Timothy O. Margheim
// See https://www.timothyomargheim.com/posts/go-tricks-benchmarks/
type Generator struct {
	r *rand.Rand // seed; if nil, a new one will be created on call.
}

const characterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const characterSetUnicode = `シャロー判事、スレンダー判事、そしてヒュー・エヴァンス卿がやって来て、シャロー判事がジョン・フォルスタッフ卿に怒りをぶちまけている理由を議論している。エヴァンスは話題を若いアン・ペイジに移し、スレンダー判事と結婚させたいと思っていることにする。彼らはペイジ判事の玄関に到着し、シャローはフォルスタッフとその取り巻きたちと対峙する。二人は食事をするために部屋に入るが、スレンダー判事は外をうろつき、アン・ペイジが中に入るまで彼女と話をしようと試みるが、うまくいかない。フォルスタッフとその取り巻きたちはガーター・インに落ち着き、そこでファルスタッフはペイジ夫人とフォード夫人を誘惑する計画を明かす。二人は夫の財産を掌握しており、ファルスタッフはそれを狙っている。彼はピストルとニムに手紙を届けさせようとするが、二人は拒否する。二人はペイジとフォードにファルスタッフの計画を阻止しようと企む`

// Returns a random string with up to `numChars` characters.
func (g *Generator) String(numChars int, unicode bool) (string, error) {
	if numChars <= 0 {
		return "", fmt.Errorf("numChars must be greater than 0, received: %d", numChars)
	}
	if numChars > 1024 {
		return "", fmt.Errorf("numChars must not be greater than 1024, received: %d", numChars)
	}

	if g.r == nil {
		src := rand.NewSource(time.Now().UnixNano())
		g.r = rand.New(src)
	}
	randomString := make([]byte, numChars)
	for i := range randomString {
		if !unicode {
			randomString[i] = characterSet[g.r.Intn(len(characterSet))]
		} else {
			randomString[i] = characterSetUnicode[g.r.Intn(len(characterSet))]
		}
	}
	return "anonymised_" + string(randomString), nil
}

// Benchmarking padExactly
func BenchmarkPadExactly(b *testing.B) {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	gen := &Generator{}

	var (
		testString string
		testStrLen int
		err        error // declared here to avoid scoping issues
	)
	for b.Loop() {
		// get one random string; avoid empty strings and too trivial ones.
		testStrLen = rnd.Intn(1000) + 3
		testString, err = gen.String(testStrLen, false)
		if err != nil {
			b.Fatalf("random string generation failed: %s\n", err)
		}
		padTo := rnd.Intn(testStrLen) + 3
		// now pad it to some value
		// b.Logf("%s - padded string of size %d to %d characters\n",
		// 	padExactly(testString, padTo),
		// 	testStrLen,
		// 	padTo,
		// )
		padExactly(testString, padTo)
	}
}

// Benchmarking padExactly
func BenchmarkPadExactlyUnicode(b *testing.B) {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	gen := &Generator{}

	var (
		testString string
		testStrLen int
		err        error // declared here to avoid scoping issues
	)
	for b.Loop() {
		// get one random string; avoid empty strings and too trivial ones.
		testStrLen = rnd.Intn(1000) + 3
		testString, err = gen.String(testStrLen, true)
		if err != nil {
			b.Fatalf("random string generation failed: %s\n", err)
		}
		padTo := rnd.Intn(testStrLen) + 3
		// now pad it to some value
		// b.Logf("%s - padded string of size %d to %d characters\n",
		// 	padExactlyUnicode(testString, padTo),
		// 	testStrLen,
		// 	padTo,
		// )
		padExactlyUnicode(testString, padTo)
	}
}

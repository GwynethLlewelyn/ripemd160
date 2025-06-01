package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/antelope-go/ripemd160"
)

// versionInfoType holds the relevant information for this build.
// It is meant to be used as a cache.
type versionInfoType struct {
	version    string    // Runtime version.
	commit     string    // Commit revision number.
	dateString string    // Commit revision time (as a RFC3339 string).
	date       time.Time // Same as before, converted to a time.Time, because that's what the cli package uses.
	builtBy    string    // User who built this (see note).
	goOS       string    // Operating system for this build (from runtime).
	goARCH     string    // Architecture, i.e., CPU type (from runtime).
	goVersion  string    // Go version used to compile this build (from runtime).
	init       bool      // Have we already initialised the cache object?
}

// NOTE: I don't know where the "builtBy" information comes from, so, right now, it gets injected
// during build time, e.g. `go build -ldflags "-X main.TheBuilder=gwyneth"` (gwyneth 20231103)

var (
	versionInfo versionInfoType // cached values for this build.
	TheBuilder  string          // to be overwritten via the linker command `go build -ldflags "-X main.TheBuilder=gwyneth"`.
	TheVersion  string          // to be overwritten with -X main.TheVersion=X.Y.Z, as above.
	debugLevel  int             // verbosity/debug level.
)

// Initialises the versionInfo variable.
func initVersionInfo() error {
	if versionInfo.init {
		// already initialised, no need to do anything else!
		return nil
	}
	// get the following entries from the runtime:
	versionInfo.goOS = runtime.GOOS
	versionInfo.goARCH = runtime.GOARCH
	versionInfo.goVersion = runtime.Version()

	// attempt to get some build info as well:
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("no valid build information found")
	}
	// use our supplied version instead of the long, useless, default Go version string.
	if TheVersion == "" {
		versionInfo.version = buildInfo.Main.Version
	} else {
		versionInfo.version = TheVersion
	}

	// Now dig through settings and extract what we can...

	var vcs, rev string // Name of the version control system name (very likely Git) and the revision.
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs":
			vcs = setting.Value
		case "vcs.revision":
			rev = setting.Value
		case "vcs.time":
			versionInfo.dateString = setting.Value
		}
	}
	versionInfo.commit = "unknown"
	if vcs != "" {
		versionInfo.commit = vcs
	}
	if rev != "" {
		versionInfo.commit += " [" + rev + "]"
	}
	// attempt to parse the date, which comes as a string in RFC3339 format, into a date.Time:
	var parseErr error
	if versionInfo.date, parseErr = time.Parse(versionInfo.dateString, time.RFC3339); parseErr != nil {
		// Note: we can safely ignore the parsing error: either the conversion works, or it doesn't, and we
		// cannot do anything about it... (gwyneth 20231103)
		// However, the AI revision bots dislike this, so we'll assign the current date instead.
		versionInfo.date = time.Now()

		if debugLevel > 1 {
			fmt.Fprintf(os.Stderr, "date parse error: %v", parseErr)
		}
	}

	// see comment above
	versionInfo.builtBy = TheBuilder

	return nil
}

// checksumOneFile calls the ripemd160 hash service with the `r` io.Reader
// and attempts to calculate its checksum.
func checksumOneFile(r io.Reader) ([]byte, error) {
	md := ripemd160.New()

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if n, err := io.WriteString(md, string(buf)); err != nil {
		printInfo("could not generate hash from input; error was %v\n", err)
		return nil, err
	} else {
		printInfo("bytes read: %d\n", n)
	}

	// Truncate huge files and add an elypsis.
	if len(buf) > 1024 {
		buf = fmt.Append(buf[:1024], " [...]")
	}
	printInfo("contents to checksum: %q\n", buf)

	res := md.Sum(nil)

	printInfo("Checksum: %x (%d bytes)\n", res, len(res))

	return res, nil
}

// printChecksum, given a valid `r` io.Reader, prints out its RIPE MD-160 checksum.
// Uses "-" to display STDIN, just like `sha256sum` and the other GNU tools.
func printChecksum(fname string, r io.Reader) (err error) {
	// potentially unneeded extra validation step.
	if len(fname) == 0 {
		fname = "-"
	}
	if checksum, err := checksumOneFile(r); err == nil {
		if setting.Tag {
			// BSD style, using --tag
			fmt.Printf("RMD160 (%s) = %x\n", fname, checksum)
		} else {
			// default style
			fmt.Printf("%x  %s\n", checksum, fname)
		}
	} else {
		return fmt.Errorf("error %v\n", err)
	}
	return nil
}

// printInfo is a simple wrapper to print debugging info, if desired, or
// skipping under normal operation.
// TODO(gwyneth): do this properly using a modern logging library.
func printInfo(fmtStr string, args ...any) {
	if !setting.Quiet && setting.Debug {
		fmt.Fprintf(os.Stderr, fmtStr, args...)
	}
}

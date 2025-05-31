package main

import (
	"fmt"
	"time"
	"os"
	"runtime"
	"runtime/debug"
)

// versionInfoType holds the relevant information for this build.
// It is meant to be used as a cache.
type versionInfoType struct {
	version		string		// Runtime version.
	commit  	string		// Commit revision number.
	dateString 	string		// Commit revision time (as a RFC3339 string).
	date		time.Time	// Same as before, converted to a time.Time, because that's what the cli package uses.
	builtBy 	string		// User who built this (see note).
	goOS		string		// Operating system for this build (from runtime).
	goARCH		string		// Architecture, i.e., CPU type (from runtime).
	goVersion	string		// Go version used to compile this build (from runtime).
	init		bool		// Have we already initialised the cache object?
}

// NOTE: I don't know where the "builtBy" information comes from, so, right now, it gets injected
// during build time, e.g. `go build -ldflags "-X main.TheBuilder=gwyneth"` (gwyneth 20231103)

var (
	versionInfo versionInfoType	// cached values for this build.
	TheBuilder string			// to be overwritten via the linker command `go build -ldflags "-X main.TheBuilder=gwyneth"`.
	debugLevel int				// verbosity/debug level.
)

// Initialises the versionInfo variable.
func initVersionInfo() error {
	if versionInfo.init {
		// already initialised, no need to do anything else!
		return nil
	}
	// get the following entries from the runtime:
	versionInfo.goOS		= runtime.GOOS
	versionInfo.goARCH		= runtime.GOARCH
	versionInfo.goVersion	= runtime.Version()

	// attempt to get some build info as well:
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("no valid build information found")
	}
	versionInfo.version = buildInfo.Main.Version

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

	// NOTE: I have no idea where the "builtBy" info is supposed to come from;
	// the way I do it is to force the variable with a compile-time option. (gwyneth 20231103)
	versionInfo.builtBy = TheBuilder

	return nil
}

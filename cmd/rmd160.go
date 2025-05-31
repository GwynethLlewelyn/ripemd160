// Very basic command-tool to do rmd160 checksums.
package main

import (
	"context"
	"fmt"
	"io"
	"net/mail"
	"os"
	"time"

	"github.com/antelope-go/ripemd160"
	"github.com/urfave/cli/v3"
)

// Setting has all the possible settings for this app.
type Setting struct {
	BinaryOrText	bool	// either read in binary or text mode; ignored.
	Check			bool	// read checksums from the FILEs and check them.
	Tag				bool	// create a BSD-style checksum.
	Zero			bool	// end each output line with NUL, not newline, and disable file name escaping.
	IgnoreMissing	bool	// TODO
	Quiet			bool	// don't print OK for each successfully verified file.
	Status			bool	// TODO
	Strict			bool	// TODO
}

// The current set of settings, available to all functions.
var setting Setting

func main() {
	// Set up the version/runtime/debug-related variables, and cache them.
	if err := initVersionInfo(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize version info: %v\n", err)
	}


	// start app
	app := &cli.Command{
		Name:      "rmd160",
		Usage:     "",
		UsageText: os.Args[0] + " [OPTION]... [FILE]...",
		Version: fmt.Sprintf(
			"%s (rev %s) [%s %s %s] [build at %s by %s]",
			versionInfo.version,
			versionInfo.commit,
			versionInfo.goOS,
			versionInfo.goARCH,
			versionInfo.goVersion,
			versionInfo.dateString,		// Date as string in RFC3339 notation.
			versionInfo.builtBy,		// `go build -ldflags "-X main.TheBuilder=[insertname here]"`
		),
		EnableShellCompletion: true,
//		Compiled: versionInfo.date,		// Converted from RFC333
		Authors: []any{
			mail.Address{Name: "pnx", Address: "henrik.hautakoski@gmail.com"},
			mail.Address{Name: "Gwyneth Llewelyn", Address: "hgwyneth.llewelyn@gwynethllewelyn.net"},
		},
		Copyright: fmt.Sprintf("© 2024-%d by Henrik Hautakoksi. All rights reserved. Freely distributed under a 3-clause-BSD license.", time.Now().Year()),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "binary",
				Aliases: []string{"d"},
				Usage:   "read in binary mode (ignored)",
				Value:   false,
				Destination: &setting.BinaryOrText,
			},
			&cli.BoolFlag{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "read checksums from the FILEs and check them",
				Value:   false,
				Destination: &setting.Check,
			},
			&cli.BoolFlag{
				Name:    "tag",
				Usage:   "create a BSD-style checksum",
				Value:   false,
				Destination: &setting.Tag,
			},
			&cli.BoolFlag{
				Name:    "text",
				Aliases: []string{"t"},
				Usage:   "read in text mode (ignored)",
				Value:   false,
				Destination: &setting.BinaryOrText,
			},
			&cli.BoolFlag{
				Name:    "zero",
				Aliases: []string{"z"},
				Usage:   "end each output line with NUL, not newline, and disable file name escaping",
				Value:   false,
				Destination: &setting.Zero,
			},
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "don't print OK for each successfully verified fileg",
				Value:   false,
				Destination: &setting.Quiet,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// Everything happens here!
			var fname string
			printInfo("Processing %d files...\n", cmd.Args().Len())
			if cmd.Args().Len() >= 0 {
				for i := range cmd.Args().Len() {
					fname = cmd.Args().Get(i)
					printInfo("\t%02d:\t%q\n", i, fname)
					if fname == "-" {
						if checksum, err := checksumOneFile(os.Stdin); err != nil {
							fmt.Printf("%s\t%s\n", checksum, fname)
						} else {
							cli.Exit(fmt.Sprintf("error checksumming file %q: %s\n", fname, err), 10)
						}
						break
					}
					if ioRead, err := os.Open(fname); err == nil {
						if checksum, err := checksumOneFile(ioRead); err != nil {
							fmt.Printf("%s\t%s\n", checksum, fname)
						} else {
							cli.Exit(fmt.Sprintf("error checksumming file %q: %s\n", fname, err), 10)
						}
					} else {
						return cli.Exit(fmt.Sprintf("error opening file %q: %s\n", fname, err), 1)
					}
				}
			} else {
				// handle stdin only
				if checksum, err := checksumOneFile(os.Stdin); err != nil {
					fmt.Printf("%s\t%s\n", checksum, fname)
				} else {
					cli.Exit(fmt.Sprintf("error checksumming file %q: %s\n", fname, err), 10)
				}
			}
			return nil
		},
	}
	if err := app.Run(context.Background(), os.Args); err != nil {
		printInfo("Run failed: %v\n", err)
	}
}

// checksumOneFile calls
func checksumOneFile(r io.Reader)([]byte, error) {
	md := ripemd160.New()

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	printInfo("contents of file: %s\n", buf)

	res := md.Sum(buf)

	printInfo("Checksum: %v (%d bytes)\n", md, len(res))

	return res, nil
}

func printInfo(fmtStr string, args ...any) {
	if !setting.Quiet {
		fmt.Fprintf(os.Stderr, fmtStr, args...)
	}
}
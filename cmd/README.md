![RIPE MD160 Logo]
# CLI for rmd160

## Usage

Pretty much as similar tools (e.g. `sha256sum` and other similar hashing tools from the GNU project).

```sh
   ./rmd160 [OPTION]... [FILE]...
```
With no `FILE`, or when `FILE` is `-`, read from standard input.

### Global Options

|	Option		|		Description												|
| ------------- | ------------------------------------------------------------- |
| --binary, -d	| read in binary mode (ignored) (default: false)				|
| --check, -c	| read checksums from the FILEs and check them (default: false)	|
| --tag			| create a BSD-style checksum (default: false)					|
| --text, -t	| read in text mode (ignored) (default: false)					|
| --zero, -z	| end each output line with NUL, not newline, and disable file name escaping (default: false) |
| --quiet, -q	| don't print OK for each successfully verified file (default: false) |
| --debug, -d	| shows additional debugging information (default: false)		|
| --help, -h	| show help														|
| --version, -v	| print the version												|


`--check` and `--zero` haven't been fully implemented yet (therefore, `--quiet` will do nothing, either).

**TBD:** more strict compliance with the options offered by the hashing GNU tools, so that it can be a drop-in replacement for them.  

## Compilation

Compile with: `go build -ldflags "-X main.TheBuilder=[INSERT YOUR NAME HERE] -X main.TheVersion=[OVERRIDE AUTOMATED GO VERSIONING]" -o rmd160`

`go test` assumes that a working executable has been placed in the same directory and will spawn it in order to execute precisely the same test vectors as the main package.

## License
Same as the original package ([BSD-3-Clause]).

## Acknowledgements

- [@pnx], of course, for his RIPE MD-160 implementation and testing battery;
- [@urfave], who never earns much credit for his fantastic CLI-building library for Go.

[![go-test](https://github.com/GwynethLlewelyn/ripemd160/actions/workflows/ci.yml/badge.svg)](https://github.com/GwynethLlewelyn/ripemd160/actions/workflows/ci.yml) [![CodeQL Advanced](https://github.com/GwynethLlewelyn/ripemd160/actions/workflows/codeql.yml/badge.svg)](https://github.com/GwynethLlewelyn/ripemd160/actions/workflows/codeql.yml)

[RIPE MD160 Logo]: ../assets/rmd160-logo-small.png
[BSD-3-Clause]: ../LICENSE.md
[@pnx]: https://github.com/pnx
[@urfave]: https://github.com/urfave
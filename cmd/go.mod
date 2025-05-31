module github.com/antelope-go/ripemd160/cmd

go 1.24.3

// This should be deleted at the end!
replace (
	github.com/antelope-go/ripemd160 => ..
	github.com/antelope-go/ripemd160/cmd => .
)

require (
	github.com/antelope-go/ripemd160 v1.0.0
	github.com/urfave/cli/v3 v3.3.3
)

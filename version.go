package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	version = "dev"
	commit  = "-"
	date    = ""
)

func displayVersion() {
	if info, available := debug.ReadBuildInfo(); available {
		if date == "" {
			version = info.Main.Version
			commit = fmt.Sprintf("(unknown, mod sum: %q)", info.Main.Sum)
			date = "(unknown)"
		}
	}

	fmt.Printf(`motoko:
 version     : %s
 commit      : %s
 build date  : %s
 go version  : %s
 go compiler : %s
 platform    : %s/%s
`, version, commit, date, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)
}

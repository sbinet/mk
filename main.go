package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const tmpl = `
## simple makefile to log workflow
.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install test


build:
	@go build $(GOFLAGS) ./...

install:
	@go get $(GOFLAGS) ./...

test: install
	@go test $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
`

const mkfname = ".makefile-mkgo.mk"

var (
	g_verbose  = flag.Bool("v", false, "enable verbose output")
	g_makefile = flag.String("f", mkfname, "alternate Makefile")
	g_race     = flag.Bool("race", false, "enable build with race detector")
	g_compiler = flag.String("compiler", "gc", "go compiler to use (gc,gccgo,llgo)")
)

func main() {
	flag.Parse()

	rc := run()
	os.Exit(rc)
}

func run() int {
	var err error
	if *g_makefile == mkfname {
		err = ioutil.WriteFile(*g_makefile, []byte(tmpl), 0644)
		defer os.Remove(*g_makefile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error** creating file [%s]: %v\n", *g_makefile, err)
			return 1
		}
	}

	cmdargs := make([]string, 0, flag.NArg())
	cmdargs = append(cmdargs, "-f", *g_makefile)
	goflags := os.Getenv("GOFLAGS")
	if *g_verbose {
		goflags += " -v"
	}
	if *g_race {
		goflags += " -race"
	}
	goflags += fmt.Sprintf(" -compiler=%s", *g_compiler)

	os.Setenv("GOFLAGS", goflags)

	cmdargs = append(cmdargs, flag.Args()...)

	cmd := exec.Command("make", cmdargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** running make %s: %v\n",
			strings.Join(cmdargs, " "),
			err,
		)
		return 1
	}

	return 0
}

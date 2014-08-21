package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// FIXME: use -ldflags to automatically get the Version ?
// go build -ldflags "-X main.Version 0.1.`date -u +.%Y%m%d%.H%M%S`" main.go

const Version = "0.1"

const tmpl = `## AUTOMATICALLY generated by mk.

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

bench: install
	@go test -bench=. $(GOFLAGS) ./...

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

	g_show         = flag.Bool("show", false, "dump the Makefile mk will use on STDOUT")
	g_show_default = flag.Bool("show-default", false, "dump the default Makefile-mk on STDOUT")
	g_version      = flag.Bool("version", false, "dump mk's version")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `mk - a simple makefile generator for Go packages.

Usage of mk:

$ mk [options] [target-name]

Examples:

$ mk
$ mk -f Makefile-my.mk
$ mk install

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *g_version {
		fmt.Fprintf(os.Stdout, "mk version: %v\n", Version)
		os.Exit(0)
	}

	if *g_show_default {
		fmt.Fprintf(os.Stdout, tmpl)
		os.Exit(0)
	}

	// detect if there is a Makefile in the current directory.
	// use it instead (if user didn't specify an alternate Makefile)
	if *g_makefile == mkfname {
		if _, err := os.Stat("Makefile"); err == nil {
			*g_makefile = "Makefile"
		}
	}

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

	if *g_show {
		mkfile, err := os.Open(*g_makefile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error** opening file [%s]: %v\n", *g_makefile, err)
			return 1
		}
		defer mkfile.Close()
		_, err = io.Copy(os.Stdout, mkfile)
		return 0
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

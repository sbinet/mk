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
.PHONY: all test clean build

GOFLAGS ?= $(GOFLAGS:-v)

all: build test
	@echo "## bye."

build:
	@go get $(GOFLAGS) ./...

test: build
	@go test $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
`

var g_verbose = flag.Bool("v", false, "enable verbose output")
var g_makefile = flag.String("f", ".makefile-mkgo.mk", "alternate Makefile")
var g_race = flag.Bool("race", false, "enable build with race detector")

func main() {
	flag.Parse()

	err := ioutil.WriteFile(*g_makefile, []byte(tmpl), 0644)
	defer os.Remove(*g_makefile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** creating file [%s]: %v\n", *g_makefile, err)
		os.Exit(1)
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
		os.Exit(1)
	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// FIXME: use -ldflags to automatically get the Version ?
// go build -ldflags "-X main.Version 0.1.`date -u +.%Y%m%d%.H%M%S`" main.go

const Version = "0.4.0"

var (
	verbose  = flag.Bool("v", false, "enable verbose output")
	race     = flag.Bool("race", false, "enable build with race detector")
	compiler = flag.String("compiler", "gc", "go compiler to use (gc,gccgo,llgo)")
	version  = flag.Bool("version", false, "dump mk's version")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `mk - a simple command to develop Go packages.

Usage of mk:

$ mk [options] [target-name]

Examples:

$ mk
$ mk install

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	log.SetPrefix("mk: ")
	log.SetFlags(0)

	if *version {
		log.Printf("version: %v\n", Version)
		os.Exit(0)
	}

	rc := run()
	os.Exit(rc)
}

func run() int {
	var (
		err error
	)

	cmdargs := make([]string, 0, flag.NArg())
	goflags := os.Getenv("GOFLAGS")
	if *verbose {
		goflags += " -v"
	}
	if *race {
		goflags += " -race"
	}
	goflags += fmt.Sprintf(" -compiler=%s", *compiler)

	os.Setenv("GOFLAGS", goflags)

	switch flag.NArg() {
	case 0:
		cmdargs = append(cmdargs, "get", "-v", "./...")
	default:
		cmdargs = append(cmdargs, flag.Args()...)
	}

	cmd := exec.Command("go", cmdargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("error running mk %s: %v\n",
			strings.Join(cmdargs, " "),
			err,
		)
		return 1
	}

	return 0
}

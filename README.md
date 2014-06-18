mk
==

`mk` is a simple `make` like command to ease the day-to-day life of
building and testing `go` based projects.

## Installation

```sh
$ go get github.com/sbinet/mk
```

## Usage

In a `go` package:

```sh
# build, install and test
$ mk

# build, install
$ mk install

# only build
$ mk build

# build and test
$ mk test
```


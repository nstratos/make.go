package main

import (
	"fmt"
	"os"
	"runtime"
)

var (
	version = "devel"
)

func main() {
	fmt.Printf("%s %s (runtime: %s)\n", os.Args[0], version, runtime.Version())
}

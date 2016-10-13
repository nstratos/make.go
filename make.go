// +build ignore

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

const (
	commandName = "example-command"
)

type platform struct {
	os   string
	arch string
}

type binary struct {
	name    string
	version string
	targets []platform
}

// Name returns the name of the binary with a certain format for a platform.
func (bin binary) Name(os, arch string) string {
	s := fmt.Sprintf("%s_%s_%s-%s", bin.name, bin.version, os, arch)
	if os == "windows" {
		s = s + ".exe"
	}
	return s
}

// Names returns the name of the binary for all the target platforms.
func (bin binary) Names() []string {
	names := make([]string, len(bin.targets))
	for i, t := range bin.targets {
		names[i] = bin.Name(t.os, t.arch)
	}
	return names
}

var (
	release   = flag.Bool("release", false, "build binaries for all target platforms")
	clean     = flag.Bool("clean", false, "remove all created binaries from current directory")
	buildOS   = flag.String("os", runtime.GOOS, "set operating system to build for")
	buildArch = flag.String("arch", runtime.GOARCH, "set architecture to build for")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: go run make.go [OPTIONS]\n\n")
	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	bin := binary{
		name: commandName,
		// Change these according to the platforms you would like to support. A
		// full list can be found here:
		// https://golang.org/doc/install/source#environment
		targets: []platform{
			{os: "linux", arch: "386"}, {os: "linux", arch: "amd64"},
			{os: "windows", arch: "386"}, {os: "windows", arch: "amd64"},
			{os: "darwin", arch: "386"}, {os: "darwin", arch: "amd64"},
		},
	}
	bin.version = getVersion()

	if *release {
		forEachBinTarget(bin, buildBinary)
		os.Exit(0)
	}

	if *clean {
		forEachBinTarget(bin, rmBinary)
		os.Exit(0)
	}

	buildBinary(bin, *buildOS, *buildArch)
}

// getVersion returns the version of the project using git.
func getVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--always")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error running git describe: %v", err)
	}
	return strings.TrimPrefix(strings.TrimSpace(string(out)), "v")
}

type binaryFunc func(bin binary, OS, arch string)

// forEachBinTarget runs a function for all the target platforms of a binary in
// parallel.
func forEachBinTarget(bin binary, fn binaryFunc) {
	var wg sync.WaitGroup
	for _, t := range bin.targets {
		wg.Add(1)
		go func(bin binary, os, arch string) {
			defer wg.Done()
			fn(bin, os, arch)
		}(bin, t.os, t.arch)
	}
	wg.Wait()
}

// buildBinary builds a binary for a certain operating system and architecture
// combination while using --ldflags to set the binary's version variable.
func buildBinary(bin binary, OS, arch string) {
	ldflags := fmt.Sprintf("--ldflags=-X main.version=%s", bin.version)
	cmd := exec.Command("go", "build", ldflags, "-o", bin.Name(OS, arch))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = copyOSEnv()
	cmd.Env = setEnv(cmd.Env, "GOOS", OS)
	cmd.Env = setEnv(cmd.Env, "GOARCH", arch)
	fmt.Println("Building binary:", bin.Name(OS, arch))
	if err := cmd.Run(); err != nil {
		log.Fatalln("Error running go build:", err)
	}
}

// rmBinary removes a binary from the current directory.
func rmBinary(bin binary, OS, arch string) {
	err := os.Remove(bin.Name(OS, arch))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Error removing binary:", err)
		}
	}
}

// copyOSEnv returns a copy of the system's environment variables.
func copyOSEnv() (environ []string) {
	return append(environ, os.Environ()...)
}

// setEnv searches in a slice of environment variables with the form key=value
// for a key and if found it sets its value, otherwise it adds the pair.
func setEnv(environ []string, key, value string) []string {
	for i, env := range environ {
		if strings.Split(env, "=")[0] == key {
			environ[i] = fmt.Sprintf("%s=%s", key, value)
			return environ
		}
	}
	return append(environ, fmt.Sprintf("%s=%s", key, value))
}

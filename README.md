# make.go

A Go script that could replace your Makefile.

## Features

- Automated versioning of produced binaries.
- Parallel cross compilation for all target platforms.

## Usage

This script is meant to live in your project's root and used with `go run
make.go`.

- `go run make.go` builds a versioned binary for the current platform.
- `go run make.go -os windows -arch amd64` builds a versioned binary for a
  specific platform.
- `go run make.go --release` builds versioned binaries for all target
  platforms.
- `go run make.go --clean` removes produced binaries.

## Rationale

Having quite a few Go projects that produce a single binary, I was looking for
a good way to have automated versioning and easily produce binaries for all the
platforms that I want to support. Of course, Go can easily cross compile
binaries and by using `--ldflags` when building we can also add a version to
our produced binary. Naturally I wrote a Makefile that did all this.

But then I thought, how would this look if it was written in Go?

## Reasons to use a make.go instead of a more succinct Makefile

- Your project has a non trivial building procedure that could benefit from
  being expressed in robust Go code.
- You want the increased portability of Go.
- You like writing Go!

## How to use this make.go

This is by no means a generic script that can be used without modifications but
you can use it as a starting point. The idea is to have a make.go specific for
your project that could act as your build script and possibly do other more
complex stuff. The source code is public domain so if you find the content
useful and want to use it as a starting point, simply copy make.go, modify it
and extend it according to the needs of your project.

## Disclaimer

Now I am not suggesting that you should go and replace all your Makefiles with
Go scripts. I wrote this as an experiment for my specific case, I thought the
code might be useful to others and decided to share it. If a Makefile or other
script gets the job done then you should use that.

## Good alternatives

- A Makefile or other script.
- [Gox](https://github.com/mitchellh/gox), my personal favorite cross compile
  tool.

## Contributing

If you would like to add to this basic make.go, you are welcome to do so.

## License

The source code is public domain.


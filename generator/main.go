package main

import (
	"flag"

	"github.com/frigus02/website/generator/build"
)

func main() {
	var in string
	var out string
	var help bool
	flag.StringVar(&in, "in", ".", "Path to the site folder")
	flag.StringVar(&out, "out", "build", "Path to output folder")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	w := build.Watch{In: in, Out: out}
	w.Watch()
}

package main

import (
	"flag"
	"log"

	"github.com/frigus02/website/generator/build"
)

func main() {
	var mode string
	var in string
	var out string
	var minify bool
	var help bool
	flag.StringVar(&mode, "mode", "build", "Mode: build or watch")
	flag.StringVar(&in, "in", ".", "Path to the site folder")
	flag.StringVar(&out, "out", "build", "Path to output folder")
	flag.BoolVar(&minify, "minify", false, "Minify HTML and CSS output")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	b := build.Build{In: in, Out: out, Minify: minify}
	switch mode {
	case "build":
		b.Build()
	case "watch":
		b.Watch()
	default:
		log.Fatal("Unknown mode")
	}
}

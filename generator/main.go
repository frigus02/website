package main

import (
	"flag"
	"log"

	"github.com/frigus02/website/generator/build"
	"github.com/frigus02/website/generator/serve"
)

func main() {
	var mode string
	var in string
	var out string
	var prod bool
	var help bool
	flag.StringVar(&mode, "mode", "build", "Mode: build or watch")
	flag.StringVar(&in, "in", ".", "Path to the site folder")
	flag.StringVar(&out, "out", "build", "Path to output folder")
	flag.BoolVar(&prod, "prod", false, "Concat CSS and minify HTML and CSS output")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	b := build.Build{
		In:        in,
		Out:       out,
		Minify:    prod,
		ConcatCSS: prod,
	}
	switch mode {
	case "build":
		b.Build()
	case "watch":
		go serve.Serve(out)
		b.Watch()
	default:
		log.Fatal("Unknown mode")
	}
}

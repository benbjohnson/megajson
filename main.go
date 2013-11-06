package main

import (
	"flag"
	"log"

	"github.com/benbjohnson/megajson/generator"
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	path := flag.Arg(0)
	options := generator.NewOptions()
	if err := generator.Generate(path, options); err != nil {
		log.Fatalln(err)
	}
}

func usage() {
	log.Fatal("usage: megajson OPTIONS FILE")
}

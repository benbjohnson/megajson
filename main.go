package main

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {

	}

	// Open the provided path.
	path := flag.Arg(0)
	f, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// Process the AST and generate a parser.
	//g := generator.NewDirGenerator()
	//err := g.Generate(path)
	//source.Codegen()
}

// Prints the CLI usage.
func usage() {
	log.Fatal("usage: megajson OPTIONS FILE")
}

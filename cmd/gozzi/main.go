package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/generator"
	"github.com/tduyng/gozzi/internal/parser"
	"github.com/tduyng/gozzi/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	site, err := config.LoadSite("config.toml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	contentParser := parser.NewParser(site)
	if err := contentParser.Parse("content"); err != nil {
		log.Fatalf("Error parsing content: %v", err)
	}

	gen, err := generator.NewGenerator(site, contentParser)
	if err != nil {
		log.Fatalf("Error creating generator: %v", err)
	}

	switch os.Args[1] {
	case "build":
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			log.Fatal(err)
		}
	case "serve":
		srv, _ := server.NewDevServer(site, gen, contentParser)
		srv.Start(1313)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Usage: gozzi <command>
Commands:
  build    Generate static site
  serve    Start development server`)
}

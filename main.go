package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/generator"
	"github.com/tduyng/gozzi/internal/parser"
	"github.com/tduyng/gozzi/internal/server"
)

var (
	version   = "dev"     // Set during build
	buildTime = "unknown" // Set during build
	commit    = "HEAD"    // Set during build
)

func main() {
	if len(os.Args) >= 2 && (os.Args[1] == "version" || os.Args[1] == "-v" || os.Args[1] == "--version") {
		printVersion()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		startTime := time.Now()
		_, contentParser, gen := initApp("config.toml", "content")
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			log.Fatal(err)
		}
		ms := time.Since(startTime).Milliseconds()
		log.Printf("Build done in %dms", ms)
	case "serve":
		site, contentParser, gen := initApp("config.toml", "content")
		srv, _ := server.NewDevServer(site, gen, contentParser)
		srv.Start(1313)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printVersion() {
	fmt.Printf("gozzi version %s\n", version)
	fmt.Printf("Build time:    %s\n", buildTime)
	fmt.Printf("Git commit:    %s\n", commit)
}

func printUsage() {
	fmt.Println(`Usage: gozzi <command> [options]

Commands:
  build    Generate static site
  serve    Start development server

Options:
  -v, --version  Show version information`)
}

func initApp(configPath, contentDir string) (*config.Site, *parser.ContentParser, *generator.Generator) {
	site, err := config.LoadSite(configPath)
	if err != nil {
		log.Fatalf("Error loading config %s: %v", configPath, err)
	}

	contentParser := parser.NewParser(site)
	if err := contentParser.Parse(contentDir); err != nil {
		log.Fatalf("Error parsing content from %s: %v", contentDir, err)
	}

	gen, err := generator.NewGenerator(site, contentParser)
	if err != nil {
		log.Fatalf("Error creating generator: %v", err)
	}
	return site, contentParser, gen
}

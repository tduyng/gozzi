// Gozzi is a fast static site generator built with Go.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/generator"
	"github.com/tduyng/gozzi/app/parser"
	"github.com/tduyng/gozzi/app/server"
)

var (
	version   string
	buildTime string
	commit    string
)

func main() {
	// Handle global flags
	flag.Usage = printUsage
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("help", false, "Show help information")

	// Parse global flags
	flag.Parse()

	// Handle global flags
	if *showVersion {
		printVersion()
		return
	}
	if *showHelp {
		printUsage()
		return
	}

	// Get subcommand
	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	// Handle subcommands
	switch args[0] {
	case "build":
		handleBuildCommand(args[1:])
	case "serve":
		handleServeCommand(args[1:])
	case "help":
		handleHelpCommand(args[1:])
	case "version":
		printVersion()
	default:
		printUsage()
		os.Exit(1)
	}
}

func handleBuildCommand(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	configPath := fs.String("config", "config.toml", "Path to config file")
	contentDir := fs.String("content", "content", "Content directory")
	fs.Usage = func() { buildUsage() }

	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}

	startTime := time.Now()
	_, contentParser, gen := initApp(*configPath, *contentDir)

	if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
		log.Fatal(err)
	}
	log.Printf("Build done in %dms", time.Since(startTime).Milliseconds())
}

func handleServeCommand(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	configPath := fs.String("config", "config.toml", "Path to config file")
	contentDir := fs.String("content", "content", "Content directory")
	port := fs.Int("port", 1313, "Port to listen on")
	fs.Usage = func() { serveUsage() }

	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}

	site, contentParser, gen := initApp(*configPath, *contentDir)
	srv, err := server.NewDevServer(
		*configPath,
		*contentDir,
		site,
		gen,
		contentParser,
	)
	if err != nil {
		log.Fatal(err)
	}
	srv.Start(*port)
}

func handleHelpCommand(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "build":
			buildUsage()
		case "serve":
			serveUsage()
		default:
			printUsage()
		}
		return
	}
	printUsage()
}

func printVersion() {
	if version == "" {
		if info, ok := debug.ReadBuildInfo(); ok {
			if v := info.Main.Version; v != "" && v != "(devel)" {
				fmt.Printf("gozzi version %s\n", v)
			}
		}
	} else {
		fmt.Printf("gozzi version %s\n", version)
		fmt.Printf("Build time:   %s\n", buildTime)
		fmt.Printf("Git commit:   %s\n", commit)
	}
}

func printUsage() {
	fmt.Println(`Usage: gozzi <command> [options]

Commands:
  build    Generate static site
  serve    Start development server
  help     Show help information
  version  Show version information

Use "gozzi help <command>" for more information about a command`)
}

func buildUsage() {
	fmt.Println(`Usage: gozzi build [options]

Options:
  --config string  Path to config file (default "config.toml")
  --content string Content directory (default "content")`)
}

func serveUsage() {
	fmt.Println(`Usage: gozzi serve [options]

Options:
  --config string  Path to config file (default "config.toml")
  --content string Content directory (default "content")
  --port int       Port to listen on (default 1313)`)
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

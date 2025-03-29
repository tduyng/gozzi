package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/generator"
	"github.com/tduyng/gozzi/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	switch os.Args[1] {
	case "build":
		if err := generator.BuildSite(cfg); err != nil {
			log.Fatal(err)
		}
	case "serve":
		startDevServer(cfg)
	default:
		printUsage()
		os.Exit(1)
	}
}

func startDevServer(cfg *config.SiteConfig) {
	port := 1313
	if p := flag.Arg(1); p != "" {
		port, _ = strconv.Atoi(p)
	}

	srv, err := server.NewDevServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	srv.Start(port)
}

func printUsage() {
	fmt.Println(`Usage: gozzi <command>
Commands:
  build    Generate static site
  serve    Start development server with live reload`)
}

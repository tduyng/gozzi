package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/generator"
	"github.com/tduyng/gozzi/internal/server"
)

func main() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	servePort := serveCmd.Int("port", 1313, "Port to serve on")

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
		buildCmd.Parse(os.Args[2:])
		if err := generator.BuildSite(cfg); err != nil {
			log.Fatal(err)
		}

	case "serve":
		serveCmd.Parse(os.Args[2:])
		startDevServer(cfg, *servePort)

	default:
		printUsage()
		os.Exit(1)
	}
}

func startDevServer(cfg *config.GlobalConfig, port int) {
	// Initial build
	if err := generator.BuildSite(cfg); err != nil {
		log.Fatal(err)
	}

	// Setup live reload server
	lrs, err := server.NewLiveReloadServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create live reload server: %v", err)
	}

	// Graceful shutdown handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel() // Always clean up the shutdown context

		lrs.Shutdown(shutdownCtx)

		cancel()
	}()

	// Start server
	if err := lrs.Start(port); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}

func printUsage() {
	fmt.Println(`Usage: gozzi <command>
Commands:
  build    Generate static site
  serve    Start development server with live reload`)
}

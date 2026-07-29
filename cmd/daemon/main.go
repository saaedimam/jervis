package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ioriimasu/jervis/internal/runtime/eventbus"
	permengine "github.com/ioriimasu/jervis/internal/runtime/permissions/engine"
	permregistry "github.com/ioriimasu/jervis/internal/runtime/permissions/registry"
)

func main() {
	fmt.Println("Starting Jervis Background Daemon Service...")

	// Initialize Runtime Core Engines
	eb := eventbus.New()
	pReg := permregistry.New()
	pEngine := permengine.New(pReg)

	if eb == nil || pEngine == nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize runtime engines\n")
		os.Exit(1)
	}

	fmt.Println("Jervis Runtime Daemon active. Press Ctrl+C to terminate.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal or exit if running non-interactively in tests
	if len(os.Args) > 1 && os.Args[1] == "--oneshot" {
		fmt.Println("Daemon oneshot complete.")
		return
	}

	<-sigChan
	fmt.Println("Daemon gracefully shutting down.")
}

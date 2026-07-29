package main

import (
	"fmt"
	"os"

	"github.com/ioriimasu/jervis/internal/runtime/buildinfo"
)

func main() {
	info := buildinfo.Get()
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("Jervis OS v%s (commit: %s, built: %s)\n", info.SemVer(), info.GitCommit(), info.BuildDate())
		return
	}

	fmt.Printf("Jervis Personal OS v%s — Production Ready Runtime Engine\n", info.SemVer())
	fmt.Println("Usage:")
	fmt.Println("  jervis version         Show build information")
	fmt.Println("  jervis daemon          Start runtime background daemon")
	fmt.Println("  jervis mcp             Start Model Context Protocol server")
}

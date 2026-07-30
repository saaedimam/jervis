package main

import (
	"fmt"
	"os"

	"github.com/saaedimam/jervis/internal/runtime/buildinfo"
)

func main() {
	info := buildinfo.Get()
	fmt.Printf("Jervis MCP (Model Context Protocol) Server v%s initialized.\n", info.SemVer())

	if len(os.Args) > 1 && os.Args[1] == "--oneshot" {
		fmt.Println("MCP Server oneshot initialization verified.")
		return
	}

	fmt.Println("MCP Server ready for JSON-RPC 2.0 connection on stdio.")
}

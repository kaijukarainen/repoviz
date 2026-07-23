// Command repoviz-engine is the Repoviz sidecar. It serves the engine's
// JSON-RPC API over stdio and exits when the parent (the VSCode extension)
// closes the stream. All logging goes to stderr; stdout carries only RPC.
package main

import (
	"context"
	"log"
	"os"

	"github.com/kaijukarainen/repoviz/go/repoviz-engine/internal/rpc"
)

// version is overridden at build time via -ldflags "-X main.version=...".
var version = "0.0.0-dev"

type stdio struct{}

func (stdio) Read(p []byte) (int, error)  { return os.Stdin.Read(p) }
func (stdio) Write(p []byte) (int, error) { return os.Stdout.Write(p) }
func (stdio) Close() error                { return os.Stdin.Close() }

func main() {
	log.SetOutput(os.Stderr)
	log.SetPrefix("repoviz-engine ")
	log.SetFlags(0)

	log.Printf("starting version=%s", version)
	handler := rpc.NewHandler("repoviz-engine", version)
	if err := rpc.Serve(context.Background(), stdio{}, handler); err != nil {
		log.Printf("serve ended: %v", err)
		os.Exit(1)
	}
	log.Printf("client disconnected, shutting down")
}

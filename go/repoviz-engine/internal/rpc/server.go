// Package rpc serves the engine's JSON-RPC API over a stream, using the
// VSCode/LSP Content-Length framing so the TypeScript shell's vscode-jsonrpc
// client can talk to it directly.
package rpc

import (
	"context"
	"io"

	"github.com/sourcegraph/jsonrpc2"
)

// Serve runs a JSON-RPC connection over rwc until the peer disconnects or ctx
// is cancelled. The caller owns rwc's lifetime; typically it wraps stdio.
func Serve(ctx context.Context, rwc io.ReadWriteCloser, h *Handler) error {
	stream := jsonrpc2.NewBufferedStream(rwc, jsonrpc2.VSCodeObjectCodec{})
	conn := jsonrpc2.NewConn(ctx, stream, jsonrpc2.HandlerWithError(h.Handle))
	select {
	case <-ctx.Done():
		conn.Close()
		return ctx.Err()
	case <-conn.DisconnectNotify():
		return nil
	}
}

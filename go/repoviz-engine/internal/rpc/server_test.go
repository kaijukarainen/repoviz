package rpc

import (
	"context"
	"net"
	"testing"

	"github.com/kaijukarainen/repoviz/go/repoviz-engine/internal/contract"
	"github.com/sourcegraph/jsonrpc2"
)

// connectLoopback wires a client and server jsonrpc2.Conn over an in-memory
// pipe using the same VSCode framing the real transport uses.
func connectLoopback(t *testing.T, h *Handler) (*jsonrpc2.Conn, func()) {
	t.Helper()
	ctx := context.Background()
	serverEnd, clientEnd := net.Pipe()

	server := jsonrpc2.NewConn(ctx,
		jsonrpc2.NewBufferedStream(serverEnd, jsonrpc2.VSCodeObjectCodec{}),
		jsonrpc2.HandlerWithError(h.Handle),
	)
	noop := jsonrpc2.HandlerWithError(func(context.Context, *jsonrpc2.Conn, *jsonrpc2.Request) (interface{}, error) {
		return nil, nil
	})
	client := jsonrpc2.NewConn(ctx,
		jsonrpc2.NewBufferedStream(clientEnd, jsonrpc2.VSCodeObjectCodec{}),
		noop,
	)
	return client, func() {
		client.Close()
		server.Close()
	}
}

func TestHandshakeRoundTrip(t *testing.T) {
	client, done := connectLoopback(t, NewHandler("repoviz-engine", "1.2.3"))
	defer done()

	var res contract.HandshakeResult
	params := contract.HandshakeParams{ProtocolVersion: ProtocolVersion, ClientName: "test", ClientVersion: "0.0.0"}
	if err := client.Call(context.Background(), "handshake", params, &res); err != nil {
		t.Fatalf("handshake call: %v", err)
	}
	if res.ProtocolVersion != ProtocolVersion {
		t.Errorf("protocolVersion = %d, want %d", res.ProtocolVersion, ProtocolVersion)
	}
	if res.EngineName != "repoviz-engine" || res.EngineVersion != "1.2.3" {
		t.Errorf("engine identity = %q/%q, want repoviz-engine/1.2.3", res.EngineName, res.EngineVersion)
	}
}

func TestHealth(t *testing.T) {
	client, done := connectLoopback(t, NewHandler("repoviz-engine", "1.2.3"))
	defer done()

	var res contract.HealthResult
	if err := client.Call(context.Background(), "health", nil, &res); err != nil {
		t.Fatalf("health call: %v", err)
	}
	if !res.Ok || res.EngineVersion != "1.2.3" {
		t.Errorf("health = %+v, want ok/1.2.3", res)
	}
}

func TestUnknownMethod(t *testing.T) {
	client, done := connectLoopback(t, NewHandler("repoviz-engine", "1.2.3"))
	defer done()

	var res map[string]any
	err := client.Call(context.Background(), "nope", nil, &res)
	if err == nil {
		t.Fatal("expected error for unknown method")
	}
}

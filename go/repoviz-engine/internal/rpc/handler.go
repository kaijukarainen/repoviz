package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kaijukarainen/repoviz/go/repoviz-engine/internal/contract"
	"github.com/sourcegraph/jsonrpc2"
)

// ProtocolVersion is the wire-protocol version this engine speaks. The client
// sends its own in the handshake; a mismatch is the shell's cue to refuse.
const ProtocolVersion = 1

// Handler dispatches JSON-RPC methods to the engine's capabilities.
type Handler struct {
	engineName    string
	engineVersion string
}

// NewHandler builds a Handler that identifies itself with the given name and
// version in handshake and health responses.
func NewHandler(engineName, engineVersion string) *Handler {
	return &Handler{engineName: engineName, engineVersion: engineVersion}
}

// Handle routes a single request and returns its result or an error. It is
// registered with jsonrpc2.HandlerWithError.
func (h *Handler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	switch req.Method {
	case "handshake":
		return h.handshake(req)
	case "health":
		return h.health()
	default:
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: fmt.Sprintf("unknown method %q", req.Method)}
	}
}

func (h *Handler) handshake(req *jsonrpc2.Request) (contract.HandshakeResult, error) {
	var params contract.HandshakeParams
	if req.Params != nil {
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return contract.HandshakeResult{}, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()}
		}
	}
	return contract.HandshakeResult{
		ProtocolVersion: ProtocolVersion,
		EngineName:      h.engineName,
		EngineVersion:   h.engineVersion,
		Capabilities:    []string{},
	}, nil
}

func (h *Handler) health() (contract.HealthResult, error) {
	return contract.HealthResult{Ok: true, EngineVersion: h.engineVersion}, nil
}

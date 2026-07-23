// Code generated from contract/schemas by scripts/codegen.sh. DO NOT EDIT.
package contract

// Root wrapper so codegen emits every message type. Not sent on the wire.
type ProtocolGo struct {
	HandshakeParams *HandshakeParams `json:"handshakeParams,omitempty"`
	HandshakeResult *HandshakeResult `json:"handshakeResult,omitempty"`
	HealthResult    *HealthResult    `json:"healthResult,omitempty"`
}

type HandshakeParams struct {
	ClientName    string `json:"clientName"`
	ClientVersion string `json:"clientVersion"`
	// Protocol version the client speaks.
	ProtocolVersion int64 `json:"protocolVersion"`
}

type HandshakeResult struct {
	// Named capabilities the engine supports.
	Capabilities  []string `json:"capabilities"`
	EngineName    string   `json:"engineName"`
	EngineVersion string   `json:"engineVersion"`
	// Protocol version the engine speaks.
	ProtocolVersion int64 `json:"protocolVersion"`
}

type HealthResult struct {
	EngineVersion string `json:"engineVersion"`
	Ok            bool   `json:"ok"`
}

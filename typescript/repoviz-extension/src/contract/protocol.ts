// Code generated from contract/schemas by scripts/codegen.sh. DO NOT EDIT.

/**
 * Root wrapper so codegen emits every message type. Not sent on the wire.
 */
export interface ProtocolTs {
    handshakeParams?: HandshakeParams;
    handshakeResult?: HandshakeResult;
    healthResult?:    HealthResult;
    [property: string]: any;
}

export interface HandshakeParams {
    clientName:    string;
    clientVersion: string;
    /**
     * Protocol version the client speaks.
     */
    protocolVersion: number;
}

export interface HandshakeResult {
    /**
     * Named capabilities the engine supports.
     */
    capabilities:  string[];
    engineName:    string;
    engineVersion: string;
    /**
     * Protocol version the engine speaks.
     */
    protocolVersion: number;
}

export interface HealthResult {
    engineVersion: string;
    ok:            boolean;
}

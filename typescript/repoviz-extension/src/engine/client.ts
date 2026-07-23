import { spawn, ChildProcessWithoutNullStreams } from "child_process";
import {
  createMessageConnection,
  MessageConnection,
  StreamMessageReader,
  StreamMessageWriter,
} from "vscode-jsonrpc/node";
import { HandshakeParams, HandshakeResult, HealthResult } from "../contract/protocol";

/** Protocol version the shell speaks; must match the engine's. */
export const PROTOCOL_VERSION = 1;

/**
 * EngineClient spawns the Repoviz sidecar and exposes its JSON-RPC methods as
 * typed calls. Engine stderr is forwarded to the supplied log sink; stdout is
 * reserved for the RPC stream.
 */
export class EngineClient {
  private child?: ChildProcessWithoutNullStreams;
  private connection?: MessageConnection;

  constructor(
    private readonly enginePath: string,
    private readonly log: (line: string) => void,
  ) {}

  /** Spawns the engine and begins listening on the connection. */
  start(): void {
    const child = spawn(this.enginePath, [], { stdio: ["pipe", "pipe", "pipe"] });
    child.stderr.on("data", (chunk: Buffer) => this.log(chunk.toString().trimEnd()));
    child.on("exit", (code) => this.log(`engine exited with code ${code ?? "null"}`));

    const connection = createMessageConnection(
      new StreamMessageReader(child.stdout),
      new StreamMessageWriter(child.stdin),
    );
    connection.onClose(() => this.log("engine connection closed"));
    connection.listen();

    this.child = child;
    this.connection = connection;
  }

  /** Performs the handshake and returns the engine's identity and protocol. */
  handshake(params: HandshakeParams): Promise<HandshakeResult> {
    return this.require().sendRequest<HandshakeResult>("handshake", params);
  }

  /** Returns the engine's liveness. */
  health(): Promise<HealthResult> {
    return this.require().sendRequest<HealthResult>("health");
  }

  /** Stops listening and terminates the engine process. */
  dispose(): void {
    this.connection?.dispose();
    this.child?.kill();
    this.connection = undefined;
    this.child = undefined;
  }

  private require(): MessageConnection {
    if (!this.connection) {
      throw new Error("EngineClient not started");
    }
    return this.connection;
  }
}

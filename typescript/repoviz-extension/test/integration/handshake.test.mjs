import { test } from "node:test";
import assert from "node:assert/strict";
import { spawn } from "node:child_process";
import * as path from "node:path";
import { fileURLToPath } from "node:url";
import { createRequire } from "node:module";

// vscode-jsonrpc is CommonJS; load its /node subpath via the CJS resolver.
const require = createRequire(import.meta.url);
const { createMessageConnection, StreamMessageReader, StreamMessageWriter } =
  require("vscode-jsonrpc/node");

const here = path.dirname(fileURLToPath(import.meta.url));
const binary =
  process.env.REPOVIZ_ENGINE_PATH ??
  path.resolve(here, "..", "..", "..", "..", "go", "repoviz-engine", "bin", "repoviz-engine");

// Spawns the real engine binary and exercises the JSON-RPC seam end-to-end
// over stdio, exactly as the extension does at runtime.
test("handshake and health against the real engine binary", async () => {
  const child = spawn(binary, [], { stdio: ["pipe", "pipe", "pipe"] });
  child.stderr.on("data", (c) => process.stderr.write(`[engine] ${c}`));

  const connection = createMessageConnection(
    new StreamMessageReader(child.stdout),
    new StreamMessageWriter(child.stdin),
  );
  connection.listen();

  try {
    const handshake = await connection.sendRequest("handshake", {
      protocolVersion: 1,
      clientName: "integration-test",
      clientVersion: "0.0.0",
    });
    assert.equal(handshake.protocolVersion, 1);
    assert.equal(handshake.engineName, "repoviz-engine");
    assert.ok(typeof handshake.engineVersion === "string");
    assert.ok(Array.isArray(handshake.capabilities));

    const health = await connection.sendRequest("health");
    assert.equal(health.ok, true);
  } finally {
    connection.dispose();
    child.kill();
  }
});

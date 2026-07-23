import * as vscode from "vscode";
import { EngineClient, PROTOCOL_VERSION } from "./engine/client";
import { resolveEnginePath } from "./engine/resolve";

let client: EngineClient | undefined;

/** Activates the extension: spawns the engine and handshakes with it. */
export async function activate(context: vscode.ExtensionContext): Promise<void> {
  const output = vscode.window.createOutputChannel("Repoviz");
  context.subscriptions.push(output);

  const engine = new EngineClient(resolveEnginePath(context.extensionPath), (line) =>
    output.appendLine(line),
  );
  engine.start();
  client = engine;
  context.subscriptions.push({ dispose: () => engine.dispose() });

  try {
    const result = await engine.handshake({
      protocolVersion: PROTOCOL_VERSION,
      clientName: "repoviz-extension",
      clientVersion: context.extension.packageJSON.version ?? "0.0.0",
    });
    if (result.protocolVersion !== PROTOCOL_VERSION) {
      output.appendLine(
        `protocol mismatch: engine speaks ${result.protocolVersion}, shell speaks ${PROTOCOL_VERSION}`,
      );
      void vscode.window.showErrorMessage("Repoviz: engine protocol version mismatch.");
    } else {
      output.appendLine(`connected to ${result.engineName} ${result.engineVersion}`);
    }
  } catch (err) {
    output.appendLine(`handshake failed: ${err instanceof Error ? err.message : String(err)}`);
    void vscode.window.showErrorMessage("Repoviz: failed to reach the engine. See the Repoviz output channel.");
  }

  context.subscriptions.push(
    vscode.commands.registerCommand("repoviz.pingEngine", async () => {
      try {
        const health = await engine.health();
        void vscode.window.showInformationMessage(
          `Repoviz engine ${health.engineVersion}: ${health.ok ? "OK" : "not OK"}`,
        );
      } catch (err) {
        void vscode.window.showErrorMessage(
          `Repoviz ping failed: ${err instanceof Error ? err.message : String(err)}`,
        );
      }
    }),
  );
}

/** Deactivates the extension and tears down the engine. */
export function deactivate(): void {
  client?.dispose();
  client = undefined;
}

import * as path from "path";

/**
 * Resolves the filesystem path to the engine binary. A dev override via the
 * REPOVIZ_ENGINE_PATH environment variable wins; otherwise the path is derived
 * from the extension location. The bundled (packaged) location will slot in
 * here in a later milestone without changing any caller.
 */
export function resolveEnginePath(extensionPath: string): string {
  const override = process.env.REPOVIZ_ENGINE_PATH;
  if (override) {
    return override;
  }
  const binary = process.platform === "win32" ? "repoviz-engine.exe" : "repoviz-engine";
  return path.join(extensionPath, "..", "..", "go", "repoviz-engine", "bin", binary);
}

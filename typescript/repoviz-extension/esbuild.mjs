import { build } from "esbuild";

// Bundles the extension host entry into a single CommonJS file. The `vscode`
// module is provided by the runtime, so it stays external.
await build({
  entryPoints: ["src/extension.ts"],
  bundle: true,
  outfile: "dist/extension.js",
  external: ["vscode"],
  platform: "node",
  target: "node20",
  format: "cjs",
  sourcemap: true,
});

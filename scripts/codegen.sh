#!/usr/bin/env bash
# Regenerates the cross-language contract types from the JSON Schemas in
# /contract. Single source of truth: edit the schema, never the generated code.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SCHEMA="$ROOT/contract/schemas/protocol.json"
GO_OUT="$ROOT/go/repoviz-engine/internal/contract/protocol.go"
TS_OUT="$ROOT/typescript/repoviz-extension/src/contract/protocol.ts"
HEADER="Code generated from contract/schemas by scripts/codegen.sh. DO NOT EDIT."

QUICKTYPE="$ROOT/typescript/repoviz-extension/node_modules/.bin/quicktype"
if [ ! -x "$QUICKTYPE" ]; then
  QUICKTYPE="npx --yes quicktype"
fi

mkdir -p "$(dirname "$GO_OUT")" "$(dirname "$TS_OUT")"

# Go: quicktype --just-types omits the package clause, so we prepend it.
$QUICKTYPE --src-lang schema --lang go --just-types --package contract \
  --out "$GO_OUT.tmp" "$SCHEMA"
{
  echo "// $HEADER"
  echo "package contract"
  echo
  cat "$GO_OUT.tmp"
} >"$GO_OUT"
rm -f "$GO_OUT.tmp"
gofmt -w "$GO_OUT"

# TypeScript.
$QUICKTYPE --src-lang schema --lang typescript --just-types \
  --out "$TS_OUT.tmp" "$SCHEMA"
{
  echo "// $HEADER"
  echo
  cat "$TS_OUT.tmp"
} >"$TS_OUT"
rm -f "$TS_OUT.tmp"

echo "codegen: wrote $GO_OUT and $TS_OUT"

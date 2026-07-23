#!/usr/bin/env bash
# The single CI entrypoint. .github/workflows calls this so pipeline logic is
# versioned here and runnable locally.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

echo "==> install extension deps"
( cd typescript/repoviz-extension && npm ci )

echo "==> codegen freshness"
bash scripts/codegen.sh
if ! git diff --quiet -- go/repoviz-engine/internal/contract typescript/repoviz-extension/src/contract; then
  echo "generated contract types are stale; run scripts/codegen.sh and commit" >&2
  git --no-pager diff -- go/repoviz-engine/internal/contract typescript/repoviz-extension/src/contract
  exit 1
fi

echo "==> engine: vet, test, build"
( cd go/repoviz-engine && go vet ./... && go test ./... && mkdir -p bin && go build -o bin/repoviz-engine ./cmd )

echo "==> extension: typecheck, lint, build"
( cd typescript/repoviz-extension && npm run typecheck && npm run lint && npm run build )

echo "==> integration test"
export REPOVIZ_ENGINE_PATH="$ROOT/go/repoviz-engine/bin/repoviz-engine"
( cd typescript/repoviz-extension && npm run test:integration )

echo "CI OK"

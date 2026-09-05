#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GEN_DIR="${ROOT_DIR}/gen"

mkdir -p "${GEN_DIR}"

echo "Generating Go code"
protoc \
  -I "${ROOT_DIR}/proto" \
  --go_out="${GEN_DIR}" \
  --go_opt=paths=source_relative \
  --go-grpc_out="${GEN_DIR}" \
  --go-grpc_opt=paths=source_relative \
  "${ROOT_DIR}/proto/git/v0/git.proto"

echo "Done."


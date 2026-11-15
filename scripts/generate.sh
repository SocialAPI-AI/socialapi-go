#!/usr/bin/env bash
# Regenerates socialapi.gen.go from openapi.json. Every other .go file is hand-written.
set -euo pipefail
cd "$(dirname "$0")/.."
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.8.2 -config oapi-codegen.yaml openapi.json
go mod tidy

#!/usr/bin/env bash
# Run all unit tests. Some extractor tests require network access.
set -uo pipefail

# The default go build cache may not be writable in sandboxed environments.
export GOCACHE="${GOCACHE:-/tmp/lux-gocache}"
mkdir -p "$GOCACHE"

failed=0

echo "==> go build ./..."
go build ./... || failed=1

echo "==> go vet ./..."
go vet ./... || failed=1

echo "==> go test -race (concurrency-critical packages)"
go test -race -count=1 ./utils/... ./downloader/... ./request/... || failed=1

echo "==> go test ./... (full suite, some tests need network)"
go test -count=1 ./... || failed=1

if [ "$failed" -ne 0 ]; then
	echo "TEST FAILED"
	exit 1
fi
echo "ALL TESTS PASSED"

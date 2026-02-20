#!/bin/bash
# Integration test script for torrent-at-home

set -e

echo "Building torrent-at-home..."
go build -o torrent-at-home ./cmd/app

echo "Running unit tests..."
go test ./...

echo "All tests passed!"

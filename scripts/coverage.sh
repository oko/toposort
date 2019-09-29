#!/bin/bash
set -eux
TMP="$(mktemp -d)"

go test -coverprofile="$TMP/coverage.out"
go tool cover -func="$TMP/coverage.out"

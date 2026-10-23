#!/usr/bin/env bash
set -euo pipefail

binary="${1:-./bin/config-analyzer}"

expect_status() {
  local expected="$1"
  shift
  local actual=0
  "$@" || actual=$?
  if [[ "$actual" -ne "$expected" ]]; then
    printf 'Expected exit code %s, got %s: %s\n' "$expected" "$actual" "$*" >&2
    return 1
  fi
}

expect_status 0 "$binary" --help
expect_status 0 "$binary" examples/safe.yaml
expect_status 1 "$binary" examples/debug.yaml
expect_status 1 "$binary" examples/insecure.yaml
expect_status 0 "$binary" examples/insecure.yaml --silent
expect_status 0 "$binary" examples/insecure.yaml -s
expect_status 1 "$binary" --stdin < examples/insecure.yaml
printf '%s\n' '{"log":{"level":"debug"}}' | expect_status 1 "$binary" --stdin
printf '%s\n' '5' | expect_status 0 "$binary" --stdin
printf '%s\n' 'log: [' | expect_status 2 "$binary" --stdin --silent
expect_status 2 "$binary"
printf 'Smoke checks passed.\n'

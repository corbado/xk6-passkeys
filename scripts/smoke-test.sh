#!/usr/bin/env bash
set -euo pipefail

# Run from the repository root after make build. Each run owns its disposable backend.
docker build -t xk6-passkeys-smoke examples/backend
container_id=$(docker run --rm -d -p 127.0.0.1::8080 xk6-passkeys-smoke)
trap 'docker stop "$container_id" >/dev/null' EXIT
address=$(docker port "$container_id" 8080/tcp)
base_url="http://$address"
ready=false
for attempt in {1..30}; do
  if curl --fail --silent "$base_url/ping" >/dev/null; then
    ready=true
    break
  fi
  sleep 1
done
if [ "$ready" != true ]; then
  docker logs "$container_id"
  exit 1
fi

./k6 run -e BASE_URL="$base_url" -e DURATION=5s examples/registration.js
./k6 run -e BASE_URL="$base_url" -e DURATION=5s examples/login.js
./k6 run -e BASE_URL="$base_url" -e DURATION=5s -e RATE=20 \
  -e PRE_ALLOCATED_VUS=10 -e MAX_VUS=20 examples/login.js
./k6 run -e BASE_URL="$base_url" -e DURATION=3s -e RATE=5 \
  -e PRE_ALLOCATED_VUS=5 -e THINK_TIME=0.1 examples/login.js

# Verify that CI fails for protocol errors, latency violations and undelivered load.
expect_threshold_failure() {
  local result=0
  ./k6 run -e BASE_URL="$base_url" "$@" || result=$?
  if [ "$result" -ne 99 ]; then
    echo "Expected k6 threshold failure (99), got $result" >&2
    exit 1
  fi
}

expect_threshold_failure -e DURATION=1s -e RATE=2 \
  -e ORIGIN=https://invalid.example examples/registration.js
expect_threshold_failure -e DURATION=1s -e RATE=2 \
  -e THINK_TIME=0.05 -e CEREMONY_P95_MS=1 examples/registration.js
expect_threshold_failure -e DURATION=1s -e RATE=20 \
  -e PRE_ALLOCATED_VUS=1 -e MAX_VUS=1 -e THINK_TIME=0.2 examples/registration.js

echo 'All smoke tests passed, including expected threshold failures.'

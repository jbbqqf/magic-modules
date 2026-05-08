#!/usr/bin/env bash
#
# reproduce.sh — single-script PR reproduction for human reviewers
#
# Issue:   https://github.com/hashicorp/terraform-provider-google/issues/18583
# Branch:  feat/18583-bigquery-routine-aggregate-function
# Service: bigquery
# Test:    TestAccBigQueryRoutine_bigqueryRoutineUdafExample
#
# Run from the root of the magic-modules clone, inside the nix dev shell.
# Falls back to a manual environment if nix is not available, but Go,
# Make and a checkout of terraform-provider-google in $GOPATH are required.
#
set -euo pipefail

step() { printf "\n\033[1;34m▸ %s\033[0m\n" "$1"; }
ok()   { printf "\033[1;32m  ✓ %s\033[0m\n" "$1"; }
fail() { printf "\033[1;31m  ✗ %s\033[0m\n" "$1" >&2; }

#-----------------------------------------------------------------------------
# Pre-flight
#-----------------------------------------------------------------------------
step "Pre-flight"
: "${GOPATH:?GOPATH not set. Run 'nix develop' from the bootstrap kit, or export GOPATH manually.}"
TPG="$GOPATH/src/github.com/hashicorp/terraform-provider-google"
if [ ! -d "$TPG" ]; then
  fail "expected TPG checkout at $TPG"
  exit 1
fi
command -v go >/dev/null || { fail "go not found"; exit 1; }
ok "go: $(go version)"
ok "tpg: $TPG"

#-----------------------------------------------------------------------------
# Step 1 — Regenerate TPG from the current magic-modules tree
#-----------------------------------------------------------------------------
step "Regenerating terraform-provider-google from this branch"
make build OUTPUT_PATH="$TPG" VERSION=ga
ok "regenerated"

#-----------------------------------------------------------------------------
# Step 2 — Build
#-----------------------------------------------------------------------------
step "Building TPG"
( cd "$TPG" && go build ./... )
ok "TPG builds"

#-----------------------------------------------------------------------------
# Step 3 — Unit tests for the touched service
#-----------------------------------------------------------------------------
step "Unit tests on google/services/bigquery"
( cd "$TPG" && go test ./google/services/bigquery/... )
ok "unit tests pass"

#-----------------------------------------------------------------------------
# Step 4 — Acceptance test
#-----------------------------------------------------------------------------
step "Acceptance test: TestAccBigQueryRoutine_bigqueryRoutineUdafExample"
if [ -z "TestAccBigQueryRoutine_bigqueryRoutineUdafExample" ]; then
  ok "no acceptance test pattern declared — skipping"
  exit 0
fi

if [ -n "${GCP_PROJECT:-}" ]; then
  export GOOGLE_PROJECT="$GCP_PROJECT"
  ok "live mode: GOOGLE_PROJECT=$GOOGLE_PROJECT"
else
  ok "no GCP_PROJECT set — running in default mode (will use cassettes if present)"
fi

(
  cd "$TPG" && \
  TF_ACC=1 go test -v -timeout 30m \
    -run "TestAccBigQueryRoutine_bigqueryRoutineUdafExample" \
    ./google/services/bigquery/...
)
ok "acceptance test pass"

step "Reproduction complete — PR is reproducible"

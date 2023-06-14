#! /usr/bin/env bash

SCRIPT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)
SCRIPT_CHECK_DOCS_DIFF=$(git diff --name-only --diff-filter=M | grep -Ec '/types\.go$')
GEN_CRD_API_REFERENCE_DOCS="${GEN_CRD_API_REFERENCE_DOCS:-${SCRIPT_ROOT}/.cache/tools/$(go env GOHOSTARCH)/gen-crd-api-reference-docs}"
SPECIFICATION_FILE=Documentation/CRDs/specification.md

run_gen() {
  if [[  "$SCRIPT_CHECK_DOCS_DIFF" -gt 0 || "$1" = "--force" ]]; then
    echo "differences found in types.go, rebuilding specification.md"
    "${GEN_CRD_API_REFERENCE_DOCS}" \
        -config="${SCRIPT_ROOT}/crd-docs-config.json" \
        -template-dir="${SCRIPT_ROOT}/Documentation/gen-crd-api-reference-docs/template" \
        -api-dir="github.com/rook/rook/pkg/apis/ceph.rook.io" \
        -out-file="${SCRIPT_ROOT}/$SPECIFICATION_FILE"
 fi
}

run_gen "$@"

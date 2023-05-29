#! /usr/bin/env bash

SCRIPT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)
GEN_CRD_API_REFERENCE_DOC_VERSION=v0.3.0
SCRIPT_CHECK_DOCS_DIFF=$(git diff --name-only --diff-filter=M | grep -Ec '/types\.go$')

install_generator() {
    go install github.com/ahmetb/gen-crd-api-reference-docs@${GEN_CRD_API_REFERENCE_DOC_VERSION}
}

run_gen() {
  if [[  "$SCRIPT_CHECK_DOCS_DIFF" -gt 0 || "$1" = "--force" ]]; then
    echo "differences found in types.go, rebuilding specification.md"
    install_generator
    gen-crd-api-reference-docs \
        -config="${SCRIPT_ROOT}/crd-docs-config.json" \
        -template-dir="${SCRIPT_ROOT}/Documentation/gen-crd-api-reference-docs/template" \
        -api-dir="github.com/rook/rook/pkg/apis/ceph.rook.io" \
        -out-file="${SCRIPT_ROOT}/Documentation/CRDs/specification.md"
 fi
}

run_gen "$@"
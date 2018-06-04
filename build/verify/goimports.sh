#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOIMPORTS="goimports"
bad_files=$(${GOIMPORTS} ${FLAGS:--l -local github.com/supergiant/supergiant} ./pkg ./cmd/ ./test)

if [[ -n "${bad_files}" ]]; then
  echo "!!! '$GOIMPORTS' needs to be run on the following files: "
  echo "${bad_files}"
  exit 1
fi


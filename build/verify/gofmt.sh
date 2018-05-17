#!/bin/env bash
# https://github.com/kubernetes/repo-infra/blob/master/verify/go-tools/verify-gofmt.sh

set -o errexit
set -o nounset
set -o pipefail

find_files() {
  find . -not \( \
      \( \
        -wholename '*/vendor/*' \
      \) -prune \
    \) -name '*.go'
}

GOFMT="gofmt -s"
GOFMT_CMD="$GOFMT ${FLAGS:--l}"
bad_files=$(find_files | xargs $GOFMT_CMD)
if [[ -n "${bad_files}" ]]; then
  echo "!!! '$GOFMT' needs to be run on the following files: "
  echo "${bad_files}"
  exit 1
fi


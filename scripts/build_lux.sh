#!/usr/bin/bash

set -euo pipefail

print_usage() {
  printf "Usage: build_lux [OPTIONS]

  Build node

  Options:

    -r  Build with race detector
"
}

race=''
while getopts 'r' flag; do
  case "${flag}" in
    r) race='-race' ;;
    *) print_usage
      exit 1 ;;
  esac
done

# Changes to the minimum golang version must also be replicated in
# scripts/build_lux.sh (here)
# Dockerfile
# README.md
# go.mod
go_version_minimum="1.20.10"

go_version() {
    go version | sed -nE -e 's/[^0-9.]+([0-9.]+).+/\1/p'
}

version_lt() {
    # Return true if $1 is a lower version than than $2,
    local ver1=$1
    local ver2=$2
    # Reverse sort the versions, if the 1st item != ver1 then ver1 < ver2
    if  [[ $(echo -e -n "$ver1\n$ver2\n" | sort -rV | head -n1) != "$ver1" ]]; then
        return 0
    else
        return 1
    fi
}

if version_lt "$(go_version)" "$go_version_minimum"; then
    echo "luxd requires Go >= $go_version_minimum, Go $(go_version) found." >&2
    exit 1
fi

# luxd root folder
LUX_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd )
# Load the constants
source "$LUX_PATH"/scripts/constants.sh

build_args="$race"
echo "Building luxd..."
go build $build_args -ldflags "-X github.com/skychains/chain/version.GitCommit=$git_commit $static_ld_flags" -o "$node_path" "$LUX_PATH/main/"*.go

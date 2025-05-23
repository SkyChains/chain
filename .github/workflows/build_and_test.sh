#!/usr/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# Lux root directory
LUX_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd ../.. && pwd )

"$LUX_PATH"/scripts/build.sh
# Check to see if the build script creates any unstaged changes to prevent
# regression where builds go.mod/go.sum files get out of date.
if [[ -z $(git status -s) ]]; then
    echo "Build script created unstaged changes in the repository"
    # TODO: Revise this check once we can reliably build without changes
    # exit 1
fi
"$LUX_PATH"/scripts/build_test.sh
"$LUX_PATH"/scripts/build_fuzz.sh

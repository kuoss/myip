#!/bin/bash
cd $(dirname $0)/..

set -xeuo pipefail
sh hack/lint.sh
sh hack/cover.sh
sh hack/licenses.sh
set +x
echo "✔️ OK"

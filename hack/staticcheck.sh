#!/bin/bash
cd $(dirname $0)/..

which staticcheck || go install honnef.co/go/tools/cmd/staticcheck@2023.1.6

staticcheck
if [[ $? != 0 ]]; then
    echo "❌ FAIL"
    exit 1
fi
echo "✔️ OK"

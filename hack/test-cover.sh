#!/bin/bash
MIN_COVER=50.0

cd $(dirname $0)/..
export PS4='[$(basename $0):$LINENO] '
set -x

go test ./... -v -failfast -race -coverprofile /tmp/coverprofile
if [[ $? != 0 ]]; then
    echo "❌ FAIL - test failed"
    exit 1
fi

COVER=$(go tool cover -func /tmp/coverprofile | tail -1 | grep -oP [0-9.]+)
rm -f /tmp/coverprofile
if [[ $COVER < $MIN_COVER ]]; then
    echo "⚠️ WARN - total COVER: ${COVER}% (<${MIN_COVER}%)"
    exit
fi
echo "✔️ OK - total COVER: ${COVER}% (>=${MIN_COVER}%)"

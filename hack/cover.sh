#!/bin/bash
cd $(dirname $0)/../

go test ./... -failfast -coverprofile /tmp/coverprofile
if [[ $? != 0 ]]; then
    echo "❌ FAIL - test failed"
    exit 1
fi

go tool cover -func /tmp/coverprofile
rm -f /tmp/coverprofile

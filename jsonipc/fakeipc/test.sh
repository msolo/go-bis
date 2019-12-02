#!/bin/bash

actual=$(echo '{"Val":"dummy"}' | ./fakeipc)
expected='{"Val":"dummy","Error":null}'
if [[ "$actual" != "$expected" ]]; then
    echo "expected: $expected" >&2
    echo "actual: $actual" >&2
    exit 1
fi

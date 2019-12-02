#!/bin/bash

function check() {
    actual=$1
    expected=$2
    if [[ "$actual" != "$expected" ]]; then
        echo "expected: $expected" >&2
        echo "actual: $actual" >&2
        exit 1
    fi
}


# Go JSON RPC allows a single param that is an object. Legal, but not likely to work with everything.
payload='{ "method": "Echo.Echo", "params": [{"Val":"Hello JSON-RPC"}], "id":1}'
actual=$(echo $payload | ./fakerpc)
expected='{"id":1,"result":{"Val":"Hello JSON-RPC"},"error":null}'
check $actual $expected

# Return order is non-deterministic (this is good)
payload='{ "method": "Echo.Echo", "params": [{"Val":"Message 1"}], "id":1} { "method": "Echo.Echo", "params": [{"Val":"Message 2"}], "id":2}'
actual=$(echo $payload | ./fakerpc | sort)
expected='{"id":1,"result":{"Val":"Message 1"},"error":null}\
{"id":2,"result":{"Val":"Message 2"},"error":null}'
check $actual $expected

# "Notifications" don't seem to work with the Go stdlib. We get a request with id=null and still send a response.
payload='{ "method": "Echo.Echo", "params": [{"Val":"Test notify"}], "id":null}'
actual=$(echo $payload | ./fakerpc)
expected='{"id":null,"result":{"Val":"Test notify"},"error":null}'
check $actual $expected

# Return an error - Go JSON RPC allows only a string error which is legal, but limiting.
payload='{ "method": "Echo.Error", "params": [{"Val":"Batfail!"}], "id":1}'
actual=$(echo $payload | ./fakerpc)
expected='{"id":1,"result":null,"error":"error with msg: Batfail!"}'
check $actual $expected

# Return a rich error - doesn't work - error can only be a string.
payload='{ "method": "Echo.Error2", "params": [{"Val":"Batfail!"}], "id":1}'
actual=$(echo $payload | ./fakerpc)
expected='{"id":1,"result":null,"error":"error with msg"}'
check $actual $expected

# actual=$(echo $payload | ./fakerpc)
# expected='{"Result":"dummy","Error":null}'
# if [[ "$actual" != "$expected" ]]; then
#     echo "expected: $expected" >&2
#     echo "actual: $actual" >&2
#     exit 1
# fi

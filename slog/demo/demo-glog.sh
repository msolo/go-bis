#!/bin/bash

# Print the structured addenda.
./demo -log.file /dev/stdout | awk -F' [|] ' '{print $NF}' | jq .

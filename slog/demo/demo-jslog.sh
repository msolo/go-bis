#!/bin/bash

# Print the standard and special fields.
./demo -log.fmt json -log.file /dev/stdout | jq -r '[.Level,.Timestamp,.Message,.Fields["field-a"]]|@tsv' | awk -F"\t" '{print $0}'

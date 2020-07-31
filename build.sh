#!/bin/bash

set -eux

PROTO_ROOT_DIRS="./pb"

# As this is a proto root, and there may be subdirectories with protos, compile the protos for each sub-directory which contains them
for protos in $(find "$PROTO_ROOT_DIRS" -name '*.proto' -exec dirname {} \; | sort | uniq) ; do
    protoc --proto_path="$PROTO_ROOT_DIRS" \
            --go_out=plugins=grpc:"$PROTO_ROOT_DIRS" \
            "$protos"/*.proto
done

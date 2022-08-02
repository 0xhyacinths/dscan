#!/bin/bash

mkdir -p ./webui/apispec
mkdir -p ./server/proto
mkdir -p ./webui/src/lib/proto

# you may need to install some stuff
# cd into the server dir, then run the command in the comment in tools.go.

# generate protos
protoc \
  -I ./proto \
  --go_out=paths=source_relative:./server/proto \
  --go-grpc_out=paths=source_relative:./server/proto \
  --grpc-gateway_out=paths=source_relative:./server/proto \
  --openapiv2_out=./webui/apispec \
  ./proto/descan.proto

pushd webui
npx openapi-typescript ./apispec/descan.swagger.json \
  --output ./src/lib/proto/descan.ts
popd

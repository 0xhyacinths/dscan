#!/bin/bash

mkdir -p ./webui/apispec
mkdir -p ./webui/src/lib/proto

# require go
# install protoc generator
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.11.0

# require protoc
protoc \
  -I ./proto \
  --openapiv2_out=./webui/apispec \
  ./proto/descan.proto

pushd webui

# require node
npm install

# require node
npx openapi-typescript ./apispec/descan.swagger.json \
  --output ./src/lib/proto/descan.ts

# require node
OUTPUT_BASE="/dscan" npm run build

popd

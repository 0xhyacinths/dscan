#!/bin/bash

mkdir -p ./apispec
mkdir -p ./src/lib/proto

# require go1.18.1
# install protoc generator
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.11.0

# require protoc
protoc \
  -I ./proto \
  --openapiv2_out=./apispec \
  ./proto/descan.proto

# require node
npx openapi-typescript ./apispec/descan.swagger.json \
  --output ./src/lib/proto/descan.ts

# require node
OUTPUT_BASE="/dscan" npm run build

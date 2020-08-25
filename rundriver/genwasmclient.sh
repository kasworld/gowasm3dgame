#!/usr/bin/env bash

BUILD_VER=${1}

echo "GOOS=js GOARCH=wasm go build -o clientdata/wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go"
GOOS=js GOARCH=wasm go build -o clientdata/wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go

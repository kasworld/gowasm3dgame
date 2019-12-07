#!/usr/bin/env bash


rm wasmclient.wasm

echo "GOOS=js GOARCH=wasm go build -o wasmclient.wasm"
GOOS=js GOARCH=wasm go build -o wasmclient.wasm wasmclient.go

echo "move files"
mv wasmclient.wasm ./www/

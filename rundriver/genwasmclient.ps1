
$BUILD_VER=$args[0]

$env:GOOS="js" 
$env:GOARCH="wasm" 
Write-Output "GOOS=js GOARCH=wasm go build -o clientdata/wasmclient.wasm -ldflags -X main.Ver=${BUILD_VER}"
go build -o wasmclient.wasm -ldflags "-X main.Ver=${BUILD_VER}" wasmclient.go
$env:GOOS=""
$env:GOARCH=""


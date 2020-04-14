#!/usr/bin/env bash


DATESTR=`date -Iseconds`
GITSTR=`git rev-parse HEAD`
BUILD_VER=${DATESTR}_${GITSTR}_release
echo "Build" ${BUILD_VER}


BuildBin() {
    local srcfile=${1}
    local dstdir=${2}
    local dstfile=${3}
    local args="-X main.Ver=${BUILD_VER}"

    echo "[BuildBin] go build -i -o ${dstdir}/${dstfile} -ldflags "${args}" ${srcfile}"

    mkdir -p ${dstdir}
    go build -i -o ${dstdir}/${dstfile} -ldflags "${args}" ${srcfile}

    if [ ! -f "${dstdir}/${dstfile}" ]; then
        echo "${dstdir}/${dstfile} build fail, build file: ${srcfile}"
        # exit 1
    fi
    strip "${dstdir}/${dstfile}"
}


cd lib
genlog -leveldatafile ./w3dlog/w3dlog.data -packagename w3dlog 
cd ..


ProtocolW3DFiles="protocol_w3d/w3d_gendata/command.data \
protocol_w3d/w3d_gendata/error.data \
protocol_w3d/w3d_gendata/noti.data \
"

PROTOCOL_W3D_VERSION=`cat ${ProtocolW3DFiles}| sha256sum | awk '{print $1}'`
echo "Protocol W3D Version:" ${PROTOCOL_W3D_VERSION}

cd protocol_w3d
genprotocol -ver=${PROTOCOL_W3D_VERSION} \
    -basedir=. \
    -prefix=w3d -statstype=int

goimports -w .

cd ..


genenum -typename=ActType -packagename=acttype -basedir=enum -statstype=int
goimports -w enum/acttype/acttype_gen.go
goimports -w enum/acttype_stats/acttype_stats_gen.go

genenum -typename=GameObjType -packagename=gameobjtype -basedir=enum -statstype=int
goimports -w enum/gameobjtype/gameobjtype_gen.go
goimports -w enum/gameobjtype_stats/gameobjtype_stats_gen.go



GameDataFiles="
game/gameconst/gameconst.go \
game/gameconst/serviceconst.go \
game/gamedata/*.go \
enum/*.enum \
"
Data_VERSION=`cat ${GameDataFiles}| sha256sum | awk '{print $1}'`
echo "Data Version:" ${Data_VERSION}

echo "
package gameconst

const DataVersion = \"${Data_VERSION}\"
" > game/gameconst/dataversion_gen.go 

# build bin

BIN_DIR="bin"
SRC_DIR="rundriver"

echo ${BUILD_VER} > ${BIN_DIR}/BUILD

BuildBin ${SRC_DIR}/server.go ${BIN_DIR} server
BuildBin ${SRC_DIR}/multiclient.go ${BIN_DIR} multiclient

cd rundriver
echo "build wasm client"
GOOS=js GOARCH=wasm go build -o www/wasmclient.wasm wasmclient.go
cd ..
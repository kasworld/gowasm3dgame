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


ProtocolW3DFiles="protocol_w3d/*.enum"

PROTOCOL_W3D_VERSION=`cat ${ProtocolW3DFiles}| sha256sum | awk '{print $1}'`
echo "Protocol W3D Version:" ${PROTOCOL_W3D_VERSION}

cd protocol_w3d
genprotocol -ver=${PROTOCOL_W3D_VERSION} \
    -basedir=. \
    -prefix=w3d -statstype=int

goimports -w .

cd ..


genenum -typename=ActType -packagename=acttype -basedir=enum -vectortype=int
genenum -typename=GameObjType -packagename=gameobjtype -basedir=enum -vectortype=int
genenum -typename=StageType -packagename=stagetype -basedir=enum -vectortype=int

cd enum 
goimports -w .
cd ..

GameDataFiles="
config/gameconst/gameconst.go \
config/gameconst/serviceconst.go \
config/gamedata/*.go \
enum/*.enum \
"
Data_VERSION=`cat ${GameDataFiles}| sha256sum | awk '{print $1}'`
echo "Data Version:" ${Data_VERSION}
mkdir -p config/dataversion
echo "
package dataversion

const DataVersion = \"${Data_VERSION}\"
" > config/dataversion/dataversion_gen.go 

# build bin

BIN_DIR="bin"
SRC_DIR="rundriver"

echo ${BUILD_VER} > ${BIN_DIR}/BUILD

BuildBin ${SRC_DIR}/server.go ${BIN_DIR} server
BuildBin ${SRC_DIR}/multiclient.go ${BIN_DIR} multiclient

cd rundriver
./genwasmclient.sh
cd ..
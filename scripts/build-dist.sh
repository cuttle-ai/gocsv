#!/usr/bin/env bash

DIST_DIR=$(pwd)/dist
SRC_DIR=gocsv
EXECUTABLE=gocsv

GIT_HASH=$(git rev-parse HEAD)
VERSION=$(cat VERSION)
LD_FLAGS="-X main.VERSION=${VERSION} -X main.GIT_HASH=${GIT_HASH}"

rm -rf ${DIST_DIR}
mkdir ${DIST_DIR}
cd ${SRC_DIR}
for os in darwin windows linux; do
  for arch in amd64; do
    basename=gocsv-${os}-${arch}
    mkdir ${DIST_DIR}/${basename}
    if [ "${os}" == "windows" ]; then
      binary="${EXECUTABLE}.exe"
    else
      binary=${EXECUTABLE}
    fi
    env GOOS=${os} GOARCH=${arch} go build -ldflags "${LD_FLAGS}" -o ${DIST_DIR}/${basename}/${binary}
    cd ${DIST_DIR} && zip -rq ${basename}.zip ${basename}
    cd ~-
    rm -r ${DIST_DIR}/${basename}
  done
done
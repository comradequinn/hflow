#! /bin/bash

rm -r ./release 2> /dev/null
mkdir release

function build() {
    echo "building binary for ${GOOS}/${GOARCH}..."
    go build -o "./release/hflow.${GOOS}.${GOARCH}" && echo "successfully built hflow.${GOOS}.${GOARCH}" || echo "error building hflow.${GOOS}.${GOARCH}"
}

ORIG_GOOS=$GOOS
ORIG_GOARCH=$GOARCH

GOOS=linux
GOARCH=amd64
build

GOOS=darwin
GOARCH=amd64
build

GOOS=darwin
GOARCH=arm64
build

GOOS=$ORIG_GOOS
GOARCH=$ORIG_GOARCH

echo "release builds completed"
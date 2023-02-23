#! /bin/bash

rm -r ./release 2> /dev/null
mkdir release

function build() {
    echo "building binary for ${GOOS}/${GOARCH}..."
    go build -o "./release/hflow.${GOOS}.${GOARCH}" && echo "successfully built hflow.${GOOS}.${GOARCH}" || echo "error building hflow.${GOOS}.${GOARCH}"
}

export GOOS=linux
export GOARCH=amd64
build

export GOOS=darwin
export GOARCH=amd64
build

export GOOS=darwin
export GOARCH=arm64
build

echo "release builds completed"

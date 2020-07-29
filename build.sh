#!/bin/bash
set -e

GIT_COMMIT=$(git rev-list -1 HEAD)
GIT_COMMIT_SHORT=${GIT_COMMIT:0:10}
GIT_DIRTY=$(git diff --quiet || echo "-dirty")
BUILD_DATE=$(date +"%Y%m%d")

OS=${1:-darwin}

if [[ $OS == "darwin" || $OS == "macos-latest" ]]; then
    GOOS=darwin
    GOARCH=amd64
elif [[ $OS == "linux" || $OS == "ubuntu-latest" ]]; then
    GOOS=linux
    GOARCH=amd64
elif [[ $OS == "windows" || $OS == "windows-latest" ]]; then
    GOOS=windows
    GOARCH=amd64
else
    echo "Unsupported OS: $OS"
    exit 1
fi

OUTPUT=bin/${GOOS}_${GOARCH}/apnstool
if [[ $GOOS == "windows" ]]; then
    OUTPUT=${OUTPUT}.exe
fi

export CGO_ENABLED=0
export GOOS=$GOOS
export GOARCH=$GOARCH
go build \
    -a -ldflags="-X 'github.com/brannon/apnstool/build_version.CommitHash=$GIT_COMMIT_SHORT$GIT_DIRTY' -X 'github.com/brannon/apnstool/build_version.BuildDate=$BUILD_DATE'" \
    -o $OUTPUT \
    .

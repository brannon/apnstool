#!/bin/sh
set -e

VERSION=$1

echo "* Creating release $VERSION"

echo "* Building"
gox -arch="amd64" -os="darwin linux windows" -output "bin/{{.OS}}_{{.Arch}}/apnstool"

echo "* Packaging"
tar czf bin/apnstool.$VERSION.macos.tar.gz -C bin/darwin_amd64 apnstool
tar czf bin/apnstool.$VERSION.linux.tar.gz -C bin/linux_amd64 apnstool
(cd bin/windows_amd64 && zip -9 ../apnstool.$VERSION.windows.zip apnstool.exe)

echo "* Uploading to release"

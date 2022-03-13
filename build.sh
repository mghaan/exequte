#!/bin/bash

#BUILD=("linux/amd64 linux/arm linux/arm64 windows/amd64 darwin/amd64 darwin/arm64")
BUILD=("linux/amd64")

for platform in ${BUILD[@]}; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    DEST="$GOOS-$GOARCH"


    echo -n "Building for $GOOS/$GOARCH ..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o ./BUILD/$DEST/exequte

    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -buildmode=plugin -o BUILD/$DEST/plugins/dummy.so plugins/dummy/dummy.go
    echo "done"
done

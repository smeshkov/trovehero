#!/bin/sh

DIST_DIR="_dist"
BINARY="trovehero"
OS="$1"
VERSION="$2"

if [ -z "$OS" ]; then
    OS="darwin"
fi

if [ ! -z "$VERSION" ]; then
    VERSION="_$VERSION"
fi

if [ ! -d "$DIST_DIR" ]; then
  mkdir -p $DIST_DIR
fi

env GOOS=${OS} GOARCH=amd64 CGO_ENABLED=0 go build -v -o ${DIST_DIR}/${BINARY}_${OS}${VERSION} .
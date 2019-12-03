#!/bin/sh

sudo apt-get update

sudo apt-get -y install libsdl2-image-dev
sudo apt-get -y install libsdl2-mixer-dev
sudo apt-get -y install libsdl2-ttf-dev
sudo apt-get -y install libsdl2-gfx-dev

go get -v github.com/veandco/go-sdl2/sdl
go get -v github.com/veandco/go-sdl2/img
go get -v github.com/veandco/go-sdl2/mix
go get -v github.com/veandco/go-sdl2/ttf
go get -v github.com/veandco/go-sdl2/gfx

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

env GOOS=${OS} GOARCH=amd64 go build -v -o ${DIST_DIR}/${BINARY}_${OS}${VERSION} cmd/game/game.go
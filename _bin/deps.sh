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
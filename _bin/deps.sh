#!/bin/sh

sudo apt-get update
sudo apt-get -y install libsdl2{,-image,-mixer,-ttf,-gfx}-dev
go get -v github.com/veandco/go-sdl2/{sdl,img,mix,ttf}
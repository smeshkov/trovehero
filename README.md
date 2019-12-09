# Trove Hero

An attempt to build a top down adventure game. 
Inspired by set of videos by [Francesc Campoy](https://campoy.cat/) - [Flappy Gopher](https://github.com/campoy/flappy-gopher).

## Setup

Follow instractions from [here](https://github.com/veandco/go-sdl2) to install SDL2, [pkg-config](https://en.wikipedia.org/wiki/Pkg-config) is also required.

## Playing game

```
go get github.com/smeshkov/trovehero
```
And run the binary generated in `$GOPATH/bin`.

OR

`make run`

Use `arrows` to move arround the green rectangle in order to collect yellow rectangles and avoid red and blue ones, you can use `space` to jump over a blue rectangle.

![Trove Hero](https://storage.googleapis.com/www.zoomio.org/trovehero.png)
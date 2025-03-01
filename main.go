package main

import (
	"asteroids/goasteroids"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {


	ebiten.SetWindowTitle("Asteroids")
	ebiten.SetWindowSize(goasteroids.ScreenWidth, goasteroids.ScreenHeight)

	err := ebiten.RunGame(&goasteroids.Game{})
	if err != nil {
		panic(err)
	}
}

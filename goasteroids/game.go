package goasteroids

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	sceneManager * SceneManager
	input Input
}

type Input struct {}

func (i *Input) Update() {}

func (g *Game) Update() error {
	if g.sceneManager == nil {
		g.sceneManager = &SceneManager{}
		g.sceneManager.GoToScene(&TitleScene{})
	}

	g.input.Update()
	if err := g.sceneManager.Update(& g.input); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneManager.Draw(screen)
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
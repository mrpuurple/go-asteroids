package goasteroids

import (
	"asteroids/assets"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
)

type TitleScene struct {}

func (t *TitleScene) Draw(screen *ebiten.Image) {
	textToDraw := "1 coin = 1 play"
	tw := widthOfText(assets.TitleFont, textToDraw)
	text.Draw(screen, textToDraw, assets.TitleFont, ScreenWidth/2-tw/2, ScreenHeight-200, color.White)
}

func (t *TitleScene) Update(state *State) error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.SceneManager.GoToScene(NewGameScene())
		return nil
	}
	return nil
}

func widthOfText(f font.Face, t string) int {
	_, textWidth := font.BoundString(f, t)
	return textWidth.Round()
}
package goasteroids

import (
	"asteroids/assets"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// LevelStartsScene is the type for our title scene. It holds the game, a timer, and a slice of stars.
type LevelStartsScene struct {
	game           *GameScene
	nextLevelTimer *Timer
	stars          []*Star
}

// Draw puts all the elements on the screen. It's called once per frame.
func (l *LevelStartsScene) Draw(screen *ebiten.Image) {
	// Draw stars.
	for _, s := range l.stars {
		s.Draw(screen)
	}

	// Draw text.
	textToDraw := fmt.Sprintf("LEVEL %d", l.game.currentLevel)
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, ScreenHeight/2)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.TitleFont,
		Size:   48,
	}, op)
}

// Update updates screen elements. It's called once per tick.
func (l *LevelStartsScene) Update(state *State) error {
	// Update timer
	l.nextLevelTimer.Update()
	if l.nextLevelTimer.IsReady() {
		l.game.meteorsForLevel += 2
		l.game.meteorCount = 0
		for k, v := range l.game.lasers {
			delete(l.game.lasers, k)
			l.game.space.Remove(v.laserObj)
		}
		state.SceneManager.GoToScene(l.game)
	}

	// Check to see  if the space bar is pressed. If it is, go to the next scene.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		l.game.meteorsForLevel += 5
		l.game.meteorCount = 0
		for k, v := range l.game.lasers {
			delete(l.game.lasers, k)
			l.game.space.Remove(v.laserObj)
		}
		state.SceneManager.GoToScene(l.game)
		return nil
	}

	return nil
}

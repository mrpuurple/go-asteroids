package goasteroids

import (
	"asteroids/assets"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// TitleScene is the type for our title scene.
type TitleScene struct {
	meteors     map[int]*Meteor // A map of meteors.
	meteorCount int             // How many meteors we currently have in the game.
	stars       []*Star         // A slice of stars.
}

var highScore int
var originalHighScore int

func init() {
	hs, err := getHighScore()
	if err != nil {
		log.Println("Error getting high score:", err)
	}
	highScore = hs
	originalHighScore = hs
}

// Draw draws all elements on the screen. It's called once per frame.
func (t *TitleScene) Draw(screen *ebiten.Image) {
	// Draw stars.
	for _, s := range t.stars {
		s.Draw(screen)
	}

	// Draw meteors.
	for _, m := range t.meteors {
		m.Draw(screen)
	}

	// Draw text.
	textToDraw := "Press space to play"

	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(float64(ScreenWidth/2), ScreenHeight-200)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.TitleFont,
		Size:   48,
	}, op)


}

// Update updates all game scene elements for the next draw. It's called once per tick.
func (t *TitleScene) Update(state *State) error {
	// Check for a spacebar press.
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.SceneManager.GoToScene(NewGameScene())
		return nil
	}

	// Draw meteors, if appropriate.
	if len(t.meteors) < 10 {
		m := NewMeteor(0.25, &GameScene{}, len(t.meteors)-1)
		t.meteorCount++
		t.meteors[t.meteorCount] = m
	}

	// Update meteors.
	for _, m := range t.meteors {
		m.Update()
	}

	return nil
}

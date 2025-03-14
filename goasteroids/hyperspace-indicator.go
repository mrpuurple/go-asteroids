package goasteroids

import (
	"asteroids/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

// HyperspaceIndicator is the type for our indicator.
type HyperspaceIndicator struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

// NewHyperspaceIndicator creates a new hyperspace indicator object.
func NewHyperspaceIndicator(pos Vector) *HyperspaceIndicator {
	return &HyperspaceIndicator{
		position: pos,
		sprite:   assets.HyperspaceIndicator,
	}
}

// Update updates the object. Since our objects are static, it does nothing.
func (hsi *HyperspaceIndicator) Update() {}

// Draw draws the object.
func (hsi *HyperspaceIndicator) Draw(screen *ebiten.Image) {
	bounds := hsi.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	// We want to manage the color, so we'll use colorm.DrawImageOptions instead of
	// the standard ebiten.DrawImageOptions we've been using so far.
	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(halfW, halfH)
	cm := colorm.ColorM{}
	cm.Scale(1.0, 1.0, 1.0, 0.2)
	op.GeoM.Translate(hsi.position.X, hsi.position.Y)
	colorm.DrawImage(screen, hsi.sprite, cm, op)
}

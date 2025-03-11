package goasteroids

import (
	"asteroids/assets"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	laserSpeedPerSecond = 1000.0
)

type Laser struct {
	game     *GameScene
	position Vector
	rotation float64
	sprite   *ebiten.Image
	laserObj *resolv.ConvexPolygon
}

func NewLaser(pos Vector, rotation float64, index int, g *GameScene) *Laser {
	// Set the sprite.
	sprite := assets.LaserSprite

	// Position X & Y coordinates from the center of the sprite.
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	// Create a laser object.
	l := &Laser{
		game:     g,
		position: pos,
		rotation: rotation,
		sprite:   sprite,
		laserObj: resolv.NewRectangle(pos.X, pos.Y, float64(sprite.Bounds().Dx()), float64(sprite.Bounds().Dy())),
	}

	// Set the position of the collision object.
	l.laserObj.SetPosition(pos.X, pos.Y)
	l.laserObj.SetData(&ObjectData{index: index})
	l.laserObj.Tags().Set(TagLaser)

	return l
}

func (l *Laser) Update() {
	// How fast should the laser go.
	speed := laserSpeedPerSecond / float64(ebiten.TPS())
	dx := math.Sin(l.rotation) * speed
	dy := math.Cos(l.rotation) * -speed

	l.position.X += dx
	l.position.Y += dy

	l.laserObj.SetPosition(l.position.X, l.position.Y)
}

func (l *Laser) Draw(screen *ebiten.Image) {
	bounds := l.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(l.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(l.position.X, l.position.Y)

	screen.DrawImage(l.sprite, op)
	
}

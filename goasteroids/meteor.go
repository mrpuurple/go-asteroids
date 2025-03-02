package goasteroids

import (
	"asteroids/assets"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rotationSpeedMin                    = -0.02
	rotationSpeedMax                    = 0.02
	numberOfSmallMeteorsFromLargeMeteor = 4
)

type Meteor struct {
	game          *GameScene
	position      Vector
	rotation      float64
	movement      Vector
	angle         float64
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewMeteor(baseVelocity float64, g *GameScene, index int) *Meteor {
	// Target the center of the screen.
	target := Vector {
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}

	// Pick a random angle.
	angle := rand.Float64() * 2 * math.Pi

	// The distance from the center that meteor should spawn at.
	// half the width, add some arbitrary distance.
	r := ScreenWidth/2.0 + 500

	// Create the position vector, using the angle and simple math.
	pos := Vector{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	// Keep the meteor moving towards the center of the screen.
	// Give it a random velocity.
	velocity := baseVelocity + rand.Float64()*1.5

	// Create the direction vector and normalize it.
	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	// Create the movement vector.
	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	// Assign a sprite to the meteor.
	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]

	// Creata a meteor object and return it.
	m := &Meteor{
		game: g,
		position: pos,
		movement: movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite: sprite,
		angle: angle,
	}

	return m
}

func (m *Meteor) Update() {
	dx := m.movement.X
	dy := m.movement.Y

	m.position.X += dx
	m.position.Y += dy
	m.rotation += m.rotationSpeed

	// Keep meteor on screen.
	m.keepOnScreen()
}

func (m *Meteor) Draw(screen *ebiten.Image) {
	bounds := m.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprite, op)
}

func (m *Meteor) keepOnScreen() {
	if m.position.X >= float64(ScreenWidth) {
		m.position.X = 0
	}
	if m.position.X < 0 {
		m.position.X = ScreenWidth
	}
	if m.position.Y >= float64(ScreenHeight) {
		m.position.Y = 0
	}
	if m.position.Y < 0 {
		m.position.Y = ScreenHeight
	}
}
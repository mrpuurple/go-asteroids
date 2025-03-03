package goasteroids

import (
	"asteroids/assets"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	rotationPerSecond = math.Pi
	maxAcceleration   = 8.0
	ScreenWidth       = 1280
	ScreenHeight      = 720
)

var curAcceleration float64

type Player struct {
	game           *GameScene
	sprite         *ebiten.Image
	rotation       float64
	position       Vector
	playerVelocity float64
	playerObj      *resolv.Circle
}

func NewPlayer(game *GameScene) *Player {
	sprite := assets.PlayerSprite

	// Center player on screen
	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: ScreenWidth/2 - halfW,
		Y: ScreenHeight/2 - halfH,
	}

	// Create a resolv object.
	playerObj := resolv.NewCircle(pos.X, pos.Y, float64(sprite.Bounds().Dx()/2))

	p := &Player{
		sprite:   sprite,
		game:     game,
		position: pos,
		playerObj: playerObj,
	}

	p.playerObj.SetPosition(pos.X, pos.Y)
	p.playerObj.Tags().Set(TagPlayer)

	return p

}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Update() {
	speed := rotationPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	p.accelerate()

	p.playerObj.SetPosition(p.position.X, p.position.Y)
}

func (p *Player) accelerate() {
	p.keepOnScreen()

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if curAcceleration < maxAcceleration {
			curAcceleration = p.playerVelocity + 4
		}

		if curAcceleration >= 8 {
			curAcceleration = 8
		}

		p.playerVelocity = curAcceleration

		// Move in the direction we are pointing.
		dx := math.Sin(p.rotation) * curAcceleration
		dy := math.Cos(p.rotation) * -curAcceleration

		// Move the player on the screen.
		p.position.X += dx
		p.position.Y += dy
	}
}

func (p *Player) keepOnScreen() {
	if p.position.X >= float64(ScreenWidth) {
		p.position.X = 0
		p.playerObj.SetPosition(0, p.position.Y)
	}
	if p.position.X < 0 {
		p.position.X = ScreenWidth
		p.playerObj.SetPosition(ScreenWidth, p.position.Y)
	}
	if p.position.Y >= float64(ScreenHeight) {
		p.position.Y = 0
		p.playerObj.SetPosition(p.position.X, 0)
	}
	if p.position.Y < 0 {
		p.position.Y = ScreenHeight
		p.playerObj.SetPosition(p.position.X, ScreenHeight)
	}
}

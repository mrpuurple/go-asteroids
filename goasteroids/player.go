package goasteroids

import (
	"asteroids/assets"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

const (
	rotationPerSecond     = math.Pi
	maxAcceleration       = 8.0
	ScreenWidth           = 1280 // The width of the screen. We use a 16/9 aspect ratio.
	ScreenHeight          = 720  // The height of the screen.
	shootCoolDown         = time.Millisecond * 150
	burstCoolDown         = time.Millisecond * 500
	laserSpawnOffset      = 50.0
	maxShotsPerBurst      = 3
	dyingAnimationAmmount = 50 * time.Millisecond
)

var curAcceleration float64
var shotsFired = 0

type Player struct {
	game           *GameScene
	sprite         *ebiten.Image
	rotation       float64
	position       Vector
	playerVelocity float64
	playerObj      *resolv.Circle
	shootCoolDown  *Timer
	burstCoolDown  *Timer
	isShielded     bool
	isDying        bool
	isDead         bool
	dyingTimer     *Timer
	dyingCounter   int
	livesRemaining int
}

func NewPlayer(game *GameScene) *Player {
	sprite := assets.PlayerSprite

	// Center player on screen.
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
		sprite:        sprite,
		game:          game,
		position:      pos,
		playerObj:     playerObj,
		shootCoolDown: NewTimer(shootCoolDown),
		burstCoolDown: NewTimer(burstCoolDown),
		isShielded: false,
		isDying: false,
		isDead: false,
		dyingTimer: NewTimer(dyingAnimationAmmount),
		dyingCounter: 0,
		livesRemaining: 1,
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

	p.isPlayerDead()

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	p.accelerate()

	p.isDoneAccelerating()

	p.updateExhaustSprite()

	p.playerObj.SetPosition(p.position.X, p.position.Y)

	p.burstCoolDown.Update()

	p.shootCoolDown.Update()

	p.fireLasers()
}

func (p *Player) isPlayerDead() {
	if p.isDead {
		p.game.playerIsDead = true
	}
}

func (p *Player) fireLasers() {
	if p.burstCoolDown.IsReady() {
		if p.shootCoolDown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
			p.shootCoolDown.Reset()
			shotsFired++
			if shotsFired <= maxShotsPerBurst {
				bounds := p.sprite.Bounds()
				halfW := float64(bounds.Dx()) / 2
				halfH := float64(bounds.Dy()) / 2

				spawnPos := Vector{
					p.position.X + halfW + math.Sin(p.rotation)*laserSpawnOffset,
					p.position.Y + halfH + math.Cos(p.rotation)*-laserSpawnOffset,
				}

				p.game.laserCount++
				laser := NewLaser(spawnPos, p.rotation, p.game.laserCount, p.game)
				p.game.lasers[p.game.laserCount] = laser
				p.game.space.Add(laser.laserObj)
			} else {
				p.burstCoolDown.Reset()
				shotsFired = 0
			}
		}
	}
}

func (p *Player) accelerate() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.keepOnScreen()

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

		// Show exhaust.
		bounds := p.sprite.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		// Where to spawn the exhaust.
		spawnPos := Vector{
			p.position.X + halfW + math.Sin(p.rotation)* exhaustSpawnOffset,
			p.position.Y + halfH + math.Cos(p.rotation)* -exhaustSpawnOffset,
		}

		p.game.exhaust = NewExhaust(spawnPos, p.rotation+180.0*math.Pi/180.0)

		// Move the player on the screen.
		p.position.X += dx
		p.position.Y += dy

		if !p.game.thrustPlayer.IsPlaying() {
			_ = p.game.thrustPlayer.Rewind()
			p.game.thrustPlayer.Play()
		}
	}
}

func (p *Player) isDoneAccelerating() {
	if inpututil.IsKeyJustReleased(ebiten.KeyUp) {
		if p.game.thrustPlayer.IsPlaying() {
			p.game.thrustPlayer.Pause()
		}
	}
}

func (p *Player) updateExhaustSprite() {
	if !ebiten.IsKeyPressed(ebiten.KeyUp) && p.game.exhaust != nil {
		p.game.exhaust = nil
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

package goasteroids

import (
	"asteroids/assets"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/solarlune/resolv"
)

const (
	baseMeteorVelocity   = 0.25                    // The base speed for meteors.
	meteorSpawnTime      = 100 * time.Millisecond  // How long before meteors spawn.
	meteorSpeedUpAmount  = 0.1                     // How much do we speed a meteor up when it's timer runs out.
	meteorSpeedUpTime    = 1000 * time.Millisecond // How long to wait to speed up meteors.
	cleanUpExplosionTime = 200 * time.Millisecond  // How long to wait before cleaning up an explosion.
)

// GameScene is the overall type for a game scene (e.g. TitleScene, GameScene, etc.).
type GameScene struct {
	player               *Player         // The player.
	baseVelocity         float64         // The base velocity for items in the game.
	meteorCount          int             // The counter for meteors.
	meteorSpawnTimer     *Timer          // The timer for spawning meteors.
	meteors              map[int]*Meteor // A map of meteors.
	meteorsForLevel      int             // # of meteors for a level.
	velocityTimer        *Timer          // The timer used for speeding up meteors.
	space                *resolv.Space   // The space for all collision objects.
	lasers               map[int]*Laser
	laserCount           int
	score                int
	explosionSmallSprite *ebiten.Image
	explosionSprite      *ebiten.Image
	explosionFrames      []*ebiten.Image
	cleanUpTimer         *Timer
	playerIsDead         bool
	audioContext         *audio.Context
	thrustPlayer         *audio.Player
}

// NewGameScene is a factory method for producing a new game. It's called once,
// when game play starts (and again when game play restarts).
func NewGameScene() *GameScene {
	g := &GameScene{
		meteorSpawnTimer:     NewTimer(meteorSpawnTime),
		baseVelocity:         baseMeteorVelocity,
		velocityTimer:        NewTimer(meteorSpeedUpTime),
		meteors:              make(map[int]*Meteor),
		meteorCount:          0,
		meteorsForLevel:      2,
		space:                resolv.NewSpace(ScreenWidth, ScreenHeight, 16, 16),
		lasers:               make(map[int]*Laser),
		laserCount:           0,
		explosionSprite:      assets.ExplosionSprite,
		explosionSmallSprite: assets.ExplosionSmallSprite,
		cleanUpTimer:         NewTimer(cleanUpExplosionTime),
	}
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)

	g.explosionFrames = assets.Explosion

	// Load audio.
	g.audioContext = audio.NewContext(48000)
	thrustPlayer, _ := g.audioContext.NewPlayer(assets.ThrustSound)
	g.thrustPlayer = thrustPlayer

	return g
}

// Update updates all game scene elements for the next draw. It's called once per tick.
func (g *GameScene) Update(state *State) error {
	g.player.Update()

	g.isPlayerDying()

	g.isPlayerDead(state)

	g.spawnMeteors()

	for _, m := range g.meteors {
		m.Update()
	}

	for _, l := range g.lasers {
		l.Update()
	}

	g.speedUpMeteors()

	g.isPlayerCollidingWithMeteor()

	g.isMeteorHitByPlayerLaser()

	g.cleanUpMeteorsAndAliens()

	return nil
}

// Draw draws all game scene elements to the screen. It's called once per frame.
func (g *GameScene) Draw(screen *ebiten.Image) {
	g.player.Draw(screen)

	// Draw meteors.
	for _, m := range g.meteors {
		m.Draw(screen)
	}

	// Draw lasers
	for _, l := range g.lasers {
		l.Draw(screen)
	}
}

// Layout is necessary to satisfy interface requirements from ebiten.
func (g *GameScene) Layout(outsideWidth, outsideHeight int) (ScreenWidth, ScreenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameScene) isMeteorHitByPlayerLaser() {
	for _, m := range g.meteors {
		for _, l := range g.lasers {
			if m.meteorObj.IsIntersecting(l.laserObj) {
				if m.meteorObj.Tags().Has(TagSmall) {
					// Small meteor hit.
					m.sprite = g.explosionSmallSprite
					g.score++
				} else {
					// Large meteor hit.
					oldPos := m.position

					m.sprite = g.explosionSprite

					g.score++

					numToSpawn := rand.IntN(numberOfSmallMeteorsFromLargeMeteor)
					for i := 0; i < numToSpawn; i++ {
						meteor := NewSmallMeteor(baseMeteorVelocity, g, len(m.game.meteors)-1)
						meteor.position = Vector{oldPos.X + float64(rand.IntN(100-50)+50), oldPos.Y + float64(rand.IntN(100-50)+50)}
						meteor.meteorObj.SetPosition(meteor.position.X, meteor.position.Y)
						g.space.Add(meteor.meteorObj)
						g.meteorCount++
						g.meteors[m.game.meteorCount] = meteor
					}
				}
			}
		}
	}
}

func (g *GameScene) isPlayerDying() {
	if g.player.isDying {
		g.player.dyingTimer.Update()
		if g.player.dyingTimer.IsReady() {
			g.player.dyingTimer.Reset()
			g.player.dyingCounter++
			if g.player.dyingCounter == 12 {
				g.player.isDying = false
				g.player.isDead = true
			} else if g.player.dyingCounter < 12 {
				g.player.sprite = g.explosionFrames[g.player.dyingCounter]
			} else {
				// Do nothing.
			}
		}
	}
}

func (g *GameScene) isPlayerDead(state *State) {
	if g.player.isDead {
		g.player.livesRemaining--
		if g.player.livesRemaining == 0 {
			state.SceneManager.GoToScene(NewGameScene())
		}
	}
}

// spawnMeteors creates meteors, up to the maximum for a level.
func (g *GameScene) spawnMeteors() {
	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()
		if len(g.meteors) < g.meteorsForLevel && g.meteorCount < g.meteorsForLevel {
			m := NewMeteor(g.baseVelocity, g, len(g.meteors)-1)
			g.space.Add(m.meteorObj)
			g.meteorCount++
			g.meteors[g.meteorCount] = m
		}
	}
}

// speedUpMeteors makes meteors move faster over time.
func (g *GameScene) speedUpMeteors() {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}
}

func (g *GameScene) isPlayerCollidingWithMeteor() {
	for _, m := range g.meteors {
		if m.meteorObj.IsIntersecting(g.player.playerObj) {
			if !g.player.isShielded {
				m.game.player.isDying = true
				break
			} else {
				// Bounce the meteor off the player.
			}
		}
	}
}

func (g *GameScene) cleanUpMeteorsAndAliens() {
	g.cleanUpTimer.Update()
	if g.cleanUpTimer.IsReady() {
		for i, m := range g.meteors {
			if m.sprite == g.explosionSprite || m.sprite == g.explosionSmallSprite {
				delete(g.meteors, i)
				g.space.Remove(m.meteorObj)
			}
		}
		g.cleanUpTimer.Reset()
	}
}

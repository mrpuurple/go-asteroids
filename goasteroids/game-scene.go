package goasteroids

import (
	"asteroids/assets"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/solarlune/resolv"
)

const (
	baseMeteorVelocity   = 0.25                    // The base speed for meteors.
	meteorSpawnTime      = 100 * time.Millisecond  // How long before meteors spawn.
	meteorSpeedUpAmount  = 0.1                     // How much do we speed a meteor up when it's timer runs out.
	meteorSpeedUpTime    = 1000 * time.Millisecond // How long to wait to speed up meteors.
	cleanUpExplosionTime = 200 * time.Millisecond  // The time to wait for cleaning up explosions.
	baseBeatWaitTime     = 1600                    // Base number of milliseconds to wait between beats of background. This is an int because we do math on it.
	numberOfStars        = 1000                    // The number of stars to display on the background.
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
	lasers               map[int]*Laser  // A map of lasers.
	laserCount           int             // A count of lasers currently in play; used as index for map lasers.
	score                int             // Current score.
	explosionSmallSprite *ebiten.Image   // A small explosion object.
	explosionSprite      *ebiten.Image   // A large explosion object.
	explosionFrames      []*ebiten.Image // The frames for explosion animation.
	cleanUpTimer         *Timer          // Timer to clean up objects.
	playerIsDead         bool            // Is the player dead.
	audioContext         *audio.Context  // The context used for our audio players.
	thrustPlayer         *audio.Player   // The audio player for thrust sound.
	exhaust              *Exhaust        // The object for exhaust (while accelerating).
	laserOnePlayer       *audio.Player   // The audio player for laser 1.
	laserTwoPlayer       *audio.Player   // The audio player for laser 2.
	laserThreePlayer     *audio.Player   // The audio player for laser 3.
	explosionPlayer      *audio.Player   // The explosion sound player.
	beatOnePlayer        *audio.Player   // The audio player for beat one (background sounds).
	beatTwoPlayer        *audio.Player   // The audio player for beat two sound (background sounds).
	beatTimer            *Timer          // The time for playing beats one and two.
	beatWaitTime         int             // The time to wait between beats. Reduced over time in each level.
	playBeatOne          bool            // Should we play beat one? Yes, if true, otherwise play beat two.
	stars                []*Star         // The stars for background.
	currentLevel         int             // The current level the player is on.
	shield               *Shield         // The player's shield.
	shieldsUpPlayer      *audio.Player   // The audio player for shields up sound.
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
		beatTimer:            NewTimer(2 * time.Second),
		beatWaitTime:         baseBeatWaitTime,
		currentLevel:         1,
	}
	g.player = NewPlayer(g)
	g.space.Add(g.player.playerObj)
	g.stars = GenerateStars(numberOfStars)

	g.explosionFrames = assets.Explosion

	// Load audio.
	g.audioContext = audio.NewContext(48000)

	thrustPlayer, _ := g.audioContext.NewPlayer(assets.ThrustSound)
	g.thrustPlayer = thrustPlayer

	laserOnePlayer, _ := g.audioContext.NewPlayer(assets.LaserOneSound)
	g.laserOnePlayer = laserOnePlayer

	laserTwoPlayer, _ := g.audioContext.NewPlayer(assets.LaserTwoSound)
	g.laserTwoPlayer = laserTwoPlayer

	laserThreePlayer, _ := g.audioContext.NewPlayer(assets.LaserThreeSound)
	g.laserThreePlayer = laserThreePlayer

	explosionPlayer, _ := g.audioContext.NewPlayer(assets.ExplosionSound)
	g.explosionPlayer = explosionPlayer

	beatOnePlayer, _ := g.audioContext.NewPlayer(assets.BeatOneSound)
	g.beatOnePlayer = beatOnePlayer

	beatTwoPlayer, _ := g.audioContext.NewPlayer(assets.BeatTwoSound)
	g.beatTwoPlayer = beatTwoPlayer

	shieldsUpPlayer, _ := g.audioContext.NewPlayer(assets.ShieldSound)
	g.shieldsUpPlayer = shieldsUpPlayer

	return g
}

// Update updates all game scene elements for the next draw. It's called once per tick.
func (g *GameScene) Update(state *State) error {
	g.player.Update()

	g.updateExhaust()

	g.updateShield()

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

	g.beatSound()

	g.isLevelComplete(state)

	return nil
}

// Draw draws all game scene elements to the screen. It's called once per frame.
func (g *GameScene) Draw(screen *ebiten.Image) {
	// Draw stars.
	for _, s := range g.stars {
		s.Draw(screen)
	}

	// Draw player.
	g.player.Draw(screen)

	// Draw exhaust.
	if g.exhaust != nil {
		g.exhaust.Draw(screen)
	}

	// Draw shield.
	if g.shield != nil {
		g.shield.Draw(screen)
	}

	// Draw meteors.
	for _, m := range g.meteors {
		m.Draw(screen)
	}

	// Draw lasers
	for _, l := range g.lasers {
		l.Draw(screen)
	}

	// Draw life indicators.
	if len(g.player.lifeIndicators) > 0 {
		for _, x := range g.player.lifeIndicators {
			x.Draw(screen)
		}
	}

	// Draw shield indicators.
	if len(g.player.shieldIndicators) > 0 {
		for _, x := range g.player.shieldIndicators {
			x.Draw(screen)
		}
	}

	// Draw hyperspace indicator.
	if g.player.hyperSpaceTimer == nil || g.player.hyperSpaceTimer.IsReady() {
		g.player.hyperspaceIndicator.Draw(screen)
	}

	// Update and draw score.
	textToDraw := fmt.Sprintf("%06d", g.score)
	op := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, 40)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.ScoreFont,
		Size:   24,
	}, op)

	// Update and draw high score.
	if g.score >= highScore {
		highScore = g.score
	}

	textToDraw = fmt.Sprintf("HIGH SCORE %06d", highScore)
	op = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, 75)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.ScoreFont,
		Size:   16,
	}, op)

	// Update and draw current level.
	textToDraw = fmt.Sprintf("LEVEL %d", g.currentLevel)
	op = &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			PrimaryAlign: text.AlignCenter,
		},
	}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(ScreenWidth/2, ScreenHeight-40)
	text.Draw(screen, textToDraw, &text.GoTextFace{
		Source: assets.LevelFont,
		Size:   16,
	}, op)
}

// Layout is necessary to satisfy interface requirements from ebiten.
func (g *GameScene) Layout(outsideWidth, outsideHeight int) (ScreenWidth, ScreenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *GameScene) updateShield() {
	if g.shield != nil {
		g.shield.Update()
	}
}

// isLevelComplete checks to see if the level is complete (all meteors destroyed).
func (g *GameScene) isLevelComplete(state *State) {
	if g.meteorCount >= g.meteorsForLevel && len(g.meteors) == 0 {
		// Level finished, so reset meteor velocity.
		g.baseVelocity = baseMeteorVelocity
		// Increase current level by one.
		g.currentLevel++

		// If we've done 5 levels, add a life.
		if g.currentLevel%5 == 0 {
			if g.player.livesRemaining < 6 {
				g.player.livesRemaining++
				x := float64(20 + len(g.player.lifeIndicators)*50.0)
				y := 20.0
				g.player.lifeIndicators = append(g.player.lifeIndicators, NewLifeIndicator(Vector{X: x, Y: y}))
			}
		}

		// Set the beat time to slowest.
		g.beatWaitTime = baseBeatWaitTime

		// Switch scenes.
		state.SceneManager.GoToScene(&LevelStartsScene{
			game:           g,
			nextLevelTimer: NewTimer(time.Second * 2),
			stars:          GenerateStars(numberOfStars),
		})
	}
}

func (g *GameScene) beatSound() {
	g.beatTimer.Update()
	if g.beatTimer.IsReady() {
		if g.playBeatOne {
			_ = g.beatOnePlayer.Rewind()
			g.beatOnePlayer.Play()
			g.beatTimer.Reset()
		} else {
			_ = g.beatTwoPlayer.Rewind()
			g.beatTwoPlayer.Play()
			g.beatTimer.Reset()
		}
		g.playBeatOne = !g.playBeatOne

		// Speed up the timer.
		if g.beatWaitTime > 400 {
			g.beatWaitTime = g.beatWaitTime - 25
			g.beatTimer = NewTimer(time.Millisecond * time.Duration(g.beatWaitTime))
		}
	}
}

func (g *GameScene) updateExhaust() {
	if g.exhaust != nil {
		g.exhaust.Update()
	}
}

func (g *GameScene) isMeteorHitByPlayerLaser() {
	for _, m := range g.meteors {
		for _, l := range g.lasers {
			if m.meteorObj.IsIntersecting(l.laserObj) {
				if m.meteorObj.Tags().Has(TagSmall) {
					// Small meteor hit.
					m.sprite = g.explosionSmallSprite
					g.score++

					if !g.explosionPlayer.IsPlaying() {
						_ = g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}
				} else {
					// Large meteor hit.
					oldPos := m.position

					m.sprite = g.explosionSprite

					g.score++

					if !g.explosionPlayer.IsPlaying() {
						_ = g.explosionPlayer.Rewind()
						g.explosionPlayer.Play()
					}

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
	if g.playerIsDead {
		g.player.livesRemaining--
		if g.player.livesRemaining == 0 {

			// New High Score?
			if g.score > originalHighScore {
				err := updateHighScore(g.score)
				if err != nil {
					log.Println(err)
				}
			}

			state.SceneManager.GoToScene(&GameOverScene{
				game:        g,
				meteors:     make(map[int]*Meteor),
				meteorCount: 5,
				stars:       GenerateStars(numberOfStars),
			})
		} else {
			score := g.score
			livesRemaining := g.player.livesRemaining
			lifeSlice := g.player.lifeIndicators[:len(g.player.lifeIndicators)-1]
			stars := g.stars
			shieldsRemaining := g.player.shieldsRemaining
			shieldIndicatorSlice := g.player.shieldIndicators

			g.Reset()

			g.player.livesRemaining = livesRemaining
			g.score = score
			g.player.lifeIndicators = lifeSlice
			g.stars = stars
			g.player.shieldsRemaining = shieldsRemaining
			g.player.shieldIndicators = shieldIndicatorSlice
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

				if !g.explosionPlayer.IsPlaying() {
					_ = g.explosionPlayer.Rewind()
					g.explosionPlayer.Play()
				}
				break
			} else {
				// Bounce the meteor.
				g.bounceMeteor(m)
			}
		}
	}
}

func (g *GameScene) bounceMeteor(m *Meteor) {
	direction := Vector{
		X: (ScreenWidth/2 - m.position.X) * -1,
		Y: (ScreenHeight/2 - m.position.Y) * -1,
	}
	normalizedDirection := direction.Normalize()
	velocity := g.baseVelocity // * 1.5 // Speed up the meteor after bounce ;).

	movement := Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	m.movement = movement
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

func (g *GameScene) Reset() {
	g.player = NewPlayer(g)
	g.meteors = make(map[int]*Meteor)
	g.meteorCount = 0
	g.lasers = make(map[int]*Laser)
	g.laserCount = 0
	g.score = 0
	g.meteorSpawnTimer.Reset()
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
	g.playerIsDead = false
	g.exhaust = nil
	g.space.RemoveAll()
	g.space.Add(g.player.playerObj)
	g.stars = GenerateStars(numberOfStars)
	g.player.shieldsRemaining = numberOfShields
	g.player.isShielded = false
}

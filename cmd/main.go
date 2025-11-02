package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// === Configuration constants ===
const (
	tileSize = 16
	gridW    = 40 // grid width in tiles
	gridH    = 30 // grid height in tiles
	screenW  = gridW * tileSize
	screenH  = gridH * tileSize

	baseTickSpeed = 8 // frames per movement (lower = faster)
	speedUpEvery  = 5 // increase speed every N food eaten
	speedDelta    = 1 // how much to reduce tickSpeed
)

// Point represents a position on the grid
type Point struct {
	X, Y int
}

// Game holds the full game state
type Game struct {
	snake     []Point
	dir       Point
	nextDir   Point
	food      Point
	score     int
	tickCount int
	tickSpeed int
	gameOver  bool
	started   bool

	tileSnake *ebiten.Image
	tileHead  *ebiten.Image
	tileFood  *ebiten.Image
}

// NewGame initializes the game instance and assets
func NewGame() *Game {
	g := &Game{
		tickSpeed: baseTickSpeed,
	}

	// Create simple colored tiles (snake, food, head)
	g.tileSnake = ebiten.NewImage(tileSize, tileSize)
	g.tileSnake.Fill(color.RGBA{R: 50, G: 200, B: 50, A: 255})
	g.tileHead = ebiten.NewImage(tileSize, tileSize)
	g.tileHead.Fill(color.RGBA{R: 0, G: 120, B: 255, A: 255})
	g.tileFood = ebiten.NewImage(tileSize, tileSize)
	g.tileFood.Fill(color.RGBA{R: 220, G: 40, B: 40, A: 255})

	rand.Seed(time.Now().UnixNano())
	g.reset()
	return g
}

// reset restarts the game to the initial state
func (g *Game) reset() {
	cx, cy := gridW/2, gridH/2
	g.snake = []Point{{cx, cy}, {cx - 1, cy}, {cx - 2, cy}}
	g.dir = Point{1, 0}
	g.nextDir = g.dir
	g.placeFood()
	g.score = 0
	g.tickCount = 0
	g.tickSpeed = baseTickSpeed
	g.gameOver = false
	g.started = true
}

// placeFood randomly positions the food on a free tile
func (g *Game) placeFood() {
	for {
		x := rand.Intn(gridW)
		y := rand.Intn(gridH)
		p := Point{x, y}
		coll := false
		for _, s := range g.snake {
			if s == p {
				coll = true
				break
			}
		}
		if !coll {
			g.food = p
			return
		}
	}
}

// Update handles game logic and user input
func (g *Game) Update() error {
	// Restart if R is pressed after game over
	if ebiten.IsKeyPressed(ebiten.KeyR) && g.gameOver {
		g.reset()
		return nil
	}

	// Input controls (arrow keys or WASD)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		if g.dir.Y != 1 {
			g.nextDir = Point{0, -1}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		if g.dir.Y != -1 {
			g.nextDir = Point{0, 1}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if g.dir.X != 1 {
			g.nextDir = Point{-1, 0}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if g.dir.X != -1 {
			g.nextDir = Point{1, 0}
		}
	}

	if g.gameOver {
		return nil
	}

	// Control the update rate (snake speed)
	g.tickCount++
	if g.tickCount < g.tickSpeed {
		return nil
	}
	g.tickCount = 0

	// Update direction
	g.dir = g.nextDir

	// Calculate new head position
	head := g.snake[0]
	newHead := Point{head.X + g.dir.X, head.Y + g.dir.Y}

	// Check wall collision
	if newHead.X < 0 || newHead.X >= gridW || newHead.Y < 0 || newHead.Y >= gridH {
		g.gameOver = true
		return nil
	}
	// Check self collision
	for _, s := range g.snake {
		if s == newHead {
			g.gameOver = true
			return nil
		}
	}

	// Add new head at the front
	g.snake = append([]Point{newHead}, g.snake...)

	// Check if food eaten
	if newHead == g.food {
		g.score++
		// Speed up every few points
		if g.score%speedUpEvery == 0 && g.tickSpeed > 2 {
			g.tickSpeed -= speedDelta
		}
		g.placeFood()
	} else {
		// Remove tail (no growth)
		g.snake = g.snake[:len(g.snake)-1]
	}

	return nil
}

// Draw renders all game elements on screen
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 10, G: 10, B: 10, A: 255})

	// Draw food
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.food.X*tileSize), float64(g.food.Y*tileSize))
	screen.DrawImage(g.tileFood, op)

	// Draw snake
	for i, p := range g.snake {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(p.X*tileSize), float64(p.Y*tileSize))
		if i == 0 {
			screen.DrawImage(g.tileHead, op)
		} else {
			screen.DrawImage(g.tileSnake, op)
		}
	}

	// HUD info
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), 4, 4)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Speed (frames/move): %d", g.tickSpeed), 4, 20)
	ebitenutil.DebugPrintAt(screen, "Controls: Arrow keys or WASD. R = restart (on game over)", 4, 36)

	// Game over message
	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER! Press R to restart.", screenW/2-120, screenH/2)
	}
}

// Layout defines the window’s internal resolution
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Snake - Go + Ebiten")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

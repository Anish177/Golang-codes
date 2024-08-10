package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	CellSize     = 10
	GridWidth    = ScreenWidth / CellSize
	GridHeight   = ScreenHeight / CellSize
)

var (
	grid          = make([][]bool, GridWidth)
	drawing       bool
	simulating    bool
	updateCounter int
)

type Game struct{}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		drawing = true
		x, y := ebiten.CursorPosition()
		drawCell(x, y)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		simulating = !simulating
	}

	if simulating && updateCounter%10 == 0 {
		updateGrid()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		drawing = false
	}

	if drawing {
		x, y := ebiten.CursorPosition()
		drawCell(x, y)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White) // Set the background color to white

	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			if grid[x][y] {
				ebitenutil.DrawRect(screen, float64(x*CellSize), float64(y*CellSize), float64(CellSize), float64(CellSize), color.Black)
			}
		}
		if !simulating {
			ebitenutil.DebugPrint(screen, "Press SPACE to start simulation")
		}
	}
	updateCounter++
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func updateGrid() {
	newGrid := make([][]bool, GridWidth)
	for x := 0; x < GridWidth; x++ {
		newGrid[x] = make([]bool, GridHeight)
	}

	for x := 0; x < GridWidth; x++ {
		for y := 0; y < GridHeight; y++ {
			liveNeighbors := 0
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if i == 0 && j == 0 {
						continue
					}
					nx, ny := x+i, y+j
					if nx >= 0 && nx < GridWidth && ny >= 0 && ny < GridHeight && grid[nx][ny] {
						liveNeighbors++
					}
				}
			}

			if grid[x][y] {
				if liveNeighbors < 2 || liveNeighbors > 3 {
					newGrid[x][y] = false
				} else {
					newGrid[x][y] = true
				}
			} else {
				if liveNeighbors == 3 {
					newGrid[x][y] = true
				}
			}
		}
	}
	grid = newGrid
}

func drawCell(x, y int) {
	gridX, gridY := x/CellSize, y/CellSize
	grid[gridX][gridY] = true
}

func main() {
	for x := 0; x < GridWidth; x++ {
		grid[x] = make([]bool, GridHeight)
	}
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Conway's Game of Life in Go")

	game := &Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2015 Martin Lindhe
// SPDX-FileCopyrightText: 2016 The Ebitengine Authors

// The original project is gol (https://github.com/martinlindhe/gol) by Martin Lindhe.

package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

var dA float64 = 1.0
var dB float64 = 0.5
var feed float64 = 0.062
var k float64 = 0.061

// World represents the game state.
type World struct {
	area   []Pix
	width  int
	height int
}

type Pix struct {
	a float64
	b float64
}

// NewWorld creates a new world.
func NewWorld(width, height int, maxInitLiveCells int) *World {
	w := &World{
		area:   make([]Pix, width*height),
		width:  width,
		height: height,
	}
	w.init(maxInitLiveCells)
	return w
}

// init inits world with a random state.
func (w *World) init(maxLiveCells int) {
	/*
		for i := 0; i < maxLiveCells; i++ {
			x := rand.IntN(w.width)
			y := rand.IntN(w.height)
			w.area[y*w.width+x]. = true
		}
	*/
	for x := 0; x < w.width; x++ {
		for y := 0; y < w.height; y++ {
			if x > 230 && x < 260 && y > 230 && y < 260 {
				w.area[y*w.width+x].a = 0
				w.area[y*w.width+x].b = 1
			} else {
				w.area[y*w.width+x].a = 1
				w.area[y*w.width+x].b = 0
			}
		}

	}
}

func laplaceA(w *World, x, y int) float64 {
	sumA := 0.0
	sumA += w.area[y*w.width+(x)].a * -1
	sumA += w.area[y*w.width+(x-1)].a * 0.2
	sumA += w.area[y*w.width+(x+1)].a * 0.2
	sumA += w.area[(y+1)*w.width+x].a * 0.2
	sumA += w.area[(y-1)*w.width+x].a * 0.2
	sumA += w.area[(y-1)*w.width+(x-1)].a * 0.05
	sumA += w.area[(y-1)*w.width+(x+1)].a * 0.05
	sumA += w.area[(y+1)*w.width+(x-1)].a * 0.05
	sumA += w.area[(y+1)*w.width+(x+1)].a * 0.05
	return sumA
}

func laplaceB(w *World, x, y int) float64 {
	sumB := 0.0
	sumB += w.area[y*w.width+(x)].b * -1
	sumB += w.area[y*w.width+(x-1)].b * 0.2
	sumB += w.area[y*w.width+(x+1)].b * 0.2
	sumB += w.area[(y+1)*w.width+x].b * 0.2
	sumB += w.area[(y-1)*w.width+x].b * 0.2
	sumB += w.area[(y-1)*w.width+(x-1)].b * 0.05
	sumB += w.area[(y-1)*w.width+(x+1)].b * 0.05
	sumB += w.area[(y+1)*w.width+(x-1)].b * 0.05
	sumB += w.area[(y+1)*w.width+(x+1)].b * 0.05
	return sumB
}

func constrain(a, b, c float64) float64 {
	if a < b {
		return b
	} else if a > c {
		return c
	} else {
		return a
	}
}

// Update game state by one tick.
func (w *World) Update(next *World) {
	var mouseX int
	var mouseY int
	mouseX, mouseY = ebiten.CursorPosition()

	leftClicked := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if leftClicked {
		w.area[mouseY*next.width+mouseX].a = 0.0
		w.area[mouseY*next.width+mouseX].b = 1.0
	}

	for x := 1; x < w.width-1; x++ {
		for y := 1; y < w.height-1; y++ {
			a := w.area[y*w.width+x].a
			b := w.area[y*w.width+x].b
			next.area[y*w.width+x].a = a + (dA * laplaceA(w, x, y)) - (a * b * b) + (feed * (1 - a))
			next.area[y*w.width+x].b = b + (dB * laplaceB(w, x, y)) + (a * b * b) - ((k + feed) * b)
			next.area[y*w.width+x].a = constrain(next.area[y*w.width+x].a, 0, 1)
			next.area[y*w.width+x].b = constrain(next.area[y*w.width+x].b, 0, 1)

		}

	}
	swap(w, next)
	/*
		width := w.width
		height := w.height
		next := make([]bool, width*height)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				pop := neighbourCount(w.area, width, height, x, y)
				switch {
				case pop < 2:
					// rule 1. Any live cell with fewer than two live neighbours
					// dies, as if caused by under-population.
					next[y*width+x] = false

				case (pop == 2 || pop == 3) && w.area[y*width+x]:
					// rule 2. Any live cell with two or three live neighbours
					// lives on to the next generation.
					next[y*width+x] = true

				case pop > 3:
					// rule 3. Any live cell with more than three live neighbours
					// dies, as if by over-population.
					next[y*width+x] = false

				case pop == 3:
					// rule 4. Any dead cell with exactly three live neighbours
					// becomes a live cell, as if by reproduction.
					next[y*width+x] = true
				}
			}
		}
		w.area = next
	*/
}

// Draw paints current game state.
func (w *World) Draw(pix []byte, next *World) {
	for i := range w.area {
		a := next.area[i].a
		b := next.area[i].b
		c := math.Floor((a - b) * 255)
		c = constrain(c, 0.0, 255.0)
		pix[4*i] = byte(c)
		pix[4*i+1] = byte(c)
		pix[4*i+2] = byte(c)
		pix[4*i+3] = 0xff
		/*
			if v.a > 0 {
				pix[4*i] = byte(c)
				pix[4*i+1] = byte(c)
				pix[4*i+2] = byte(c)
				pix[4*i+3] = 0

			} else {
				pix[4*i] = 0xff
				pix[4*i+1] = 0xff
				pix[4*i+2] = 0xff
				pix[4*i+3] = 0xff
			}
		*/
	}
}

func swap(w *World, next *World) {
	temp := *w
	*w = *next
	*next = temp
}

const (
	screenWidth  = 500
	screenHeight = 500
)

type Game struct {
	grid   *World
	next   *World
	pixels []byte
}

func (g *Game) Update() error {
	g.grid.Update(g.next)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.grid.Draw(g.pixels, g.next)
	screen.WritePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetTPS(150)
	g := &Game{
		grid: NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
		next: NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Reaction Diffusion")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

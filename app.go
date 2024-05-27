/*
@author: sk
@date: 2024/5/26
*/
package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type App struct {
	Chip8 *Chip8
}

func NewApp() *App {
	return &App{Chip8: NewChip8()}
}

var (
	keys = []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4,
		ebiten.KeyQ, ebiten.KeyW, ebiten.KeyE, ebiten.KeyR,
		ebiten.KeyA, ebiten.KeyS, ebiten.KeyD, ebiten.KeyF,
		ebiten.KeyZ, ebiten.KeyX, ebiten.KeyC, ebiten.KeyV}
)

func (a *App) Update() error {
	for i, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			a.Chip8.Keypad[i] = true
		}
		if inpututil.IsKeyJustReleased(key) {
			a.Chip8.Keypad[i] = false
		}
	}
	a.Chip8.Update()
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{A: 0xff})
	for i := 0; i < VideoW; i++ {
		for j := 0; j < VideoH; j++ {
			if a.Chip8.Video[i][j] {
				screen.Set(i, j, color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff})
			}
		}
	}
}

func (a *App) Layout(w, h int) (int, int) {
	return w / VideoS, h / VideoS
}

func (a *App) Load(file string) {
	a.Chip8.Load(file)
}

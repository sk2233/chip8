/*
@author: sk
@date: 2024/5/26
*/
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// https://github.com/dmatlack/chip8/tree/master/roms/games game roms
// https://austinmorlan.com/posts/chip8_emulator/

func main() {
	app := NewApp()
	app.Load("res/Tetris [Fran Dachille, 1991].ch8")
	ebiten.SetWindowSize(VideoW*VideoS, VideoH*VideoS)
	err := ebiten.RunGame(app)
	HandleErr(err)
}

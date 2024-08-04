/*
@author: sk
@date: 2024/5/26
*/
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Q 变形 W 左 E 右 R 加速下 仅针对 Tetris [Fran Dachille, 1991].ch8
// https://github.com/dmatlack/chip8/tree/master/roms/games game roms
// https://austinmorlan.com/posts/chip8_emulator/

func main() {
	app := NewApp()
	app.Load("res/Tetris [Fran Dachille, 1991].ch8")
	ebiten.SetWindowSize(VideoW*VideoS, VideoH*VideoS)
	err := ebiten.RunGame(app)
	HandleErr(err)
}

/*
@author: sk
@date: 2024/5/26
*/
package main

import (
	"math/rand"
	"os"
)

type Chip8 struct {
	Rs   [16]uint8   // 寄存器
	Mem  [4096]uint8 // 可用内存
	Addr uint16      // 临时地址，因为内存大小需要16位的
	Pc   uint16
	// 方法栈，记录 pc的位置
	Stack [16]uint16
	Sp    uint8
	// 触发计时器&声音计时器，每秒 60 次递减直到为0
	DelayTimer uint8
	SoundTimer uint8
	Keypad     [16]bool             // 按键映射
	Video      [VideoW][VideoH]bool // 屏幕像素
	OpCode     uint16
}

func (c *Chip8) Load(file string) {
	bs, err := os.ReadFile(file)
	HandleErr(err)
	copy(c.Mem[StartAddr:], bs)
}

// 00e0
func (c *Chip8) op00e0() {
	for i := 0; i < len(c.Video); i++ {
		for j := 0; j < len(c.Video[i]); j++ {
			c.Video[i][j] = false
		}
	}
}

// 00ee
func (c *Chip8) op00ee() {
	c.Sp--
	c.Pc = c.Stack[c.Sp]
}

// 1xxx 都匹配
func (c *Chip8) op1xxx() {
	addr := c.OpCode & 0x0fff
	c.Pc = addr
}

// 2xxx 都匹配
func (c *Chip8) op2xxx() {
	addr := c.OpCode & 0x0fff
	c.Stack[c.Sp] = c.Pc
	c.Sp++
	c.Pc = addr
}

// 3ivv 都匹配
func (c *Chip8) op3ivv() {
	ri := (c.OpCode & 0xf00) >> 8 // 寄存器索引
	val := uint8(c.OpCode & 0xff) // 对应的值
	if c.Rs[ri] == val {          // 跳过下条指令
		c.Pc += 2
	}
}

// 4ivv 都匹配
func (c *Chip8) op4ivv() {
	ri := (c.OpCode & 0xf00) >> 8 // 寄存器索引
	val := uint8(c.OpCode & 0xff) // 对应的值
	if c.Rs[ri] != val {          // 跳过下条指令
		c.Pc += 2
	}
}

// 5ij0 都匹配
func (c *Chip8) op5ij0() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	if c.Rs[ri] == c.Rs[rj] {
		c.Pc += 2
	}
}

// 6ivv 匹配
func (c *Chip8) op6ivv() {
	ri := (c.OpCode & 0xf00) >> 8
	val := uint8(c.OpCode & 0xff)
	c.Rs[ri] = val
}

// 7ivv
func (c *Chip8) op7ivv() {
	ri := (c.OpCode & 0xf00) >> 8
	val := uint8(c.OpCode & 0xff)
	c.Rs[ri] += val
}

// 8ij0
func (c *Chip8) op8ij0() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	c.Rs[ri] = c.Rs[rj]
}

// 8ij1
func (c *Chip8) op8ij1() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	c.Rs[ri] |= c.Rs[rj]
}

// 8ij2
func (c *Chip8) op8ij2() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	c.Rs[ri] &= c.Rs[rj]
}

// 8ij3
func (c *Chip8) op8ij3() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	c.Rs[ri] ^= c.Rs[rj]
}

// 8ij4
func (c *Chip8) op8ij4() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	sum := uint16(c.Rs[ri]) + uint16(c.Rs[rj])
	if sum > 0xff { // 更新条件寄存器
		c.Rs[0xf] = 1
	} else {
		c.Rs[0xf] = 0
	}
	c.Rs[ri] = uint8(sum & 0xff)
}

// 8ij5
func (c *Chip8) op8ij5() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	if c.Rs[ri] > c.Rs[rj] {
		c.Rs[0xf] = 1
	} else {
		c.Rs[0xf] = 0
	}
	c.Rs[ri] -= c.Rs[rj]
}

// 8ij6  j ???
func (c *Chip8) op8ij6() {
	ri := (c.OpCode & 0xf00) >> 8
	c.Rs[0xf] = c.Rs[ri] & 0x1
	c.Rs[ri] >>= 1
}

// 8ij7
func (c *Chip8) op8ij7() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	if c.Rs[rj] > c.Rs[ri] {
		c.Rs[0xf] = 1
	} else {
		c.Rs[0xf] = 0
	}
	c.Rs[ri] = c.Rs[rj] - c.Rs[ri]
}

// 8ijE j ???
func (c *Chip8) op8ijE() {
	ri := (c.OpCode & 0xf00) >> 8
	c.Rs[0xf] = (c.Rs[ri] & 0x80) >> 7
	c.Rs[ri] <<= 1
}

// 9ij0
func (c *Chip8) op9ij0() {
	ri := (c.OpCode & 0xf00) >> 8
	rj := (c.OpCode & 0xf0) >> 4
	if c.Rs[ri] != c.Rs[rj] {
		c.Pc += 2
	}
}

// Axxx
func (c *Chip8) opAxxx() {
	addr := c.OpCode & 0xfff
	c.Addr = addr
}

// Bxxx
func (c *Chip8) opBxxx() {
	offset := c.OpCode & 0xfff
	c.Pc = uint16(c.Rs[0]) + offset
}

// Cimm
func (c *Chip8) opCimm() {
	ri := (c.OpCode & 0xf00) >> 8
	mask := uint8(c.OpCode & 0xff)
	c.Rs[ri] = uint8(rand.Uint32()) & mask
}

// Dxyh
func (c *Chip8) opDxyh() {
	rx := (c.OpCode & 0xf00) >> 8
	ry := (c.OpCode & 0xf0) >> 4
	xPos := c.Rs[rx] % VideoW
	yPos := c.Rs[ry] % VideoH
	h := c.OpCode & 0xf
	c.Rs[0xf] = 0                    // 记录是否碰撞
	for i := uint16(0); i < h; i++ { // 逐行检查
		temp := c.Mem[c.Addr+i] // 每行固定宽 8
		for j := uint8(0); j < 8; j++ {
			if (temp & (0x80 >> j)) > 0 { // 需要绘制
				if c.Video[xPos+j][yPos+uint8(i)] { // 发生了碰撞 标记
					c.Rs[0xf] = 1
				} // 像素每次绘制都是取反的
				c.Video[xPos+j][yPos+uint8(i)] = !c.Video[xPos+j][yPos+uint8(i)]
			}
		}
	}
}

// Ei9E
func (c *Chip8) opEi9E() {
	ri := (c.OpCode & 0xf00) >> 8
	key := c.Rs[ri]
	if c.Keypad[key] { // 按键判断
		c.Pc += 2
	}
}

// EiA1
func (c *Chip8) opEiA1() {
	ri := (c.OpCode & 0xf00) >> 8
	key := c.Rs[ri]
	if !c.Keypad[key] { // 按键判断
		c.Pc += 2
	}
}

// Fi07
func (c *Chip8) opFi07() {
	ri := (c.OpCode & 0xf00) >> 8
	c.Rs[ri] = c.DelayTimer
}

// Fi0A
func (c *Chip8) opFi0A() {
	ri := (c.OpCode & 0xf00) >> 8
	for i := 0; i < len(c.Keypad); i++ {
		if c.Keypad[i] {
			c.Rs[ri] = uint8(i)
			return
		}
	} // 等待按钮按下，否则重复前一个指令，即等待按钮按下
	c.Pc -= 2
}

// Fi15
func (c *Chip8) opFi15() {
	ri := (c.OpCode & 0xf00) >> 8
	c.DelayTimer = c.Rs[ri]
}

// Fi18
func (c *Chip8) opFi18() {
	ri := (c.OpCode & 0xf00) >> 8
	c.SoundTimer = c.Rs[ri]
}

// Fi1E
func (c *Chip8) opFi1E() {
	ri := (c.OpCode & 0xf00) >> 8
	c.Addr += uint16(c.Rs[ri])
}

// Fi29
func (c *Chip8) opFi29() {
	ri := (c.OpCode & 0xf00) >> 8
	num := c.Rs[ri]
	c.Addr = FontStartAddr + uint16(num*5) // 设置第 num 个数字预定义精灵的位置 一个数字宽 8（1byte）高 5
}

// Fi33
func (c *Chip8) opFi33() {
	ri := (c.OpCode & 0xf00) >> 8
	val := c.Rs[ri]
	// 把对应的值转换为 10 进制存储到对应位置上 只到百位
	c.Mem[c.Addr+2] = val % 10
	val /= 10
	c.Mem[c.Addr+1] = val % 10
	val /= 10
	c.Mem[c.Addr] = val % 10
}

// Fi55
func (c *Chip8) opFi55() { // 寄存器数据保存
	ri := (c.OpCode & 0xf00) >> 8
	for i := uint16(0); i <= ri; i++ {
		c.Mem[c.Addr+i] = c.Rs[i]
	}
}

// Fi65
func (c *Chip8) opFi65() { // 寄存器数据加载
	ri := (c.OpCode & 0xf00) >> 8
	for i := uint16(0); i <= ri; i++ {
		c.Rs[i] = c.Mem[c.Addr+i]
	}
}

func (c *Chip8) Update() {
	c.OpCode = (uint16(c.Mem[c.Pc]) << 8) | uint16(c.Mem[c.Pc+1])
	c.Pc += 2
	nums := []uint8{uint8((c.OpCode & 0xf000) >> 12), uint8((c.OpCode & 0xf00) >> 8),
		uint8((c.OpCode & 0xf0) >> 4), uint8(c.OpCode & 0xf)}
	switch nums[0] {
	case 0:
		switch c.OpCode {
		case 0x00E0:
			c.op00e0()
		case 0x00EE:
			c.op00ee()
		}
	case 1:
		c.op1xxx()
	case 2:
		c.op2xxx()
	case 3:
		c.op3ivv()
	case 4:
		c.op4ivv()
	case 5:
		if nums[3] == 0 {
			c.op5ij0()
		}
	case 6:
		c.op6ivv()
	case 7:
		c.op7ivv()
	case 8:
		switch nums[3] {
		case 0:
			c.op8ij0()
		case 1:
			c.op8ij1()
		case 2:
			c.op8ij2()
		case 3:
			c.op8ij3()
		case 4:
			c.op8ij4()
		case 5:
			c.op8ij5()
		case 6:
			c.op8ij6()
		case 7:
			c.op8ij7()
		case 0xE:
			c.op8ijE()
		}
	case 9:
		if nums[3] == 0 {
			c.op9ij0()
		}
	case 0xA:
		c.opAxxx()
	case 0xB:
		c.opBxxx()
	case 0xC:
		c.opCimm()
	case 0xD:
		c.opDxyh()
	case 0xE:
		switch (nums[2] << 4) | nums[3] {
		case 0x9E:
			c.opEi9E()
		case 0xA1:
			c.opEiA1()
		}
	case 0xF:
		switch (nums[2] << 4) | nums[3] {
		case 0x07:
			c.opFi07()
		case 0x0A:
			c.opFi0A()
		case 0x15:
			c.opFi15()
		case 0x18:
			c.opFi18()
		case 0x1E:
			c.opFi1E()
		case 0x29:
			c.opFi29()
		case 0x33:
			c.opFi33()
		case 0x55:
			c.opFi55()
		case 0x65:
			c.opFi65()
		}
	}
	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
	if c.SoundTimer > 0 {
		c.SoundTimer--
	}
}

func NewChip8() *Chip8 {
	res := &Chip8{Pc: StartAddr}
	copy(res.Mem[FontStartAddr:], FontData)
	return res
}

package main

import (
	"fmt"
	"strings"
	"github.com/nsf/termbox-go"
)

type JTerm struct {
	x, y          int
	width, height int
	shift         int
	jtree         *JTree
}

func NewJTerm(jtree *JTree) *JTerm {
	termbox.Init()
	w, h := termbox.Size()
	return &JTerm{0, 0, w, h, 0, jtree}
}

//鼠标移动到指定位置
func (this *JTerm) MoveTo(x, y, shift int) {
	this.x = x
	this.y = y
	this.shift = shift

	termbox.SetCursor(this.x, this.y)
}

func (this *JTerm) MoveToRela(x, y int) {
	this.MoveTo(x, y, this.shift)
}

//指定距离位移
func (this *JTerm) Move(x, y int) {
	this.x += x

	tempY := this.y + y
	tempS := this.shift + y
	if tempY > this.height {
		this.shift = tempS
	} else if tempY < 0 && tempS >= 0 {
		this.shift = tempS
	} else if tempY >= 0 && tempY <= this.height {
		this.y = tempY
	}

	termbox.SetCursor(this.x, this.y)
}

func tprint(x, y int, str string, fg termbox.Attribute, bg termbox.Attribute) int {
	for i, v := range str {
		termbox.SetCell(x+i, y, v, fg, bg)
	}

	return len(str)
}

func writeNum(x, y, n, max int) int{
	s := fmt.Sprintf("%d", n)
	s += strings.Repeat(" ", max - len(s))
	return tprint(x, y, s, NUMBER, DEFAULT)
}

//渲染
func (this *JTerm) Render() {
	termbox.Clear(DEFAULT, DEFAULT)

	for i := 0; i < this.height; i++ {
		line := this.jtree.Line(i + this.shift)
		if line == nil {
			continue
		}
		maxlen := len(fmt.Sprintf("%d", len(this.jtree.lines)))
		s := writeNum(0, i, i + 1 + this.shift, maxlen)
		for j, v := range line {
			termbox.SetCell(j+s, i, v.Val, v.Attribute, DEFAULT)
		}

	}

	termbox.Flush()
}

//窗口调整
func (this *JTerm) Resize(x, y int) {
	this.width = x
	this.height = y
	if this.y > y {
		this.y = this.height
	}
}

func (this *JTerm) ListenAction() {
	termbox.SetInputMode(termbox.InputEsc)
	defer termbox.Close()
	this.Render()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 0 {
				switch ev.Key {
				case termbox.KeyArrowLeft:
					this.Move(-1, 0)
				case termbox.KeyArrowRight:
					this.Move(1, 0)
				case termbox.KeyArrowUp:
					this.Move(0, -1)
				case termbox.KeyArrowDown:
					this.Move(0, 1)
				case termbox.KeyEnter:
					this.jtree.Toggle(this.y + this.shift)
				case termbox.KeyF3:
					this.jtree.OpenToggleAll()
				case termbox.KeyF2:
					this.jtree.EmptyToggleAll()
					this.MoveTo(0, 0, 0)
				case termbox.KeyCtrlC:
					break loop
				}
			} else {
				switch ev.Ch {
				case 'w', 'W':
					this.Move(0, -5)
				case 's', 'S':
					this.Move(0, 5)
				}
			}

		case termbox.EventResize:
			this.Resize(ev.Width, ev.Height)
		case termbox.EventMouse:

		}

		this.Render()

	}
}

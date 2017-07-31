package main

import (
	"fmt"
	"strings"
	"github.com/nsf/termbox-go"
)

type JTerm struct {
	x, y           int
	width, height  int
	shiftY, shiftX int
	jtree          *JTree
}

func NewJTerm(jtree *JTree) *JTerm {
	termbox.Init()
	w, h := termbox.Size()
	return &JTerm{0, 0, w, h, 0, 0, jtree}
}

//鼠标移动到指定位置
func (this *JTerm) MoveTo(x, y, shiftY int, shiftX int) {
	this.x = x
	this.y = y
	this.shiftY = shiftY
	this.shiftX = shiftX
	termbox.SetCursor(this.x, this.y)
}

func (this *JTerm) MoveToRela(x, y int) {
	this.MoveTo(x, y, this.shiftY, this.shiftX)
}

//返回一个值的正常边界值
func correctUnit(x int, lborder, rborder int) int{
	if(x < lborder) {
		return lborder
	}

	if(x >= rborder) {
		return rborder - 1
	}

	return x
}


//移动某一个元素通过他的偏移和边界
func MoveUnit(x int, objX, shiftX int, lborder, rborder int) (int, int){
	tempX := objX + x
	tempSX := shiftX + x

	if tempX >= rborder {
		objX = correctUnit(tempX, lborder, rborder)
		shiftX = tempSX
	}else if tempX < 0{
		objX = correctUnit(tempX, lborder, rborder)
		if tempSX > 0{
			shiftX = tempSX
		}else {
			shiftX = 0
		}
	}else {
		objX = correctUnit(tempX, lborder, rborder)
	}

	return objX, shiftX
}

//指定距离位移
func (this *JTerm) Move(x, y int) {
	this.x, this.shiftX = MoveUnit(x, this.x, this.shiftX, 0, this.width)
	this.y, this.shiftY = MoveUnit(y, this.y, this.shiftY, 0, this.height)
	termbox.SetCursor(this.x, this.y)
}

func tprint(x, y int, str string, fg termbox.Attribute, bg termbox.Attribute) int {
	for i, v := range str {
		termbox.SetCell(x+i, y, v, fg, bg)
	}

	return len(str)
}

//画下标
func writeNum(x, y, n, max int, equal bool) int{
	s := fmt.Sprintf("%d", n)
	s += strings.Repeat(" ", max - len(s))

	bg := DEFAULT
	if equal {
		bg = BOOLEAN
	}

	return tprint(x, y, s, NUMBER, bg)
}

//渲染
func (this *JTerm) Render() {
	termbox.Clear(DELIM, DEFAULT)

	for i := 0; i < this.height; i++ {
		line := this.jtree.Line(i + this.shiftY)
		if line == nil {
			continue
		}
		//渲染下标
		maxlen := len(fmt.Sprintf("%d", len(this.jtree.lines)))
		//数字+空格位数
		s := writeNum(0 - this.shiftX, i, i + 1 + this.shiftY, maxlen, this.y == i)

		for j, v := range line {
			termbox.SetCell(j + s - this.shiftX, i, v.Val, v.Attribute, DEFAULT)
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
	if this.x > x {
		this.x = this.width
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
					this.jtree.Toggle(this.y + this.shiftY)
				case termbox.KeyCtrlC:
					break loop
				}
			} else {
				switch ev.Ch {
				case 'w', 'W':
					this.Move(0, -5)
				case 's', 'S':
					this.Move(0, 5)
				case 'a', 'A':
					this.Move(-5, 0)
				case 'd', 'D':
					this.Move(5, 0)
				case 'q', 'Q':
					this.jtree.OpenToggleAll()
				case 'e', 'E':
					this.jtree.EmptyToggleAll()
					this.MoveTo(0, 0, 0, 0)
				}
			}

		case termbox.EventResize:
			this.Resize(ev.Width, ev.Height)
		case termbox.EventMouse:

		}

		this.Render()

	}
}

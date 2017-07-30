package main

import(
	"github.com/nsf/termbox-go"
	//"fmt"
)

type JTerm struct {
	x, y int
	width, height int
	shift int
	jtree *JTree
}

func NewJTerm(jtree *JTree) *JTerm{
	termbox.Init()
	w, h := termbox.Size()
	return &JTerm{0, 0, w, h, 0, jtree}
} 

//鼠标移动到指定位置
func (this *JTerm) MoveTo(x, y int) {
	this.x = x
	this.y = y
	this.shift = 0

	termbox.SetCursor(this.x, this.y)
}

//指定距离位移
func (this *JTerm) Move(x, y int) {
	this.x += x

	tempY := this.y + y

	if tempY > this.height || this.shift < 0{
		this.shift++
	}else if tempY < 0 || this.shift > 0{
		this.shift--
	}else {
		this.y = tempY
	}
	
	
	termbox.SetCursor(this.x, this.y)
}

//渲染
func (this *JTerm) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	//fmt.Println(this.y, this.jtree.lineMap, this.jtree.extendedLine, this.jtree.startLine)
	for i := 0; i < this.height; i++ {
		line := this.jtree.Line(i)
		if line == nil {
			continue
		}
		for j, v := range line {
			termbox.SetCell(j, i - this.shift, v.Val, v.Attribute, termbox.ColorDefault)
		}

	} 

	termbox.Flush()
}

//窗口调整
func (this *JTerm) Resize(x, y int) {
	this.width = x
	this.height = y
}

func (this *JTerm) ListenAction() {
	termbox.SetInputMode(termbox.InputEsc)
	defer termbox.Close()
	this.Render()

	loop :
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyArrowLeft: this.Move(-1, 0)
				case termbox.KeyArrowRight : this.Move(1, 0)
				case termbox.KeyArrowUp : this.Move(0, -1)
				case termbox.KeyArrowDown : this.Move(0, 1)
				case termbox.KeyEnter : this.jtree.Toggle(this.y)
				case termbox.KeyF3 : this.jtree.OpenToggleAll()
				case termbox.KeyF2 : this.jtree.EmptyToggleAll(); this.MoveTo(0, 0)
				case termbox.KeyCtrlC : break loop
				}
			case termbox.EventResize : this.Resize(ev.Width, ev.Height)
			}

			this.Render()
		}
}

package main

import(
	"github.com/nsf/termbox-go"
)


type JWriter interface {
	WriteLine(string, termbox.Attribute)
	NewLine()
} 

type Char struct {
	Val rune
	Attribute termbox.Attribute
}

type Line []Char

type JWriterMan struct {
	Lines []Line
	line int
}

func NewJWriterMan() *JWriterMan{
	p := &JWriterMan {
		line : 0,
	}

	p.Lines = append(p.Lines, Line{})

	return p
}

func (this *JWriterMan) NewLine() {
	this.Lines = append(this.Lines, Line{})
	this.line++
}

func (this *JWriterMan) WriteLine(buf string, attr termbox.Attribute) {
	for _, v := range buf {
		this.Lines[this.line] = append(this.Lines[this.line], Char{v, attr})
	}
}

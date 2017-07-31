package main

import (
	"encoding/json"
	"fmt"
	"github.com/nsf/termbox-go"
	"strconv"
	"strings"
)

const (
	KEY     = termbox.ColorCyan
	STRING  = termbox.ColorGreen
	NUMBER  = termbox.ColorYellow
	DELIM   = termbox.ColorWhite
	BOOLEAN = termbox.ColorRed
	DEFAULT = termbox.ColorBlack
)

type JParse struct {
	content []byte
	depth   int
	JWriter
}

func NewJParse(content []byte, wman JWriter) *JParse {
	return &JParse{content, 0, wman}
}

func (this *JParse) Format() {
	var m map[string]interface{}

	err := json.Unmarshal(this.content, &m)

	if err != nil {
		panic("parse json failed!")
	}

	this.format(m)

}

//格式化基本数据
func (this *JParse) format(i interface{}) {
	switch v := i.(type) {
	case map[string]interface{}:
		this.formatObject(v)
	case []interface{}:
		this.formatArray(v)
	case float64:
		this.WriteLine(strconv.FormatFloat(v, 'f', -1, 64), NUMBER)
	case string:
		this.WriteLine(fmt.Sprintf("\"%s\"", v), STRING)
	case bool:
		this.WriteLine(strconv.FormatBool(v), BOOLEAN)
	case nil:
		this.WriteLine("null", DELIM)
	}
}

//格式化json对象数据
func (this *JParse) formatObject(m map[string]interface{}) {
	if len(m) == 0 {
		this.WriteLine("{}", DELIM)
		this.NewLine()
		return
	}

	this.WriteLine("{", DELIM)

	this.depth++
	this.NewLine()

	size := len(m)
	index := 0
	for k, v := range m {
		this.writeIndet()
		this.WriteLine(fmt.Sprintf("\"%s\"", k), KEY)
		this.WriteLine(":", DELIM)
		this.format(v)
		if index != size-1 {
			this.WriteLine(",", DELIM)
		}
		this.NewLine()
		index++
	}

	this.depth--
	this.writeIndet()
	this.WriteLine("}", DELIM)
}

//格式化数组数据
func (this *JParse) formatArray(a []interface{}) {
	if len(a) == 0 {
		this.WriteLine("[]", DELIM)
		this.NewLine()
		return
	}

	this.WriteLine("[", DELIM)

	this.depth++
	this.NewLine()

	size := len(a)
	index := 0
	for _, v := range a {
		this.writeIndet()
		this.format(v)
		if index != size-1 {
			this.WriteLine(",", DELIM)
		}
		this.NewLine()
		index++
	}

	this.depth--
	this.writeIndet()
	this.WriteLine("]", DELIM)
}

func (this *JParse) writeIndet() {
	this.WriteLine(strings.Repeat(" ", this.depth*4), DELIM)
}

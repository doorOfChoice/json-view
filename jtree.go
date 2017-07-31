package main

import (
	"unicode"
)

type JTree struct {
	lines        []Line      //完整的结果
	extendedLine map[int]int //展开的行
	startLine    map[int]int //最开始行的起始点和结束点
	lineMap      map[int]int //虚拟地址对应的实际地址
}

func NewJTree(lines []Line) *JTree {
	tree := &JTree{
		lines:        lines,
		extendedLine: make(map[int]int),
		lineMap:      make(map[int]int),
		startLine:    parseStartMap(lines),
	}

	tree.parseLineMap()

	return tree
}

//返回指定虚拟地址的实际数据
func (this *JTree) Line(virLn int) Line {
	actualLn, ok := this.lineMap[virLn]
	if ok {
		//如果是开始行并且已经展开就返回省略号
		if !this.isExtended(actualLn) && this.isStart(actualLn) {
			return this.linedot(actualLn)
		}

		return this.lines[actualLn]
	}

	return nil
}

//折叠: 闭则开，开则闭
func (this *JTree) Toggle(virLn int) {
	actualLn, ok := this.lineMap[virLn]
	if ok {
		//判断是否是开始行， 不是则无效
		if this.isStart(actualLn) {
			if !this.isExtended(actualLn) {
				this.extendedLine[actualLn] = actualLn
			} else {
				delete(this.extendedLine, actualLn)
			}
		}
		this.parseLineMap()
	}
}

//合上所有的展开数据
func (this *JTree) EmptyToggleAll() {
	this.extendedLine = make(map[int]int)
	this.parseLineMap()
}

//展开所有的闭合数据
func (this *JTree) OpenToggleAll() {
	for key, _ := range this.startLine {
		this.extendedLine[key] = 1
	}

	this.parseLineMap()
}

//生成省略号
func (this *JTree) linedot(actualLn int) Line {
	line := this.lines[actualLn]
	line = append(line, Char{'…', DELIM})

	for _, v1 := range this.lines[this.startLine[actualLn]] {
		if !unicode.IsSpace(v1.Val) {
			line = append(line, Char{v1.Val, DELIM})
		}
	}

	return line
}

//是否展开
func (this *JTree) isExtended(actualLn int) bool {
	_, ok := this.extendedLine[actualLn]
	return ok
}

//是否是开始行
func (this *JTree) isStart(actualLn int) bool {
	_, ok := this.startLine[actualLn]
	return ok
}

//转化虚拟地址对应的实际地址
func (this *JTree) parseLineMap() {
	n := make(map[int]int)
	//跳转到的位置
	skipTill := 0

	virLn := 0
	for actualLn, _ := range this.lines {
		if actualLn < skipTill {
			continue
		}
		//如果是开始行且没展开则跳转
		if !this.isExtended(actualLn) && this.isStart(actualLn) {
			skipTill = this.startLine[actualLn] + 1
		}

		n[virLn] = actualLn
		virLn++
	}

	this.lineMap = n
}

//寻找到所有的开始行
func parseStartMap(lines []Line) map[int]int {
	startMap := make(map[int]int)
	resultMap := make(map[int]int)
	index := 0

	for i, v1 := range lines {
		for _, v2 := range v1 {
			switch v2.Val {
			case '{', '[':
				startMap[index] = i
				index++
			case '}', ']':
				if s := startMap[index-1]; s != i {
					resultMap[s] = i
				}
				index--
			}
		}
	}

	return resultMap
}

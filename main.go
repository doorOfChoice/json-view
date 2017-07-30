package main

import(
	// "encoding/json"
	 "io/ioutil"
	 "os"
	"fmt"
//"github.com/nsf/termbox-go"
)


func main() {
	f, err := os.Open("../a.json")

	if(err != nil) {
		fmt.Println(err.Error())
		return
	}

	buf, err := ioutil.ReadAll(f)

	if(err != nil) {
		fmt.Println(err.Error())
		return
	}

	g := NewJWriterMan()
	g2 := NewJParse(buf, g)
	g2.Format()

	tree := NewJTree(g.Lines)
	
	tree.isStart(0)

	term := NewJTerm(tree)
	//printChars(g.Lines)
	term.ListenAction()
}

func printChars(lines []Line) {
	for _, v1 := range lines {
		for _, v2 := range v1 {
			fmt.Printf("%c", v2.Val)
		}
		fmt.Println()
	}
}
package main

import (
	"fmt"

	"github.com/peta-skr/NAND2Tetris/Chapter7/parser"
)

var Constant [32768]int

func main() {
	parseData, err := parser.Constructor("./test/SimpleAdd/SimpleAdd.vm")

	if err != nil {
		fmt.Println("some Error")
		return
	}
}
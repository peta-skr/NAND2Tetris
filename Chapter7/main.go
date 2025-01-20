package main

import (
	"fmt"
	"strconv"

	codewriter "github.com/peta-skr/NAND2Tetris/Chapter7/codeWriter"
	"github.com/peta-skr/NAND2Tetris/Chapter7/parser"
)

// var Constant [32768]int
var Stack = make([]int, 10)

func main() {
	parseData, err := parser.Constructor("./test/SimpleAdd/SimpleAdd.vm")
	output := codewriter.Constructor("./test/SimpleAdd/SimpleAdd.asm")

	if err != nil {
		fmt.Println("some Error")
		return
	}

	for parseData.HasMoreCommands() {
		parseData.Advance()

		switch parseData.CommandType() {
		case parser.C_ARITHMETIC:
			command := parseData.Arg1()
			if command == "add" {
				num1 := Stack[len(Stack) - 1]
				num2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]
				add := num1 + num2
				Stack = append(Stack, add)
			}
		case parser.C_PUSH:
			str := parseData.Arg2()
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("some Error")
				return
			}
			Stack = append(Stack, num)
		}
	}
}
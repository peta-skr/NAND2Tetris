package main

import (
	"fmt"
	"strconv"

	codewriter "github.com/peta-skr/NAND2Tetris/Chapter7/codeWriter"
	"github.com/peta-skr/NAND2Tetris/Chapter7/parser"
)

// var Constant [32768]int
var Stack = make([]int, 0)

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
				// num1 := Stack[len(Stack) - 1]
				// num2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]

				// str := "Add " + strconv.Itoa(num1) + " " + strconv.Itoa(num2)
				
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D+M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			}
		case parser.C_PUSH:
			str := parseData.Arg2()
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("some Error")
				return
			}
			str1 := "@" + strconv.Itoa(len(Stack) + 256)
			output.WriteArithmetic(str1)
			Stack = append(Stack, num)
			output.WriteArithmetic("D=A")
			output.WriteArithmetic("@SP")
			output.WriteArithmetic("M=D")
			output.WriteArithmetic("@"+str)
			output.WriteArithmetic("D=A")
			output.WriteArithmetic(str1)
			output.WriteArithmetic("M=D")
			output.WriteArithmetic("@SP")
			output.WriteArithmetic("M=M+1")
		}
	}
}
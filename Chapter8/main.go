package main

import (
	"fmt"

	codewriter "github.com/peta-skr/NAND2Tetris/Chapter7/codeWriter"
	"github.com/peta-skr/NAND2Tetris/Chapter7/parser"
)

// var Constant [32768]int
var Stack = make([]int, 0)
var Local = make([]int, 0)

func main() {
	vm("./test/FibonacciSeries/FibonacciSeries.vm", "./test/FibonacciSeries/FibonacciSeries.asm")
}

func vm(inputfile string, outputfile string) {
	parseData, err := parser.Constructor(inputfile)
	output := codewriter.Constructor(outputfile)

	if err != nil {
		fmt.Println("some Error")
		return
	}

	output.WriteArithmetic("@256")
	output.WriteArithmetic("D=A")
	output.WriteArithmetic("@SP")
	output.WriteArithmetic("M=D")

	output.WriteArithmetic("@300")
	output.WriteArithmetic("D=A")
	output.WriteArithmetic("@LCL")
	output.WriteArithmetic("M=D")

	output.WriteArithmetic("@400")
	output.WriteArithmetic("D=A")
	output.WriteArithmetic("@ARG")
	output.WriteArithmetic("M=D")

	output.WriteArithmetic("@3000")
	output.WriteArithmetic("D=A")
	output.WriteArithmetic("@THIS")
	output.WriteArithmetic("M=D")

	output.WriteArithmetic("@3010")
	output.WriteArithmetic("D=A")
	output.WriteArithmetic("@THAT")
	output.WriteArithmetic("M=D")
	fmt.Println(parseData)

	for parseData.HasMoreCommands() {
		parseData.Advance()

		switch parseData.CommandType() {
		case parser.C_ARITHMETIC:
			command := parseData.Arg1()
			switch command {
			case "add":
				output.WriteArithmetic("add")
			case "sub":
				output.WriteArithmetic("sub")
			case "eq":
				output.WriteArithmetic("eq")
			case "lt":
				output.WriteArithmetic("lt")
			case "gt":
				output.WriteArithmetic("gt")
			case "neg":
				output.WriteArithmetic("neg")
			case "and":
				output.WriteArithmetic("and")
			case "or":
				output.WriteArithmetic("or")
			case "not":
				output.WriteArithmetic("not")
			}
		case parser.C_PUSH:
			switch parseData.Arg1() {
			case "constant":
				output.WritePushPop(parser.C_PUSH, "constant", parseData.Arg2())
			case "local":
				output.WritePushPop(parser.C_PUSH, "local", parseData.Arg2())
			case "argument":
				output.WritePushPop(parser.C_PUSH, "argument", parseData.Arg2())
			case "this":
				output.WritePushPop(parser.C_PUSH, "this", parseData.Arg2())
			case "that":
				output.WritePushPop(parser.C_PUSH, "that", parseData.Arg2())
			case "temp":
				output.WritePushPop(parser.C_PUSH, "temp", parseData.Arg2())
			case "pointer":
				output.WritePushPop(parser.C_PUSH, "pointer", parseData.Arg2())
			case "static":
				output.WritePushPop(parser.C_PUSH, "static", parseData.Arg2())
			}
		case parser.C_POP:
			switch parseData.Arg1() {
			case "constant":
				output.WritePushPop(parser.C_POP, "constant", parseData.Arg2())

			case "local":
				output.WritePushPop(parser.C_POP, "local", parseData.Arg2())
			case "argument":
				output.WritePushPop(parser.C_POP, "argument", parseData.Arg2())
			case "this":
				output.WritePushPop(parser.C_POP, "this", parseData.Arg2())
			case "that":
				output.WritePushPop(parser.C_POP, "that", parseData.Arg2())
			case "temp":
				output.WritePushPop(parser.C_POP, "temp", parseData.Arg2())
			case "pointer":
				output.WritePushPop(parser.C_POP, "pointer", parseData.Arg2())
			case "static":
				output.WritePushPop(parser.C_POP, "static", parseData.Arg2())
			}
		case parser.C_LABEL:
			output.WriteLabel(parseData.Arg1())
		case parser.C_IF:
			output.WriteIf(parseData.Arg1())
		case parser.C_GOTO:
			output.WriteGoto(parseData.Arg1())
		}
	}
}

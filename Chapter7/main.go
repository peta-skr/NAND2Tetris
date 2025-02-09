package main

import (
	"fmt"
	"math/rand"
	"strconv"

	codewriter "github.com/peta-skr/NAND2Tetris/Chapter7/codeWriter"
	"github.com/peta-skr/NAND2Tetris/Chapter7/parser"
)

// var Constant [32768]int
var Stack = make([]int, 0)
var Local = make([]int, 0)

func main() {
	vm("./test/BasicTest/BasicTest.vm", "./test/BasicTest/BasicTest.asm")
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
				str := parseData.Arg2()
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("some Error")
					return
				}
				Stack = append(Stack, num)
				output.WriteArithmetic("@" + str)
				output.WriteArithmetic("D=A")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			case "local":
				str := parseData.Arg2()
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("some Error")
					return
				}
				Stack = append(Stack, num)
				output.WriteArithmetic("@LCL")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@" + str)
				output.WriteArithmetic("D=D+A")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=0")
			case "argument":
				str := parseData.Arg2()
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("some Error")
					return
				}
				Stack = append(Stack, num)
				output.WriteArithmetic("@ARG")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@" + str)
				output.WriteArithmetic("D=D+A")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=0")
			case "this":
				str := parseData.Arg2()
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("some Error")
					return
				}
				Stack = append(Stack, num)
				output.WriteArithmetic("@THIS")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@" + str)
				output.WriteArithmetic("D=D+A")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=0")
			case "that":
				str := parseData.Arg2()
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("some Error")
					return
				}
				Stack = append(Stack, num)
				output.WriteArithmetic("@THAT")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@" + str)
				output.WriteArithmetic("D=D+A")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=0")
			case "temp":
				str := parseData.Arg2()
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("some Error")
					return
				}
				Stack = append(Stack, num)
				output.WriteArithmetic("@R5")
				output.WriteArithmetic("D=A")
				output.WriteArithmetic("@" + str)
				output.WriteArithmetic("D=D+A")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=0")
			case "pointer":
				str := parseData.Arg2()
				num, err := strconv.Atoi(str)
				if err != nil {
					fmt.Println("some Error")
					return
				}

				if num == 0 {
					Stack = append(Stack, num)
					output.WriteArithmetic("@SP")
					output.WriteArithmetic("M=M-1")
					output.WriteArithmetic("@SP")
					output.WriteArithmetic("A=M")
					output.WriteArithmetic("D=M")
					output.WriteArithmetic("M=0")
					output.WriteArithmetic("@THIS")
					output.WriteArithmetic("M=D")
				} else {
					output.WriteArithmetic("@SP")
					output.WriteArithmetic("M=M-1")
					output.WriteArithmetic("@SP")
					output.WriteArithmetic("A=M")
					output.WriteArithmetic("D=M")
					output.WriteArithmetic("M=0")
					output.WriteArithmetic("@THAT")
					output.WriteArithmetic("M=D")
				}
			case "static":
				str := parseData.Arg2()

				output.WriteArithmetic("@16")
				output.WriteArithmetic("D=A")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@" + str)
				output.WriteArithmetic("D=A")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=D+M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				output.WriteArithmetic("@R13")
				output.WriteArithmetic("M=0")
			}
		}
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// https://qiita.com/srtkkou/items/ccbddc881d6f3549baf1
func RandomString() string {

	// r := rand.New(time.Now().UnixNano())

	b := make([]rune, 20)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

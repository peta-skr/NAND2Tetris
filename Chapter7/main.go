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
	parseData, err := parser.Constructor("./test/StaticTest/StaticTest.vm")
	output := codewriter.Constructor("./test/StaticTest/StaticTest.asm")

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
				s1 := Stack[len(Stack) - 1]
				s2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]

				Stack = append(Stack, s1 + s2)

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
			case "sub":
				s1 := Stack[len(Stack) - 1]
				s2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]

				Stack = append(Stack, s1 - s2)

				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=M-D")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
				
			case "eq":
				s1 := Stack[len(Stack) - 1]
				s2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]

				if s1 == s2 {
					Stack = append(Stack, 0)
				}else {
					Stack = append(Stack, -1)
				}

				label1 := RandomString()
				label2 := RandomString()


				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M-D")
				output.WriteArithmetic("@"+label1)
				output.WriteArithmetic("D;JEQ")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@"+label2)
				output.WriteArithmetic("0;JMP")
				output.WriteArithmetic("(" + label1 + ")")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=-1")
				output.WriteArithmetic("(" + label2 + ")")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			case "lt":
				s1 := Stack[len(Stack) - 1]
				s2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]

				if s2 < s1 {
					Stack = append(Stack, 0)
				}else {
					Stack = append(Stack, -1)
				}

				label1 := RandomString()
				label2 := RandomString()


				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M-D")
				output.WriteArithmetic("@"+label1)
				output.WriteArithmetic("D;JLT")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@"+label2)
				output.WriteArithmetic("0;JMP")
				output.WriteArithmetic("(" + label1 + ")")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=-1")
				output.WriteArithmetic("(" + label2 + ")")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			case "gt":
				s1 := Stack[len(Stack) - 1]
				s2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]

				if s2 > s1 {
					Stack = append(Stack, 0)
				}else {
					Stack = append(Stack, -1)
				}

				label1 := RandomString()
				label2 := RandomString()


				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M-D")
				output.WriteArithmetic("@"+label1)
				output.WriteArithmetic("D;JGT")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=0")
				output.WriteArithmetic("@"+label2)
				output.WriteArithmetic("0;JMP")
				output.WriteArithmetic("(" + label1 + ")")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=-1")
				output.WriteArithmetic("(" + label2 + ")")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			case "neg":
				s1 := Stack[len(Stack) - 1]
				Stack = Stack[:len(Stack) - 1]

				Stack = append(Stack, -s1)

				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=-M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			case "and":
				s1 := Stack[len(Stack) - 1]
				s2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]

				Stack = append(Stack, s1 & s2)

				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D&M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			case "or":
				s1 := Stack[len(Stack) - 1]
				s2 := Stack[len(Stack) - 2]
				Stack = Stack[:len(Stack) - 2]

				Stack = append(Stack, s1 | s2)

				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D|M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			case "not":
				s1 := Stack[len(Stack) - 1]
				Stack = Stack[:len(Stack) - 1]

				Stack = append(Stack, ^s1)

				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M-1")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=!M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			}
		case parser.C_PUSH:
			switch parseData.Arg1() {
			case "constant":
				str := parseData.Arg2()
				num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("some Error")
				return
			}
			Stack = append(Stack, num)
			output.WriteArithmetic("@"+str)
			output.WriteArithmetic("D=A")
			output.WriteArithmetic("@SP")
			output.WriteArithmetic("A=M")
			output.WriteArithmetic("M=D")
			output.WriteArithmetic("@SP")
			output.WriteArithmetic("M=M+1")
		case "local":
			str := parseData.Arg2()
		// 	num, err := strconv.Atoi(str)
		// if err != nil {
		// 	fmt.Println("some Error")
		// 	return
		// }
		output.WriteArithmetic("@LCL")
		output.WriteArithmetic("D=M")
		output.WriteArithmetic("@" + str)
		output.WriteArithmetic("D=D+A")
		output.WriteArithmetic("A=D")
		output.WriteArithmetic("D=M")
		output.WriteArithmetic("@SP")
		output.WriteArithmetic("A=M")
		output.WriteArithmetic("M=D")
		output.WriteArithmetic("@LCL")
		output.WriteArithmetic("D=M")
		// output.WriteArithmetic("@" + str)
		// output.WriteArithmetic("D=D+A")
		// output.WriteArithmetic("A=D")
		// output.WriteArithmetic("M=0")
		output.WriteArithmetic("@SP")
		output.WriteArithmetic("M=M+1")
		case "argument":
			str := parseData.Arg2()
		// 	num, err := strconv.Atoi(str)
		// if err != nil {
		// 	fmt.Println("some Error")
		// 	return
		// }
		output.WriteArithmetic("@ARG")
		output.WriteArithmetic("D=M")
		output.WriteArithmetic("@" + str)
		output.WriteArithmetic("D=D+A")
		output.WriteArithmetic("A=D")
		output.WriteArithmetic("D=M")
		output.WriteArithmetic("@SP")
		output.WriteArithmetic("A=M")
		output.WriteArithmetic("M=D")
		output.WriteArithmetic("@ARG")
		output.WriteArithmetic("D=M")
		// output.WriteArithmetic("@" + str)
		// output.WriteArithmetic("D=D+A")
		// output.WriteArithmetic("A=D")
		// output.WriteArithmetic("M=0")
		output.WriteArithmetic("@SP")
		output.WriteArithmetic("M=M+1")
		case "this":
			str := parseData.Arg2()
		// 	num, err := strconv.Atoi(str)
		// if err != nil {
		// 	fmt.Println("some Error")
		// 	return
		// }
		output.WriteArithmetic("@THIS")
		output.WriteArithmetic("D=M")
		output.WriteArithmetic("@" + str)
		output.WriteArithmetic("D=D+A")
		output.WriteArithmetic("A=D")
		output.WriteArithmetic("D=M")
		output.WriteArithmetic("@SP")
		output.WriteArithmetic("A=M")
		output.WriteArithmetic("M=D")
		output.WriteArithmetic("@THIS")
		output.WriteArithmetic("D=M")
		// output.WriteArithmetic("@" + str)
		// output.WriteArithmetic("D=D+A")
		// output.WriteArithmetic("A=D")
		// output.WriteArithmetic("M=0")
		output.WriteArithmetic("@SP")
		output.WriteArithmetic("M=M+1")
		case "that":
			str := parseData.Arg2()
		// 	num, err := strconv.Atoi(str)
		// if err != nil {
		// 	fmt.Println("some Error")
		// 	return
		// }
		output.WriteArithmetic("@THAT")
		output.WriteArithmetic("D=M")
		output.WriteArithmetic("@" + str)
		output.WriteArithmetic("D=D+A")
		output.WriteArithmetic("A=D")
		output.WriteArithmetic("D=M")
		output.WriteArithmetic("@SP")
		output.WriteArithmetic("A=M")
		output.WriteArithmetic("M=D")
		output.WriteArithmetic("@THAT")
		output.WriteArithmetic("D=M")
		// output.WriteArithmetic("@" + str)
		// output.WriteArithmetic("D=D+A")
		// output.WriteArithmetic("A=D")
		// output.WriteArithmetic("M=0")
		output.WriteArithmetic("@SP")
		output.WriteArithmetic("M=M+1")
		case "temp":
			str := parseData.Arg2()
			// 	num, err := strconv.Atoi(str)
			// if err != nil {
			// 	fmt.Println("some Error")
			// 	return
			// }
			output.WriteArithmetic("@R5")
			output.WriteArithmetic("D=A")
			output.WriteArithmetic("@" + str)
			output.WriteArithmetic("D=D+A")
			output.WriteArithmetic("A=D")
			output.WriteArithmetic("D=M")
			output.WriteArithmetic("@SP")
			output.WriteArithmetic("A=M")
			output.WriteArithmetic("M=D")
			output.WriteArithmetic("@R5")
			output.WriteArithmetic("D=M")
			// output.WriteArithmetic("@" + str)
			// output.WriteArithmetic("D=D+A")
			// output.WriteArithmetic("A=D")
			// output.WriteArithmetic("M=0")
			output.WriteArithmetic("@SP")
			output.WriteArithmetic("M=M+1")
		case "pointer":
			str := parseData.Arg2()
				num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Println("some Error")
				return
			}

			if num == 0 {
				output.WriteArithmetic("@THIS")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				// output.WriteArithmetic("@" + str)
				// output.WriteArithmetic("D=D+A")
				// output.WriteArithmetic("A=D")
				// output.WriteArithmetic("M=0")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			}else {
				output.WriteArithmetic("@THAT")
				output.WriteArithmetic("D=M")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("A=M")
				output.WriteArithmetic("M=D")
				// output.WriteArithmetic("@" + str)
				// output.WriteArithmetic("D=D+A")
				// output.WriteArithmetic("A=D")
				// output.WriteArithmetic("M=0")
				output.WriteArithmetic("@SP")
				output.WriteArithmetic("M=M+1")
			}
		case "static":

			str := parseData.Arg2()

			output.WriteArithmetic("@16")
			output.WriteArithmetic("D=A")
			output.WriteArithmetic("@" + str)
			output.WriteArithmetic("A=D+A")
			output.WriteArithmetic("D=M")
			output.WriteArithmetic("@SP")
			output.WriteArithmetic("A=M")
			output.WriteArithmetic("M=D")
			output.WriteArithmetic("@SP")
			output.WriteArithmetic("M=M+1")
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
				output.WriteArithmetic("@"+str)
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
				}else {
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
				output.WriteArithmetic("@"+str)
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
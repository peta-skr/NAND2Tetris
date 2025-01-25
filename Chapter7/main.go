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

func main() {
	parseData, err := parser.Constructor("./test/StackTest/StackTest.vm")
	output := codewriter.Constructor("./test/StackTest/StackTest.asm")

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
				Stack = Stack[:len(Stack) - 2]

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
			}else if command == "eq" {
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
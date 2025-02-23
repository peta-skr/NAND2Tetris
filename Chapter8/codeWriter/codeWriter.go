package codewriter

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/peta-skr/NAND2Tetris/Chapter7/parser"
)

type Output struct {
	file     *os.File
	filename string
}

func Constructor(filepath string) Output {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		fmt.Println("ファイルを開けませんでした:", err)
		return Output{}
	}

	output := Output{
		file: file,
	}

	return output
}

func (o *Output) SetFileName(filename string) {
	o.filename = filename
}

func (o *Output) WriteArithmetic(command string) {

	switch command {
	case "add":
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("M=D+M")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	case "sub":
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("M=M-D")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")

	case "eq":
		label1 := RandomString()
		label2 := RandomString()

		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M-D")
		o.WriteCode("@" + label1)
		o.WriteCode("D;JEQ")
		o.WriteCode("@SP")
		o.WriteCode("M=M")
		o.WriteCode("A=M")
		o.WriteCode("M=0")
		o.WriteCode("@" + label2)
		o.WriteCode("0;JMP")
		o.WriteCode("(" + label1 + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M")
		o.WriteCode("A=M")
		o.WriteCode("M=-1")
		o.WriteCode("(" + label2 + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	case "lt":
		label1 := RandomString()
		label2 := RandomString()

		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M-D")
		o.WriteCode("@" + label1)
		o.WriteCode("D;JLT")
		o.WriteCode("@SP")
		o.WriteCode("M=M")
		o.WriteCode("A=M")
		o.WriteCode("M=0")
		o.WriteCode("@" + label2)
		o.WriteCode("0;JMP")
		o.WriteCode("(" + label1 + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M")
		o.WriteCode("A=M")
		o.WriteCode("M=-1")
		o.WriteCode("(" + label2 + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	case "gt":

		label1 := RandomString()
		label2 := RandomString()

		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M-D")
		o.WriteCode("@" + label1)
		o.WriteCode("D;JGT")
		o.WriteCode("@SP")
		o.WriteCode("M=M")
		o.WriteCode("A=M")
		o.WriteCode("M=0")
		o.WriteCode("@" + label2)
		o.WriteCode("0;JMP")
		o.WriteCode("(" + label1 + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M")
		o.WriteCode("A=M")
		o.WriteCode("M=-1")
		o.WriteCode("(" + label2 + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	case "neg":
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("M=-M")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	case "and":

		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("M=D&M")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	case "or":
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("M=D|M")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	case "not":

		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("M=!M")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	default:
		o.WriteCode(command)
	}

}

func (o *Output) WriteCode(command string) {
	_, _ = o.file.WriteString(command + "\n")
}

func (o *Output) WritePushPop(cmdType parser.CmdType, command string, arg2 string) {
	switch cmdType {
	case parser.C_PUSH:
		switch command {
		case "constant":
			o.WriteCode("@" + arg2)
			o.WriteCode("D=A")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@SP")
			o.WriteCode("M=M+1")
		case "local":
			o.WriteCode("@LCL")
			o.WriteCode("D=M")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("A=D")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@LCL")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("M=M+1")
		case "argument":
			o.WriteCode("@ARG")
			o.WriteCode("D=M")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("A=D")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@ARG")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("M=M+1")
		case "this":
			o.WriteCode("@THIS")
			o.WriteCode("D=M")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("A=D")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@THIS")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("M=M+1")

		case "that":
			o.WriteCode("@THAT")
			o.WriteCode("D=M")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("A=D")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@THAT")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("M=M+1")
		case "temp":
			o.WriteCode("@R5")
			o.WriteCode("D=A")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("A=D")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@R5")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("M=M+1")
		case "pointer":

			if arg2 == "0" {
				o.WriteCode("@THIS")
				o.WriteCode("D=M")
				o.WriteCode("@SP")
				o.WriteCode("A=M")
				o.WriteCode("M=D")
				o.WriteCode("@SP")
				o.WriteCode("M=M+1")
			} else {
				o.WriteCode("@THAT")
				o.WriteCode("D=M")
				o.WriteCode("@SP")
				o.WriteCode("A=M")
				o.WriteCode("M=D")
				o.WriteCode("@SP")
				o.WriteCode("M=M+1")
			}
		case "static":
			o.WriteCode("@16")
			o.WriteCode("D=A")
			o.WriteCode("@" + arg2)
			o.WriteCode("A=D+A")
			o.WriteCode("D=M")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@SP")
			o.WriteCode("M=M+1")
		}
	case parser.C_POP:
		switch command {

		// case "constant":
		// 	o.WriteCode("@" + arg2)
		// 	o.WriteCode("D=A")
		// 	o.WriteCode("@SP")
		// 	o.WriteCode("A=M")
		// 	o.WriteCode("M=D")
		// 	o.WriteCode("@SP")
		// 	o.WriteCode("M=M+1")
		case "local":
			o.WriteCode("@LCL")
			o.WriteCode("D=M")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("@R13")
			o.WriteCode("M=D")
			o.WriteCode("@SP")
			o.WriteCode("M=M-1")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("D=M")
			o.WriteCode("M=0")
			o.WriteCode("@R13")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@R13")
			o.WriteCode("M=0")
		case "argument":
			o.WriteCode("@ARG")
			o.WriteCode("D=M")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("@R13")
			o.WriteCode("M=D")
			o.WriteCode("@SP")
			o.WriteCode("M=M-1")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("D=M")
			o.WriteCode("M=0")
			o.WriteCode("@R13")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@R13")
			o.WriteCode("M=0")
		case "this":
			o.WriteCode("@THIS")
			o.WriteCode("D=M")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("@R13")
			o.WriteCode("M=D")
			o.WriteCode("@SP")
			o.WriteCode("M=M-1")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("D=M")
			o.WriteCode("M=0")
			o.WriteCode("@R13")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@R13")
			o.WriteCode("M=0")
		case "that":
			o.WriteCode("@THAT")
			o.WriteCode("D=M")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("@R13")
			o.WriteCode("M=D")
			o.WriteCode("@SP")
			o.WriteCode("M=M-1")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("D=M")
			o.WriteCode("M=0")
			o.WriteCode("@R13")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@R13")
			o.WriteCode("M=0")
		case "temp":
			o.WriteCode("@R5")
			o.WriteCode("D=A")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=D+A")
			o.WriteCode("@R13")
			o.WriteCode("M=D")
			o.WriteCode("@SP")
			o.WriteCode("M=M-1")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("D=M")
			o.WriteCode("M=0")
			o.WriteCode("@R13")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@R13")
			o.WriteCode("M=0")
		case "pointer":
			if arg2 == "0" {
				o.WriteCode("@SP")
				o.WriteCode("M=M-1")
				o.WriteCode("@SP")
				o.WriteCode("A=M")
				o.WriteCode("D=M")
				o.WriteCode("M=0")
				o.WriteCode("@THIS")
				o.WriteCode("M=D")
			} else {
				o.WriteCode("@SP")
				o.WriteCode("M=M-1")
				o.WriteCode("@SP")
				o.WriteCode("A=M")
				o.WriteCode("D=M")
				o.WriteCode("M=0")
				o.WriteCode("@THAT")
				o.WriteCode("M=D")
			}
		case "static":
			o.WriteCode("@16")
			o.WriteCode("D=A")
			o.WriteCode("@R13")
			o.WriteCode("M=D")
			o.WriteCode("@" + arg2)
			o.WriteCode("D=A")
			o.WriteCode("@R13")
			o.WriteCode("M=D+M")
			o.WriteCode("@SP")
			o.WriteCode("M=M-1")
			o.WriteCode("@SP")
			o.WriteCode("A=M")
			o.WriteCode("D=M")
			o.WriteCode("M=0")
			o.WriteCode("@R13")
			o.WriteCode("A=M")
			o.WriteCode("M=D")
			o.WriteCode("@R13")
			o.WriteCode("M=0")
		}
	}
}

func (o *Output) close() {
	o.file.Close()

}

func (o *Output) WriteInit() {}

func (o *Output) WriteLabel(label string) {
	o.WriteCode("(" + label + ")")
}

func (o *Output) WriteGoto(label string) {
	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode("D=M")
	o.WriteCode("@" + label)
	o.WriteCode("0;JMP")
}

func (o *Output) WriteIf(label string) {

	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode("D=M")
	o.WriteCode("@" + label)
	o.WriteCode("D;JNE")

}

func (o *Output) WriteCall() {

}

func (o *Output) WriteReturn() {

}

func (o *Output) WriteFunction() {

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

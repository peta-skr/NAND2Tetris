package codewriter

import (
	"fmt"
	"math/rand"
	"os"
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
	}

}

func (o *Output) WriteCode(command string) {
	_, _ = o.file.WriteString(command + "\n")
}

func writePushPop() {

}

func (o *Output) close() {
	o.file.Close()

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

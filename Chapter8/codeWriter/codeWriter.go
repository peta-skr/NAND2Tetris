package codewriter

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/peta-skr/NAND2Tetris/Chapter7/parser"
)

type Output struct {
	file            *os.File
	filename        string
	isInFunction    bool
	currentFunction string
}

func Constructor(filepath string) Output {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) // ファイルが存在しなければ新規作成＋書き込み専用

	if err != nil {
		fmt.Println("ファイルを開けませんでした:", err)
		return Output{}
	}

	return Output{file: file}
}

func (o *Output) SetFileName(filename string) {
	o.filename = filename
}

func (o *Output) WriteArithmetic(command string) {
	switch command {
	case "add", "sub", "and", "or":
		o.binaryOperation(command)
	case "neg", "not":
		o.unaryOperation(command)
	case "eq", "lt", "gt":
		o.comparisonOperation(command)
	default:
		o.WriteCode(command)
	}
}

func (o *Output) binaryOperation(command string) {
	op := map[string]string{
		"add": "M=D+M",
		"sub": "M=M-D",
		"and": "M=D&M",
		"or":  "M=D|M",
	}[command]

	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode("D=M")
	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode(op)
	o.WriteCode("@SP")
	o.WriteCode("M=M+1")
}

func (o *Output) unaryOperation(command string) {
	op := map[string]string{
		"neg": "M=-M",
		"not": "M=!M",
	}[command]

	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode(op)
	o.WriteCode("@SP")
	o.WriteCode("M=M+1")
}

func (o *Output) comparisonOperation(command string) {
	labelTrue := "TRUE_" + RandomString()
	labelEnd := "END_" + RandomString()
	jump := map[string]string{
		"eq": "JEQ",
		"lt": "JLT",
		"gt": "JGT",
	}[command]

	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode("D=M")
	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode("D=M-D")
	o.WriteCode("@" + labelTrue)
	o.WriteCode("D;" + jump)
	o.WriteCode("@SP")
	o.WriteCode("A=M")
	o.WriteCode("M=0")
	o.WriteCode("@" + labelEnd)
	o.WriteCode("0;JMP")
	o.WriteCode("(" + labelTrue + ")")
	o.WriteCode("@SP")
	o.WriteCode("A=M")
	o.WriteCode("M=-1")
	o.WriteCode("(" + labelEnd + ")")
	o.WriteCode("@SP")
	o.WriteCode("M=M+1")
}

func (o *Output) WriteCode(command string) {
	_, _ = o.file.WriteString(command + "\n")
}

func (o *Output) WritePushPop(cmdType parser.CmdType, segment, index, fileName string) {
	if cmdType == parser.C_PUSH {
		o.push(segment, index, fileName)
	} else {
		o.pop(segment, index, fileName)
	}
}

func (o *Output) push(segment, index, fileName string) {
	switch segment {
	case "constant":
		o.WriteCode("@" + index)
		o.WriteCode("D=A")
	case "local", "argument", "this", "that":
		base := map[string]string{
			"local":    "LCL",
			"argument": "ARG",
			"this":     "THIS",
			"that":     "THAT",
		}[segment]
		o.WriteCode("@" + base)
		o.WriteCode("D=M")
		o.WriteCode("@" + index)
		o.WriteCode("A=D+A")
		o.WriteCode("D=M")
	case "temp":
		o.WriteCode("@R5")
		o.WriteCode("D=A")
		o.WriteCode("@" + index)
		o.WriteCode("A=D+A")
		o.WriteCode("D=M")
	case "pointer":
		base := "THIS"
		if index == "1" {
			base = "THAT"
		}
		o.WriteCode("@" + base)
		o.WriteCode("D=M")
	case "static":
		o.WriteCode("@" + fileName + "." + index)
		o.WriteCode("D=M")
	}
	o.WriteCode("@SP")
	o.WriteCode("A=M")
	o.WriteCode("M=D")
	o.WriteCode("@SP")
	o.WriteCode("M=M+1")
}

func (o *Output) pop(segment, index, fileName string) {
	switch segment {
	case "local", "argument", "this", "that":
		base := map[string]string{
			"local":    "LCL",
			"argument": "ARG",
			"this":     "THIS",
			"that":     "THAT",
		}[segment]
		o.WriteCode("@" + base)
		o.WriteCode("D=M")
		o.WriteCode("@" + index)
		o.WriteCode("D=D+A")
	case "temp":
		o.WriteCode("@R5")
		o.WriteCode("D=A")
		o.WriteCode("@" + index)
		o.WriteCode("D=D+A")
	case "pointer":
		base := "THIS"
		if index == "1" {
			base = "THAT"
		}
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M")
		o.WriteCode("@" + base)
		o.WriteCode("M=D")
		return
	case "static":
		o.WriteCode("@" + fileName + "." + index)
		o.WriteCode("D=A")
	}
	o.WriteCode("@R13")
	o.WriteCode("M=D")
	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode("D=M")
	o.WriteCode("@R13")
	o.WriteCode("A=M")
	o.WriteCode("M=D")
}

func (o *Output) Close() {

	err := o.file.Sync()
	if err != nil {
		fmt.Println("Sync エラー:", err)
		return
	}

	defer o.file.Close()
}

func (o *Output) WriteInit() {}

func (o *Output) WriteLabel(label string) {
	var fullLabel string
	if o.isInFunction {
		fullLabel = o.currentFunction + "$" + label
	} else {
		fullLabel = label
	}
	o.WriteCode("(" + fullLabel + ")")
}

func (o *Output) WriteGoto(label string) {
	var fullLabel string
	if o.isInFunction {
		fullLabel = o.currentFunction + "$" + label
	} else {
		fullLabel = label
	}
	o.WriteCode("@" + fullLabel)
	o.WriteCode("0;JMP")
}

func (o *Output) WriteIf(label string) {
	var fullLabel string
	if o.isInFunction {
		fullLabel = o.currentFunction + "$" + label
	} else {
		fullLabel = label
	}
	o.WriteCode("@SP")
	o.WriteCode("M=M-1")
	o.WriteCode("A=M")
	o.WriteCode("D=M")
	o.WriteCode("@" + fullLabel)
	o.WriteCode("D;JNE")
}

func (o *Output) WriteCall(functionName string, numArgs string) {

	retAddrLabel := functionName + "." + RandomString()

	// 戻りアドレスをプッシュ
	o.WriteCode("@" + retAddrLabel)
	o.WriteCode("D=A")
	o.WriteCode("@SP")
	o.WriteCode("A=M")
	o.WriteCode("M=D")
	o.WriteCode("@SP")
	o.WriteCode("M=M+1")

	// LCL, ARG, THIS, THATをプッシュ
	for _, segment := range []string{"LCL", "ARG", "THIS", "THAT"} {
		o.WriteCode("@" + segment)
		o.WriteCode("D=M")
		o.WriteCode("@SP")
		o.WriteCode("A=M")
		o.WriteCode("M=D")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
	}

	// ARG = SP - 5 - nArgs
	o.WriteCode("@SP")
	o.WriteCode("D=M")
	o.WriteCode("@5")
	o.WriteCode("D=D-A")
	o.WriteCode("@" + numArgs)
	o.WriteCode("D=D-A")
	o.WriteCode("@ARG")
	o.WriteCode("M=D")

	// LCL = SP
	o.WriteCode("@SP")
	o.WriteCode("D=M")
	o.WriteCode("@LCL")
	o.WriteCode("M=D")

	// goto functionName
	o.WriteCode("@" + functionName)
	o.WriteCode("0;JMP")

	o.WriteCode("(" + retAddrLabel + ")")
}

func (o *Output) WriteReturn() {
	// フレームを設定 (LCLの値をR13に保存)
	o.WriteCode("@LCL")
	o.WriteCode("D=M")
	o.WriteCode("@R13")
	o.WriteCode("M=D")

	// 戻りアドレスを取得 (フレーム - 5)
	o.WriteCode("@5")
	o.WriteCode("A=D-A")
	o.WriteCode("D=M")
	o.WriteCode("@R14")
	o.WriteCode("M=D")

	// 戻り値をARG[0]に配置
	o.WriteCode("@SP")
	o.WriteCode("AM=M-1")
	o.WriteCode("D=M")
	o.WriteCode("@ARG")
	o.WriteCode("A=M")
	o.WriteCode("M=D")

	// SPを(ARG+1)に設定
	o.WriteCode("@ARG")
	o.WriteCode("D=M+1")
	o.WriteCode("@SP")
	o.WriteCode("M=D")

	// THATを復元
	o.WriteCode("@R13")
	o.WriteCode("AM=M-1")
	o.WriteCode("D=M")
	o.WriteCode("@THAT")
	o.WriteCode("M=D")

	// THISを復元
	o.WriteCode("@R13")
	o.WriteCode("AM=M-1")
	o.WriteCode("D=M")
	o.WriteCode("@THIS")
	o.WriteCode("M=D")

	// ARGを復元
	o.WriteCode("@R13")
	o.WriteCode("AM=M-1")
	o.WriteCode("D=M")
	o.WriteCode("@ARG")
	o.WriteCode("M=D")

	// LCLを復元
	o.WriteCode("@R13")
	o.WriteCode("AM=M-1")
	o.WriteCode("D=M")
	o.WriteCode("@LCL")
	o.WriteCode("M=D")

	// 戻りアドレスにジャンプ
	o.WriteCode("@R14")
	o.WriteCode("A=M")
	o.WriteCode("0;JMP")

	o.isInFunction = false // 関数から抜ける
	o.currentFunction = "" // 現在の関数名をリセット
}

func (o *Output) WriteFunction(functionName string, numLocal int) {
	o.currentFunction = functionName
	o.isInFunction = true
	o.WriteCode("(" + functionName + ")")

	for i := 0; i < numLocal; i++ {
		o.WriteCode("@SP")
		o.WriteCode("A=M")
		o.WriteCode("M=0")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1")
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

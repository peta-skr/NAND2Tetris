package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	codewriter "github.com/peta-skr/NAND2Tetris/Chapter7/codeWriter"
	"github.com/peta-skr/NAND2Tetris/Chapter7/parser"
)

// var Constant [32768]int
var Stack = make([]int, 0)
var Local = make([]int, 0)

func main() {
	dir := "./test/BasicLoop"
	output := "main.asm"

	// outputと同名のものがある場合は、削除する
	if _, err := os.Stat(dir + "/" + output); err == nil {
		err := os.Remove(dir + "/" + output)
		if err != nil {
			fmt.Println("ファイル削除エラー:", err)
			return
		}
	}

	// フォルダ内の.vmファイルをすべて処理する
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}

	// ファイル名を表示
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".vm") {
			vm(dir+"/"+entry.Name(), dir+"/"+output, entry.Name())
		}
	}
}

func vm(inputfile string, outputfile string, filename string) {
	parseData, err := parser.Constructor(inputfile)
	output := codewriter.Constructor(outputfile)

	if err != nil {
		fmt.Println("some Error")
		return
	}

	/*
	*  ブートストラップコードの記載
	 */
	// SP = 256
	output.WriteArithmetic("@256")
	output.WriteArithmetic("D=A")
	output.WriteArithmetic("@SP")
	output.WriteArithmetic("M=D")
	// call Sys.init
	output.WriteCall("Sys.init", "0")

	// output.WriteArithmetic("@300")
	// output.WriteArithmetic("D=A")
	// output.WriteArithmetic("@LCL")
	// output.WriteArithmetic("M=D")

	// output.WriteArithmetic("@400")
	// output.WriteArithmetic("D=A")
	// output.WriteArithmetic("@ARG")
	// output.WriteArithmetic("M=D")

	// output.WriteArithmetic("@3000")
	// output.WriteArithmetic("D=A")
	// output.WriteArithmetic("@THIS")
	// output.WriteArithmetic("M=D")

	// output.WriteArithmetic("@3010")
	// output.WriteArithmetic("D=A")
	// output.WriteArithmetic("@THAT")
	// output.WriteArithmetic("M=D")
	// fmt.Println(parseData)

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
				output.WritePushPop(parser.C_PUSH, "constant", parseData.Arg2(), filename)
			case "local":
				output.WritePushPop(parser.C_PUSH, "local", parseData.Arg2(), filename)
			case "argument":
				output.WritePushPop(parser.C_PUSH, "argument", parseData.Arg2(), filename)
			case "this":
				output.WritePushPop(parser.C_PUSH, "this", parseData.Arg2(), filename)
			case "that":
				output.WritePushPop(parser.C_PUSH, "that", parseData.Arg2(), filename)
			case "temp":
				output.WritePushPop(parser.C_PUSH, "temp", parseData.Arg2(), filename)
			case "pointer":
				output.WritePushPop(parser.C_PUSH, "pointer", parseData.Arg2(), filename)
			case "static":
				output.WritePushPop(parser.C_PUSH, "static", parseData.Arg2(), filename)
			}
		case parser.C_POP:
			switch parseData.Arg1() {
			case "constant":
				output.WritePushPop(parser.C_POP, "constant", parseData.Arg2(), filename)

			case "local":
				output.WritePushPop(parser.C_POP, "local", parseData.Arg2(), filename)
			case "argument":
				output.WritePushPop(parser.C_POP, "argument", parseData.Arg2(), filename)
			case "this":
				output.WritePushPop(parser.C_POP, "this", parseData.Arg2(), filename)
			case "that":
				output.WritePushPop(parser.C_POP, "that", parseData.Arg2(), filename)
			case "temp":
				output.WritePushPop(parser.C_POP, "temp", parseData.Arg2(), filename)
			case "pointer":
				output.WritePushPop(parser.C_POP, "pointer", parseData.Arg2(), filename)
			case "static":
				output.WritePushPop(parser.C_POP, "static", parseData.Arg2(), filename)
			}
		case parser.C_LABEL:
			output.WriteLabel(parseData.Arg1())
		case parser.C_IF:
			output.WriteIf(parseData.Arg1())
		case parser.C_GOTO:
			output.WriteGoto(parseData.Arg1())
		case parser.C_CALL:
			output.WriteCall(parseData.Arg1(), parseData.Arg2())
		case parser.C_FUNCTION:
			arg2, err := strconv.Atoi(parseData.Arg2())
			if err != nil {
				fmt.Println(err)
				return
			}
			output.WriteFunction(parseData.Arg1(), arg2)
		case parser.C_RETURN:
			output.WriteReturn()
		}
	}

	output.WriteArithmetic("(END)")
	output.WriteArithmetic("@END")
	output.WriteArithmetic("0;JMP")
	// ファイルを閉じる
	output.Close()
}

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
		// label1 := RandomString()
		// label2 := RandomString()

		// o.WriteCode("@SP")
		// o.WriteCode("M=M-1")
		// o.WriteCode("A=M")
		// o.WriteCode("D=M")
		// o.WriteCode("@SP")
		// o.WriteCode("M=M-1")
		// o.WriteCode("A=M")
		// o.WriteCode("D=M-D")
		// o.WriteCode("@" + label1)
		// o.WriteCode("D;JLT")
		// o.WriteCode("@SP")
		// o.WriteCode("M=M")
		// o.WriteCode("A=M")
		// o.WriteCode("M=0")
		// o.WriteCode("@" + label2)
		// o.WriteCode("0;JMP")
		// o.WriteCode("(" + label1 + ")")
		// o.WriteCode("@SP")
		// o.WriteCode("M=M")
		// o.WriteCode("A=M")
		// o.WriteCode("M=-1")
		// o.WriteCode("(" + label2 + ")")
		// o.WriteCode("@SP")
		// o.WriteCode("M=M+1")

		// より安全な比較演算(lt)の実装
		// ランダムなユニークラベルを生成
		label_y_neg := "Y_NEG_" + RandomString()
		label_y_pos := "Y_POS_" + RandomString()
		label_x_neg_y_neg := "X_NEG_Y_NEG_" + RandomString()
		label_x_pos_y_pos := "X_POS_Y_POS_" + RandomString()
		label_push_true := "PUSH_TRUE_" + RandomString()
		label_push_false := "PUSH_FALSE_" + RandomString()
		label_end := "END_" + RandomString()

		// より安全な比較演算(lt)の実装
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M") // y値を取得
		o.WriteCode("@" + label_y_neg)
		o.WriteCode("D;JLT") // yが負数の場合分岐
		o.WriteCode("@" + label_y_pos)
		o.WriteCode("0;JMP") // yが正数の場合分岐

		o.WriteCode("(" + label_y_neg + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M") // x値を取得
		o.WriteCode("@" + label_x_neg_y_neg)
		o.WriteCode("D;JLT") // xも負数の場合
		// xは正、yは負 → x > y は常に真なので、x < y は常に偽(0)
		o.WriteCode("@SP")
		o.WriteCode("A=M")
		o.WriteCode("M=0") // falseをプッシュ
		o.WriteCode("@" + label_end)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_x_neg_y_neg + ")")
		// 両方負の場合、通常の比較で問題なし（オーバーフローしない）
		o.WriteCode("@SP")
		o.WriteCode("M=M+1") // スタックを戻す
		o.WriteCode("A=M")
		o.WriteCode("D=M") // yを再度取得
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")   // xにポイント
		o.WriteCode("D=M-D") // x-y
		o.WriteCode("@" + label_push_true)
		o.WriteCode("D;JLT")
		o.WriteCode("@" + label_push_false)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_y_pos + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M") // x値を取得
		o.WriteCode("@" + label_x_pos_y_pos)
		o.WriteCode("D;JGE") // xが正または0の場合
		// xが負、yが正 → x < y は常に真
		o.WriteCode("@" + label_push_true)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_x_pos_y_pos + ")")
		// 両方正の場合、通常の比較で問題なし
		o.WriteCode("@SP")
		o.WriteCode("M=M+1") // スタックを戻す
		o.WriteCode("A=M")
		o.WriteCode("D=M") // yを再度取得
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")   // xにポイント
		o.WriteCode("D=M-D") // x-y
		o.WriteCode("@" + label_push_true)
		o.WriteCode("D;JLT")
		o.WriteCode("@" + label_push_false)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_push_true + ")")
		o.WriteCode("@SP")
		o.WriteCode("A=M")
		o.WriteCode("M=-1") // trueをプッシュ (-1)
		o.WriteCode("@" + label_end)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_push_false + ")")
		o.WriteCode("@SP")
		o.WriteCode("A=M")
		o.WriteCode("M=0") // falseをプッシュ (0)

		o.WriteCode("(" + label_end + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1") // スタックポインタを増加
	case "gt":

		// label1 := RandomString()
		// label2 := RandomString()

		// o.WriteCode("@SP")
		// o.WriteCode("M=M-1")
		// o.WriteCode("A=M")
		// o.WriteCode("D=M")
		// o.WriteCode("@SP")
		// o.WriteCode("M=M-1")
		// o.WriteCode("A=M")
		// o.WriteCode("D=M-D")
		// o.WriteCode("@" + label1)
		// o.WriteCode("D;JGT")
		// o.WriteCode("@SP")
		// o.WriteCode("M=M")
		// o.WriteCode("A=M")
		// o.WriteCode("M=0")
		// o.WriteCode("@" + label2)
		// o.WriteCode("0;JMP")
		// o.WriteCode("(" + label1 + ")")
		// o.WriteCode("@SP")
		// o.WriteCode("M=M")
		// o.WriteCode("A=M")
		// o.WriteCode("M=-1")
		// o.WriteCode("(" + label2 + ")")
		// o.WriteCode("@SP")
		// o.WriteCode("M=M+1")

		// ランダムなユニークラベルを生成
		label_y_neg := "Y_NEG_" + RandomString()
		label_y_pos := "Y_POS_" + RandomString()
		label_x_neg_y_neg := "X_NEG_Y_NEG_" + RandomString()
		label_x_pos_y_pos := "X_POS_Y_POS_" + RandomString()
		label_push_true := "PUSH_TRUE_" + RandomString()
		label_push_false := "PUSH_FALSE_" + RandomString()
		label_end := "END_" + RandomString()

		// より安全な比較演算(gt)の実装
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M") // y値を取得
		o.WriteCode("@" + label_y_neg)
		o.WriteCode("D;JLT") // yが負数の場合分岐
		o.WriteCode("@" + label_y_pos)
		o.WriteCode("0;JMP") // yが正数の場合分岐

		o.WriteCode("(" + label_y_neg + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M") // x値を取得
		o.WriteCode("@" + label_x_neg_y_neg)
		o.WriteCode("D;JLT") // xも負数の場合
		// xは正、yは負 → x > y は常に真
		o.WriteCode("@SP")
		o.WriteCode("A=M")
		o.WriteCode("M=-1") // trueをプッシュ
		o.WriteCode("@" + label_end)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_x_neg_y_neg + ")")
		// 両方負の場合、通常の比較で問題なし（オーバーフローしない）
		o.WriteCode("@SP")
		o.WriteCode("M=M+1") // スタックを戻す
		o.WriteCode("A=M")
		o.WriteCode("D=M") // yを再度取得
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")   // xにポイント
		o.WriteCode("D=M-D") // x-y
		o.WriteCode("@" + label_push_true)
		o.WriteCode("D;JGT") // ここがltとは逆（JGT）
		o.WriteCode("@" + label_push_false)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_y_pos + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")
		o.WriteCode("D=M") // x値を取得
		o.WriteCode("@" + label_x_pos_y_pos)
		o.WriteCode("D;JGE") // xが正または0の場合
		// xが負、yが正 → x < y は常に真なので、x > y は常に偽
		o.WriteCode("@" + label_push_false)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_x_pos_y_pos + ")")
		// 両方正の場合、通常の比較で問題なし
		o.WriteCode("@SP")
		o.WriteCode("M=M+1") // スタックを戻す
		o.WriteCode("A=M")
		o.WriteCode("D=M") // yを再度取得
		o.WriteCode("@SP")
		o.WriteCode("M=M-1")
		o.WriteCode("A=M")   // xにポイント
		o.WriteCode("D=M-D") // x-y
		o.WriteCode("@" + label_push_true)
		o.WriteCode("D;JGT") // ここがltとは逆（JGT）
		o.WriteCode("@" + label_push_false)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_push_true + ")")
		o.WriteCode("@SP")
		o.WriteCode("A=M")
		o.WriteCode("M=-1") // trueをプッシュ (-1)
		o.WriteCode("@" + label_end)
		o.WriteCode("0;JMP")

		o.WriteCode("(" + label_push_false + ")")
		o.WriteCode("@SP")
		o.WriteCode("A=M")
		o.WriteCode("M=0") // falseをプッシュ (0)

		o.WriteCode("(" + label_end + ")")
		o.WriteCode("@SP")
		o.WriteCode("M=M+1") // スタックポインタを増加
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

func (o *Output) WritePushPop(cmdType parser.CmdType, command string, arg2 string, fileName string) {
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
			o.WriteCode("@" + fileName + "." + arg2)
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
			// o.WriteCode("@R13")
			// o.WriteCode("M=0")
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
			// o.WriteCode("@R13")
			// o.WriteCode("M=0")
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
			// o.WriteCode("@R13")
			// o.WriteCode("M=0")
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
			// o.WriteCode("@R13")
			// o.WriteCode("M=0")
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
			// o.WriteCode("@R13")
			// o.WriteCode("M=0")
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
			o.WriteCode("@" + fileName + "." + arg2)
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
			// o.WriteCode("@R13")
			// o.WriteCode("M=0")
		}
	}
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

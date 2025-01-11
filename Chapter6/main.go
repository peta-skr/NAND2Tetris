/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/peta-skr/NAND2Tetris/code"
	"github.com/peta-skr/NAND2Tetris/parser"
)

func main() {

	/*
	コマンドからファイル名を渡す処理は後回し
	*/
	//cmd.Execute()


	// parserの初期化
	
	// 1行ずつバイナリコードに変換する
}

func Assemble(filepath string) string {
	// parserの初期化
	var data parser.InputData
	var err error
	data, err = parser.Initialize(filepath)
	var parsedData []string

	if err != nil {
		fmt.Println(err)
		return ""
	}

	for data.HasMoreCommands() {
		data.Advance()


		switch data.CommandType() {
		case parser.A_COMMAND:
			str, found := strings.CutPrefix(data.Text[data.Index], "@")

			if !found {
				fmt.Errorf("some error")
			}

			num, err := strconv.Atoi(str)

			if err != nil {
				fmt.Errorf("some error")
			}

			binary_str := strconv.FormatInt(int64(num), 2)

			//16桁になるまで0を追加する
			binary_str = fmt.Sprintf("%016s", binary_str)

			parsedData = append(parsedData, binary_str)

		case parser.C_COMMAND:
			dest := code.Dest(data.Dest())
			comp := code.Comp(data.Comp())
			jump := code.Jump(data.Jump())

			binary_str := "111" + comp + dest + jump

			parsedData = append(parsedData, binary_str)


		case parser.L_COMMAND:
		}
	}
	
	output := ""

	for i := range len(parsedData) {
		output += parsedData[i]
		if len(parsedData) - 1 > i {
			output += "\n"
		}
	}

	return output
}

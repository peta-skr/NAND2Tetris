package assemble

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/peta-skr/NAND2Tetris/code"
	"github.com/peta-skr/NAND2Tetris/parser"
	"github.com/peta-skr/NAND2Tetris/symbolTable"
)

func Assemble(filepath string) string {
	// parserの初期化
	data, err := parser.Initialize(filepath)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var parsedData []string
	var addressCounter int
	symbolTable := symbolTable.Initialize()


	/*
	* 1回目のループ
	* シンボルテーブルにラベルシンボルの情報を追加していく
	*/
	for data.HasMoreCommands() {
		data.Advance()

		switch data.CommandType() {
		case parser.A_COMMAND:
			addressCounter++

		case parser.C_COMMAND:
			addressCounter++

		case parser.L_COMMAND:
			symbol := data.Symbol()
			symbolTable.AddEntry(symbol, addressCounter)

		}
	}

	// もう一度ループするため
	data.Index = -1

	// 変数シンボルは16番目のアドレスから割り当てるため
	addressCounter = 16

	/*
	* 2回目のループは実際に変換
	*/
	for data.HasMoreCommands() {
		data.Advance()


		switch data.CommandType() {
		case parser.A_COMMAND:
			str, found := strings.CutPrefix(data.Text[data.Index], "@")
			if !found {
				fmt.Println("fail to cutPrefix")
			}

			num, err := strconv.Atoi(str)
			var binary_str string

			if err != nil {
				if symbolTable.Contains(str) {
				 	address := symbolTable.GetAddress(str)
					binary_str = strconv.FormatInt(int64(address), 2)
				}else {
					symbolTable.AddEntry(str, addressCounter)
					binary_str = strconv.FormatInt(int64(addressCounter), 2)
					addressCounter++
				}
			}else {
				binary_str = strconv.FormatInt(int64(num), 2)
			}

			//16桁になるまで0を追加する
			binary_str = fmt.Sprintf("%016s", binary_str)
			parsedData = append(parsedData, binary_str)

		case parser.C_COMMAND:
			dest := code.Dest(data.Dest())
			comp := code.Comp(data.Comp())
			jump := code.Jump(data.Jump())

			binary_str := "111" + comp + dest + jump
			parsedData = append(parsedData, binary_str)
		}
	}
	
	output := ""
	parsedLen := len(parsedData)

	for i := range parsedLen {
		output += parsedData[i]
		if parsedLen - 1 > i {
			output += "\n"
		}
	}

	return output
}

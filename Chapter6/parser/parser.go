package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Parser() {

}

type CommandType int

const (
	A_COMMAND CommandType = iota
	C_COMMAND
	L_COMMAND
)

type InputData struct {
	text []string
	index int
	length int
	commandType CommandType
}


/*
ファイルを開きパースを行う準備をする
*/
func Initialize(filePath string) (inputData InputData, err error) {
	
	f, err := os.Open(filePath);
	
	if err != nil {
		return inputData, err
	}
	
	defer f.Close()
	
	inputData = InputData{}
	inputData.index = -1
	
	scanner := bufio.NewScanner(f)
	length := 0
	for scanner.Scan() {
		length++
		inputData.text = append(inputData.text, scanner.Text())
	}
	
	inputData.length = length

	return inputData, nil
}

/*
次のコマンドがあるかチェック
*/
func (i *InputData) HasMoreCommands() bool {
	return i.length > i.index + 1
}

/*
次のコマンドを読み込む
*/
func (i *InputData) Advance(){

	if !i.HasMoreCommands() {
		return
	}

	i.index++

	// ここで次のコマンドを変換する
	cmd := i.text[i.index]

	if strings.HasPrefix(cmd, "@") {
		i.commandType = A_COMMAND
	} else if strings.HasPrefix(cmd, "(") {
		i.commandType = L_COMMAND
	}else {
	}

}

func (i *InputData) CommandType() CommandType {
	return i.commandType
}

func (i *InputData) Symbol() string {

	cmd := i.text[i.index]

	switch i.CommandType() {
	case A_COMMAND:
		symbol, found := strings.CutPrefix(cmd, "@")
		if !found {
			fmt.Println("fail to  cut `@` from A_COMMAND")
			return ""
		}

		return symbol
	case L_COMMAND:
		s, f := strings.CutPrefix(cmd, "(")

		if !f {
			fmt.Println("fail to  cut `(` from L_COMMAND")
			return ""
		}

		symbol, found := strings.CutSuffix(s, "(")
		if !found {
			fmt.Println("fail to  cut `)` from L_COMMAND")
			return ""
		}

		return symbol
	default:
		fmt.Println("current CommandType is neither A_COMMAND nor L_COMMAND")
		return ""
	}
}

func (i *InputData) dest() string {
	if i.commandType != C_COMMAND {
		return "current CommandType is not C_COMMAND"
	}

	cmd := i.text[i.index]

	dest, _, found := strings.Cut(cmd, "=")

	if !found {
		return "null"
	}

	return dest
}


func (i *InputData) comp() string {
	if i.commandType != C_COMMAND {
		return "current CommandType is not C_COMMAND"
	}

	cmd := i.text[i.index]

	_, after, found := strings.Cut(cmd, "=")

	if !found {
		return "some error"
	}

	var comp string

	if strings.HasSuffix(after, ";") {
		c, found := strings.CutPrefix(after, ";")

		if !found {
			return "some error"
		}

		comp = c
	}else {
		comp = after
	}

	return comp

}

func (i *InputData) jump() string {
	if i.commandType != C_COMMAND {
		return "current CommandType is not C_COMMAND"
	}

	cmd := i.text[i.index]

	jump, found := strings.CutSuffix(cmd, ";")

	if !found {
		return "some error"
	}

	return jump
}
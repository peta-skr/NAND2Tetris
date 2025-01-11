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
	IGNORE_COMMAND CommandType = iota
	A_COMMAND
	C_COMMAND
	L_COMMAND
)

type InputData struct {
	Text []string
	Index int
	Length int
	CmdType CommandType
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
	inputData.Index = -1
	
	scanner := bufio.NewScanner(f)
	length := 0
	for scanner.Scan() {
		length++
		inputData.Text = append(inputData.Text, scanner.Text())
	}
	
	inputData.Length = length

	return inputData, nil
}

/*
次のコマンドがあるかチェック
*/
func (i *InputData) HasMoreCommands() bool {
	return i.Length > i.Index + 1
}

/*
次のコマンドを読み込む
*/
func (i *InputData) Advance(){

	if !i.HasMoreCommands() {
		return
	}

	i.Index++

	cmd := i.Text[i.Index]

	if strings.HasPrefix(cmd, "@") {
		i.CmdType = A_COMMAND
	} else if strings.HasPrefix(cmd, "(") {
		i.CmdType = L_COMMAND
	}else if !strings.HasPrefix(cmd, "//") && cmd != "" {
		i.CmdType = C_COMMAND
	}

}

func (i *InputData) CommandType() CommandType {
	return i.CmdType
}

func (i *InputData) Symbol() string {

	cmd := i.Text[i.Index]

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

func (i *InputData) Dest() string {
	if i.CommandType() != C_COMMAND {
		return "current CommandType is not C_COMMAND"
	}

	cmd := i.Text[i.Index]

	if strings.Contains(cmd, "=") {
		dest, _, _ := strings.Cut(cmd, "=")

		return strings.Trim(dest, " ")
	}

	return "null"

}


func (i *InputData) Comp() string {
	if i.CommandType() != C_COMMAND {
		return "current CommandType is not C_COMMAND"
	}

	cmd := i.Text[i.Index]
	var cutDest string

	if strings.Contains(cmd, "=") {

		_, c, found := strings.Cut(cmd, "=")

		if !found {
			return "some error"
		}
		
		cutDest = c
	}else {
		cutDest = cmd
	}

	var comp string

	if strings.Contains(cutDest, ";") {
		c, _, found := strings.Cut(cutDest, ";")

		if !found {
			return "some error"
		}

		comp = c
	}else {
		comp = cutDest
	}

	return strings.Trim(comp, " ")

}

func (i *InputData) Jump() string {
	if i.CmdType != C_COMMAND {
		return "current CommandType is not C_COMMAND"
	}

	cmd := i.Text[i.Index]

	_, jump, found := strings.Cut(cmd, ";")

	if !found {
		return "null"
	}

	return strings.Trim(jump, " ")
}
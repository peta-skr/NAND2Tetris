package parser

import (
	"bufio"
	"os"
	"strings"
)

type CmdType int

const (
	C_ARITHMETIC CmdType = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

type VMCode struct {
	data []string
	length int
	index int
	cmdType CmdType
}

func Constructor(filePath string) (VMCode, error) {
	inputData := VMCode{}

	f, err := os.Open(filePath)

	if err != nil {
		return inputData, err
	}

	defer f.Close()

	inputData.index = -1

	scanner := bufio.NewScanner(f)
	length := 0
	for scanner.Scan() {
		length++
		inputData.data = append(inputData.data, strings.Trim(scanner.Text(), " "))
	}

	inputData.
}

func HasMoreCommands() bool {
	return true
}

func CommandType() Command {
	return C_ARITHMETIC
}

func Arg1() string {
	return ""
}

func Arg2() int {
	return 0
}
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
	data    []string
	length  int
	index   int
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
	for scanner.Scan() {
		// コメントの削除と、前後の空白の削除
		text := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(text, "//") {
			inputData.data = append(inputData.data, text)
		}
	}

	inputData.length = len(inputData.data)

	return inputData, nil
}

func (v *VMCode) HasMoreCommands() bool {
	return v.length > v.index+1
}

func (v *VMCode) Advance() {
	if !v.HasMoreCommands() {
		return
	}

	v.index++
	line := v.data[v.index]

	switch {
	case strings.HasPrefix(line, "push"):
		v.cmdType = C_PUSH
	case strings.HasPrefix(line, "pop"):
		v.cmdType = C_POP
	case strings.HasPrefix(line, "label"):
		v.cmdType = C_LABEL
	case strings.HasPrefix(line, "if-goto"):
		v.cmdType = C_IF
	case strings.HasPrefix(line, "goto"):
		v.cmdType = C_GOTO
	case strings.HasPrefix(line, "call"):
		v.cmdType = C_CALL
	case strings.HasPrefix(line, "function"):
		v.cmdType = C_FUNCTION
	case strings.HasPrefix(line, "return"):
		v.cmdType = C_RETURN
	default:
		v.cmdType = C_ARITHMETIC
	}
}

func (v *VMCode) CommandType() CmdType {
	return v.cmdType
}

func (v *VMCode) Arg1() string {
	parts := strings.Split(v.data[v.index], " ")
	if v.cmdType == C_ARITHMETIC {
		return parts[0]
	}
	return parts[1]
}

func (v *VMCode) Arg2() string {
	if v.cmdType == C_PUSH || v.cmdType == C_POP || v.cmdType == C_CALL || v.cmdType == C_FUNCTION {
		parts := strings.Split(v.data[v.index], " ")
		return parts[2]
	}
	return ""
}

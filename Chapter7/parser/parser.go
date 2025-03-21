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
		text := strings.Trim(scanner.Text(), " ")
		if !strings.HasPrefix(text, "//") {
			length++
			inputData.data = append(inputData.data, text)
		}
	}

	inputData.length = length

	return inputData, nil
}

func (v *VMCode) HasMoreCommands() bool {
	return v.length > v.index + 1
}

func (v *VMCode) Advance() {
	if !v.HasMoreCommands() {
		return
	}

	v.index++

	if strings.HasPrefix(v.data[v.index], "push") {
		v.cmdType = C_PUSH
	}else if strings.HasPrefix(v.data[v.index], "pop") {
		v.cmdType = C_POP
	}else {
		v.cmdType = C_ARITHMETIC
	}
}

func (v *VMCode) CommandType() CmdType {
	return v.cmdType
}

func (v *VMCode) Arg1() string {
	if v.cmdType == C_ARITHMETIC {
		l := strings.Split(v.data[v.index], " ")
		return l[0]
	}else if v.cmdType == C_PUSH {
		l := strings.Split(v.data[v.index], " ")
		return l[1]
	}else if v.cmdType == C_POP {
		l := strings.Split(v.data[v.index], " ")
		return l[1]
	}
	return ""
}

func (v *VMCode) Arg2() string {	
	if v.cmdType == C_PUSH {
		l := strings.Split(v.data[v.index], " ")
		return l[2]
	}else if v.cmdType == C_POP {
		l := strings.Split(v.data[v.index], " ")
		return l[2]
	}
	return ""
}
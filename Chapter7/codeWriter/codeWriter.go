package codewriter

import (
	"fmt"
	"os"
)

type Output struct {
	file *os.File
	filename string
}

func Constructor(filepath string) Output {
	file, err := os.OpenFile(filepath, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644)

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
		_, _ = o.file.WriteString("")
	case "sub":
	case "neg":
	}

}

func writePushPop() {
	
}

func (o *Output) close() {
	o.file.Close()
	
}
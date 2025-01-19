package codewriter

import (
	"fmt"
	"os"
)

type Output struct {
	file *os.File
	filename string
}

func Constructor() Output {
	file, err := os.OpenFile("./example.txt", os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644)

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

func WriteArithmetic(command string) {
	
}

func writePushPop() {
	
}

func close() {
	defer file.Close()
	
}
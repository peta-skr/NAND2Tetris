package jackanalyzer

import (
	codegenerator "Chapter11/CodeGenerator"
	compilationengine "Chapter11/CompilationEngine"
	jacktokenizer "Chapter11/JackTokenizer"
	"fmt"
	"os"
	"path/filepath"
)

func Analyzer(source string) {

	info, err := os.Stat(source)
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}

	if info.IsDir() {
		// ディレクトリの場合
		files, err := os.ReadDir(source)
		if err != nil {
			fmt.Println("エラー:", err)
			return
		}

		first := true

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".jack" {
				fmt.Println("Dir: " + file.Name())
				Analyze(filepath.Join(source, file.Name()), first)
				first = false
			}
		}
	} else if filepath.Ext(source) == ".jack" {
		// ファイルの場合
		fmt.Println(filepath.Base(source))
		Analyze(source, true)

	} else {
		fmt.Println("ファイルまたはディレクトリを指定してください")
	}
}

func Analyze(source string, fristFlag bool) {
	tokenizer := jacktokenizer.Tokenizer(source)
	vmFilePath := source[:len(source)-5] + ".vm"
	// vmFilePath := filepath.Dir(source) + "/Main.vm"

	if fristFlag {

		// 同名のファイルがある場合は削除
		if _, err := os.Stat(vmFilePath); err == nil {
			err := os.Remove(vmFilePath)
			if err != nil {
				fmt.Println("エラー:", err)
				return
			}
		}
	}

	parseTree, symboltable := compilationengine.Compile(*tokenizer, vmFilePath)

	codegenerator.GenerateCode(parseTree, symboltable, vmFilePath)

}

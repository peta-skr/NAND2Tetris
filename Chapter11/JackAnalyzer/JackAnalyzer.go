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

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".jack" {
				fmt.Println(file.Name())
				Analyze(filepath.Join(source, file.Name()))
			}
		}
	} else if filepath.Ext(source) == ".jack" {
		// ファイルの場合
		fmt.Println(filepath.Base(source))
		Analyze(source)

	} else {
		fmt.Println("ファイルまたはディレクトリを指定してください")
	}
}

func Analyze(source string) {
	tokenizer := jacktokenizer.Tokenizer(source)
	vmFilePath := source[:len(source)-5] + ".vm"
	parseTree, symboltable := compilationengine.Compile(*tokenizer, vmFilePath)
	fmt.Println("parseTree:", parseTree)
	fmt.Println("symboltable:", symboltable)

	codegenerator.GenerateCode(parseTree, symboltable, vmFilePath)

}

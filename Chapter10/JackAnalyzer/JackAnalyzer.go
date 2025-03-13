package jackanalyzer

import (
	jacktokenizer "Chapter10/JackTokenizer"
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
				jacktokenizer.Tokenizer(filepath.Join(source, file.Name()))
			}
		}
	} else if filepath.Ext(source) == ".jack" {
		// ファイルの場合
		fmt.Println(filepath.Base(source))
		jacktokenizer.Tokenizer(source)
	} else {
		fmt.Println("ファイルまたはディレクトリを指定してください")
	}
}

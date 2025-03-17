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

	xmlFilePath := source[:len(source)-5] + "T_test.xml" // test用
	xmlFile, err := os.Create(xmlFilePath)
	if err != nil {
		fmt.Println("XMLファイルを作成できませんでした: ", err)
		return
	}
	defer xmlFile.Close()

	// ファイルにトークンを書き込んでいく
	xmlFile.WriteString("<tokens>\r\n")

	for tokenizer.HasMoreTokens() {
		token := tokenizer.GetTokenValue()
		tokenType := tokenizer.GetTokenType()
		switch tokenType {
		case jacktokenizer.KEYWORD:
			xmlFile.WriteString("<keyword> " + token + " </keyword>\r\n")
		case jacktokenizer.SYMBOL:
			xmlFile.WriteString("<symbol> " + token + " </symbol>\r\n")
		case jacktokenizer.IDENTIFIER:
			xmlFile.WriteString("<identifier> " + token + " </identifier>\r\n")
		case jacktokenizer.INT_CONST:
			xmlFile.WriteString("<integerConstant> " + token + " </integerConstant>\r\n")
		case jacktokenizer.STRING_CONST:
			xmlFile.WriteString("<stringConstant> " + token + " </stringConstant>\r\n")
		}
		tokenizer.Advance()
	}

	xmlFile.WriteString("</tokens>\r\n")

}

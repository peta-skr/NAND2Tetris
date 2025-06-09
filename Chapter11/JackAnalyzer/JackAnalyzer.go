package jackanalyzer

import (
	codegenerator "Chapter11/CodeGenerator"
	compilationengine "Chapter11/CompilationEngine"
	jacktokenizer "Chapter11/JackTokenizer"
	symboltable "Chapter11/SymbolTable"
	"fmt"
	"os"
	"path/filepath"
)

var subroutineKindMap = make(map[string]map[string]string) // クラス名 -> サブルーチン名 -> サブルーチンの種類("function"/"method"/"constructor")
var compileMap = make(map[string]any)                      // クラス名 -> オブジェクト

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
		// mainPath := "" // Main.jackのパスを保持するための変数

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".jack" {
				fmt.Println("Dir: " + file.Name())
				// if file.Name() == "Main.jack" {
				// Main.jackは最後に処理する
				// 	mainPath = filepath.Join(source, file.Name())
				// 	continue
				// }
				Analyze(filepath.Join(source, file.Name()), first)
				first = false
			}
		}

		for _, file := range files {
			if filepath.Ext(file.Name()) == ".jack" {

				// fileName := filepath.Base(source)
				// nameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]

				data := compileMap[file.Name()].([]any)
				parseTree := data[0].(compilationengine.ParseTree)
				symTable := data[1].(symboltable.SymbolTable)

				filePath := filepath.Join(source, file.Name())
				vmFilePath := filePath[:len(filePath)-5] + ".vm"

				generateVMCode(
					parseTree,
					symTable,
					vmFilePath)
			}
		}
		// Main.jackを最後に処理
		// if mainPath != "" {
		// 	Analyze(mainPath, first)
		// }
	} else if filepath.Ext(source) == ".jack" {
		// ファイルの場合
		Analyze(source, true)

	} else {
		fmt.Println("ファイルまたはディレクトリを指定してください")
	}
}

func Analyze(source string, fristFlag bool) {
	tokenizer := jacktokenizer.Tokenizer(source)
	vmFilePath := source[:len(source)-5] + ".vm"
	// vmFilePath := filepath.Dir(source) + "/Main.vm"

	// if fristFlag {

	// 	// 同名のファイルがある場合は削除
	// 	if _, err := os.Stat(vmFilePath); err == nil {
	// 		err := os.Remove(vmFilePath)
	// 		if err != nil {
	// 			fmt.Println("エラー:", err)
	// 			return
	// 		}
	// 	}
	// }

	parseTree, symboltable, subroutineTable := compilationengine.Compile(*tokenizer, vmFilePath)
	fileName := filepath.Base(source)
	nameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	compileMap[filepath.Base(source)] = []any{parseTree, symboltable} // クラス名 -> [parseTree, symboltable]

	subroutineKindMap[nameWithoutExt] = subroutineTable

	// subroutineMap := codegenerator.GenerateCode(parseTree, symboltable, vmFilePath, subroutineKindMap)
	// subroutineKindMap = subroutineMap // クラス名ごとにサブルーチンの種類を記録

}

func generateVMCode(parseTree compilationengine.ParseTree, symboltable symboltable.SymbolTable, vmFilePath string) {
	codegenerator.GenerateCode(parseTree, symboltable, vmFilePath, subroutineKindMap)
}

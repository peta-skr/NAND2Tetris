package main

import (
	jackanalyzer "Chapter10/JackAnalyzer"
	"fmt"
	"os"
)

func main() {
	// コマンドライン引数を取得
	args := os.Args
	fmt.Println(args) // 引数を表示
	// 引数がない場合は、エラーを表示して終了
	if len(args) < 2 {
		fmt.Println("引数が足りません")
		return
	}

	path := args[1]

	jackanalyzer.Analyzer(path)

}

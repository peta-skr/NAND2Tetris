package main

import (
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

	folerPath := args[1]

	files, err := os.ReadDir(folerPath)
	if err != nil {
		fmt.Println("エラー:", err)
		return
	}

	for _, file := range files {
		fmt.Println(file.Name())
	}

}

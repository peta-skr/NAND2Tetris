package jackanalyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTokenizer(t *testing.T) {
	// テストディレクトリのパス
	testDir := "../test"

	// テストディレクトリが存在しない場合はエラー
	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != testDir {
			// Analyze関数を呼び出す
			Analyzer(path)
		} else if !info.IsDir() && filepath.Ext(info.Name()) == ".jack" { // ファイルがディレクトリでなく、拡張子が .jack の場合にテストを実行
			t.Run(info.Name(), func(t *testing.T) {
				// 生成されたVMファイルのパス
				generatedVMPath := path[:len(path)-5] + ".vm"
				// 期待されるVMファイルのパス
				expectedVMPath := path[:len(path)-5] + ".vm.ans"

				// 生成されたVMファイルを読み込む
				generatedVM, err := os.ReadFile(generatedVMPath)
				if err != nil {
					t.Fatalf("生成されたVMファイルを読み取れませんでした: %v", err)
				}
				// 期待されるVMファイルを読み込む
				expectedVM, err := os.ReadFile(expectedVMPath)
				if err != nil {
					t.Fatalf("期待されるVMファイルを読み取れませんでした: %v", err)
				}

				// 比較部分
				generatedNormalized := normalizeContent(string(generatedVM))
				expectedNormalized := normalizeContent(string(expectedVM))

				if generatedNormalized != expectedNormalized {
					// 文字ごとに比較してどこが違うか特定
					for i := 0; i < len(generatedNormalized) && i < len(expectedNormalized); i++ {
						if generatedNormalized[i] != expectedNormalized[i] {
							fmt.Printf("違いが見つかりました位置 %d: generated='%c'(%d), expected='%c'(%d)\n",
								i, generatedNormalized[i], generatedNormalized[i],
								expectedNormalized[i], expectedNormalized[i])
							break
						}
					}
					t.Errorf("VMファイルの内容が一致しません: %s", path)
				}
			})

		}
		return nil
	})
	if err != nil {
		t.Fatalf("テストディレクトリを読み取れませんでした: %v", err)
	}
}

func normalizeContent(s string) string {
	// 改行コード統一
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	// 行末の空白削除
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	// 最後の空行削除
	s = strings.TrimRight(strings.Join(lines, "\n"), "\n")
	return s
}

package jackanalyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTokenizer(t *testing.T) {
	// テストディレクトリのパス
	testDir := "../test"

	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == ".jack" {
			t.Run(info.Name(), func(t *testing.T) {
				Analyze(path)

				// 生成されたXMLファイルのパス
				generatedXMLPath := path[:len(path)-5] + "T_test.xml"
				// 期待されるXMLファイルのパス
				expectedXMLPath := path[:len(path)-5] + "T.xml"

				// 生成されたXMLファイルを読み込む
				generatedXML, err := os.ReadFile(generatedXMLPath)
				if err != nil {
					t.Fatalf("生成されたXMLファイルを読み取れませんでした: %v", err)
				}

				// 期待されるXMLファイルを読み込む
				expectedXML, err := os.ReadFile(expectedXMLPath)
				if err != nil {
					t.Fatalf("期待されるXMLファイルを読み取れませんでした: %v", err)
				}

				// fmt.Println()

				// 生成されたXMLファイルと期待されるXMLファイルの内容を比較
				if string(generatedXML) != string(expectedXML) {
					t.Errorf("XMLファイルの内容が一致しません: %s", path)
				}
			})
		}
		return nil
	})
	if err != nil {
		t.Fatalf("テストディレクトリを読み取れませんでした: %v", err)
	}
}

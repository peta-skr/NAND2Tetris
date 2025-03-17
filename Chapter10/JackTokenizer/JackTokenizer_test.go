package jacktokenizer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTokenizer(t *testing.T) {
	// テストディレクトリのパス
	testDir := "./test"

	// テストディレクトリ内のファイルを取得
	files, err := os.ReadDir(testDir)
	if err != nil {
		t.Fatalf("テストディレクトリを読み取れませんでした: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".jack" {
			t.Run(file.Name(), func(t *testing.T) {
				// トークナイザを初期化
				tokenizer := Tokenizer(filepath.Join(testDir, file.Name()))
				// tokenizer.Tokenize()

				// トークンが正しく生成されているかを確認
				if !tokenizer.HasMoreTokens() {
					t.Errorf("トークンが生成されていません: %s", file.Name())
				}

				// トークンを進めて確認
				for tokenizer.HasMoreTokens() {
					token := tokenizer.TokenValue()
					tokenType := tokenizer.TokenType()
					t.Logf("トークン: %s, タイプ: %v", token, tokenType)
					tokenizer.Advance()
				}
			})
		}
	}
}

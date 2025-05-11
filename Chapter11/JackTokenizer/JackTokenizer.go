package jacktokenizer

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type TokenType int

const (
	KEYWORD TokenType = iota
	SYMBOL
	IDENTIFIER
	INT_CONST
	STRING_CONST
)

type Keyword string

const (
	CLASS       = "CLASS"
	METHOD      = "METHOD"
	FUNCTION    = "FUNCTION"
	CONSTRUCTOR = "CONSTRUCTOR"
	INT         = "INT"
	BOOLEAN     = "BOOLEAN"
	CHAR        = "CHAR"
	VOID        = "VOID"
	VAR         = "VAR"
	STATIC      = "STATIC"
	FIELD       = "FIELD"
	LET         = "LET"
	DO          = "DO"
	IF          = "IF"
	ELSE        = "ELSE"
	WHILE       = "WHILE"
	RETURN      = "RETURN"
	TRUE        = "TRUE"
	FALSE       = "FALSE"
	NULL        = "NULL"
	THIS        = "THIS"
)

var keywords = map[string]Keyword{
	"class":       CLASS,
	"method":      METHOD,
	"function":    FUNCTION,
	"constructor": CONSTRUCTOR,
	"int":         INT,
	"boolean":     BOOLEAN,
	"char":        CHAR,
	"void":        VOID,
	"var":         VAR,
	"static":      STATIC,
	"field":       FIELD,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
}

type Token struct {
	Type  TokenType
	Value string
}

type JackTokenizer struct {
	source []string
	tokens []Token
	index  int
	last   int
}

func Tokenizer(source string) *JackTokenizer {
	// ファイルを開く
	file, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content []string
	tokens := make([]Token, 0)
	inBlockComment := false

	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\t")
		// コメントや空行をスキップ
		if isComment(line, &inBlockComment) || strings.TrimSpace(line) == "" {
			continue
		}

		// トークンを取得
		arr := strings.Split(line, "")

		for index := 0; index < len(arr); index++ {
			str := arr[index]
			if str == "" || str == " " {
				continue
			}

			// トークンの種類を判定
			switch {
			case isKeyword(str):
				tokens = append(tokens, Token{Type: KEYWORD, Value: str})
			case isSymbol(str):
				switch str {
				case "<":
					tokens = append(tokens, Token{Type: SYMBOL, Value: "&lt;"})
				case ">":
					tokens = append(tokens, Token{Type: SYMBOL, Value: "&gt;"})
				case "&":
					tokens = append(tokens, Token{Type: SYMBOL, Value: "&amp;"})
				case "/":
					index++
					if arr[index] == "/" || arr[index] == "*" {
						index = len(arr) - 1
						break
					} else {
						index--
						tokens = append(tokens, Token{Type: SYMBOL, Value: str})
					}
				default:
					tokens = append(tokens, Token{Type: SYMBOL, Value: str})
				}
			case isIntConst(str):
				indeToken := str
				for {
					index++
					str = arr[index]
					if index < len(arr) && isIntConst(str) {
						indeToken += str
					} else {
						index--
						break
					}
				}
				tokens = append(tokens, Token{Type: INT_CONST, Value: indeToken})
			case isStringConst(str):
				stringToken := ""
				for {
					index++
					if index < len(arr) {
						str = arr[index]
						if strings.HasSuffix(str, "\"") {
							break
						}
						stringToken += str
					} else {
						break
					}
				}
				tokens = append(tokens, Token{Type: STRING_CONST, Value: stringToken})
			default:
				// value := str
				value := ""

				if !isSymbol(str) && str != " " && str != "\t" && str != "\n" && str != "\r" {
					value += str
				}

				for {
					index++
					if index < len(arr) {
						str = arr[index]
						if isSymbol(str) || str == " " || str == "\t" || str == "\n" || str == "\r" {
							// if isSymbol(str) || unicode.IsSpace(rune(str[0])) {
							index-- // シンボルやスペースに出会ったらインデックスを戻す
							break
						}
						value += str
					} else {
						break
					}
				}
				if isKeyword(value) {
					tokens = append(tokens, Token{Type: KEYWORD, Value: value})
				} else {
					if len(value) != 0 {
						tokens = append(tokens, Token{Type: IDENTIFIER, Value: value})
					}
				}
			}
		}
	}

	return &JackTokenizer{
		source: content,
		tokens: tokens,
		index:  0,
		last:   len(tokens),
	}
}

func isComment(line string, inBlockComment *bool) bool {
	trimmedLine := strings.TrimSpace(line)

	// ブロックコメント内にいる場合
	if *inBlockComment {
		if strings.Contains(trimmedLine, "*/") {
			*inBlockComment = false
		}
		return true
	}

	// 行コメント
	if strings.HasPrefix(trimmedLine, "//") {
		return true
	}

	// ブロックコメントの開始
	if strings.HasPrefix(trimmedLine, "/*") {
		if !strings.HasSuffix(trimmedLine, "*/") {
			*inBlockComment = true
		}
		return true
	}

	return false
}

func isKeyword(str string) bool {
	if _, ok := keywords[str]; ok {
		return true
	}

	return false
}

func isSymbol(str string) bool {
	return strings.Contains("{}()[].,;+-*/&|<>=~", str)
}

func isIntConst(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func isStringConst(str string) bool {
	if strings.HasPrefix(str, "\"") {
		return true
	}

	return false
}

func (j *JackTokenizer) HasMoreTokens() bool {
	return j.index < j.last
}

func (j *JackTokenizer) Advance() {
	if j.HasMoreTokens() {
		j.index++
	}
}

func (j *JackTokenizer) GetTokenType() TokenType {
	return j.tokens[j.index].Type
}

func (j *JackTokenizer) GetTokenValue() string {
	return j.tokens[j.index].Value
}

func (j *JackTokenizer) GetKeyword() Keyword {
	if j.GetTokenType() == KEYWORD {
		return keywords[j.tokens[j.index].Value]
	}
	return ""
}
func (j *JackTokenizer) Symbol() string {
	if j.GetTokenType() == SYMBOL {
		return j.tokens[j.index].Value
	}

	return ""
}

func (j *JackTokenizer) Identifier() string {
	if j.GetTokenType() == IDENTIFIER {
		return j.tokens[j.index].Value
	}

	return ""
}

func (j *JackTokenizer) IntVal() int {
	if j.GetTokenType() == INT_CONST {
		i, _ := strconv.Atoi(j.tokens[j.index].Value)
		return i
	}

	return 0
}

func (j *JackTokenizer) StringVal() string {
	if j.GetTokenType() == STRING_CONST {
		return j.tokens[j.index].Value
	}

	return ""
}

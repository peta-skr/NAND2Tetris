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

type Keyword int

const (
	CLASS Keyword = iota
	METHOD
	FUNCTION
	CONSTRUCTOR
	INT
	BOOLEAN
	CHAR
	VOID
	VAR
	STATIC
	FIELD
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	TRUE
	FALSE
	NULL
	THIS
)

type Token struct {
	Type  TokenType
	Value string
}

type JackTokenizer struct {
	source string
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
	var content strings.Builder
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}

	return &JackTokenizer{
		source: content.String(),
		tokens: make([]Token, 0),
		index:  0,
		last:   0,
	}
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

func (j *JackTokenizer) GetKeyword() Keyword {
	if j.GetTokenType() == KEYWORD {
		switch j.tokens[j.index].Value {
		case "class":
			return CLASS
		case "method":
			return METHOD
		case "function":
			return FUNCTION
		case "constructor":
			return CONSTRUCTOR
		case "int":
			return INT
		case "boolean":
			return BOOLEAN
		case "char":
			return CHAR
		case "void":
			return VOID
		case "var":
			return VAR
		case "static":
			return STATIC
		case "field":
			return FIELD
		case "let":
			return LET
		case "do":
			return DO
		case "if":
			return IF
		case "else":
			return ELSE
		case "while":
			return WHILE
		case "return":
			return RETURN
		case "true":
			return TRUE
		case "false":
			return FALSE
		case "null":
			return NULL
		case "this":
			return THIS
		}
	}

	return 0
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

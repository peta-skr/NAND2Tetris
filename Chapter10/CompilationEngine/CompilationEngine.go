package compilationengine

import (
	jacktokenizer "Chapter10/JackTokenizer"
	"fmt"
	"os"
	"strings"
)

// Node インターフェース
type Node interface {
	isNode()
}

// ParseTree 構造体
type ParseTree struct {
	Nodes []Node
}

type ParseNode struct {
	Type     jacktokenizer.TokenType
	Value    string
	Children []Node
}

func (n ParseNode) isNode() {}

type ContainerNode struct {
	Name     string
	Children []Node
}

func (n ContainerNode) isNode() {}

func (tree *ParseTree) AddNode(node Node) {
	tree.Nodes = append(tree.Nodes, node)
}

// エラーチェック関数
func expectToken(tokenize *jacktokenizer.JackTokenizer, expectedType jacktokenizer.TokenType, expectedValue string) bool {
	return tokenize.GetTokenType() == expectedType && tokenize.GetTokenValue() == expectedValue
}

// トークンをノードとして追加するヘルパー関数
func addTokenNode(tokenize *jacktokenizer.JackTokenizer, parentNode *ContainerNode) {
	node := ParseNode{Type: tokenize.GetTokenType(), Value: tokenize.GetTokenValue()}
	parentNode.Children = append(parentNode.Children, node)
	tokenize.Advance()
}

func Compile(tokenize jacktokenizer.JackTokenizer, outputFile string) {
	parseTree := ParseTree{}

	// トークンを解析してパースツリーを構築
	for tokenize.HasMoreTokens() {
		tokenType := tokenize.GetTokenType()

		switch tokenType {
		case jacktokenizer.KEYWORD:
			switch tokenize.GetKeyword() {
			case jacktokenizer.CLASS:
				parseTree.CompileClass(&tokenize)
			}
		}
		tokenize.Advance()
	}

	// パースツリーを文字列に変換
	xmlContent := buildXML(parseTree)

	// ファイルに書き込む
	err := os.WriteFile(outputFile, []byte(xmlContent), 0644)
	if err != nil {
		fmt.Println("XMLファイルを作成できませんでした: ", err)
	}
}

func buildXML(tree ParseTree) string {
	var builder strings.Builder

	for _, node := range tree.Nodes {
		writeNode(&builder, node, 0)
	}

	return builder.String()
}

func writeNode(builder *strings.Builder, node Node, indent int) {
	indentaion := strings.Repeat("  ", indent)

	switch n := node.(type) {
	case ParseNode:
		switch n.Type {
		case jacktokenizer.KEYWORD:
			builder.WriteString(fmt.Sprintf("%s<keyword> %s </keyword>\r\n", indentaion, n.Value))
		case jacktokenizer.SYMBOL:
			builder.WriteString(fmt.Sprintf("%s<symbol> %s </symbol>\r\n", indentaion, n.Value))
		case jacktokenizer.IDENTIFIER:
			builder.WriteString(fmt.Sprintf("%s<identifier> %s </identifier>\r\n", indentaion, n.Value))
		case jacktokenizer.INT_CONST:
			builder.WriteString(fmt.Sprintf("%s<integerConstant> %s </integerConstant>\r\n", indentaion, n.Value))
		case jacktokenizer.STRING_CONST:
			builder.WriteString(fmt.Sprintf("%s<stringConstant> %s </stringConstant>\r\n", indentaion, n.Value))
		}
	case ContainerNode:
		builder.WriteString(fmt.Sprintf("%s<%s>\r\n", indentaion, n.Name))
		for _, child := range n.Children {
			writeNode(builder, child, indent+1)
		}
		builder.WriteString(fmt.Sprintf("%s</%s>\r\n", indentaion, n.Name))

	}
}

func (p *ParseTree) CompileClass(tokenize *jacktokenizer.JackTokenizer) {
	classContentNode := ContainerNode{Name: "class", Children: []Node{}}

	// 'class' キーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "class") {
		return // エラー
	}
	addTokenNode(tokenize, &classContentNode)

	// クラス名
	if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	classNameNode := ParseNode{Type: jacktokenizer.IDENTIFIER, Value: tokenize.GetTokenValue()}
	classContentNode.Children = append(classContentNode.Children, classNameNode)
	tokenize.Advance()

	// '{' シンボル
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "{") {
		return // エラー
	}
	addTokenNode(tokenize, &classContentNode)

	for tokenize.HasMoreTokens() {
		tokenType := tokenize.GetTokenType()
		// '}' を検出したら終了
		if expectToken(tokenize, jacktokenizer.SYMBOL, "}") {
			break
		}

		switch tokenType {
		case jacktokenizer.KEYWORD:
			switch tokenize.GetKeyword() {
			case jacktokenizer.STATIC, jacktokenizer.FIELD:
				CompileClassVarDec(tokenize, &classContentNode)
			case jacktokenizer.METHOD, jacktokenizer.FUNCTION, jacktokenizer.CONSTRUCTOR:
				CompileSubroutine(tokenize, &classContentNode)
			default:
				return
			}
		default:
			return
		}
	}

	// // '}' シンボル
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "}") {
		return // エラー
	}
	addTokenNode(tokenize, &classContentNode)

	p.AddNode(classContentNode)
}

func CompileClassVarDec(tokenize *jacktokenizer.JackTokenizer, classContentNode *ContainerNode) {
	classVarDecNode := ContainerNode{Name: "classVarDec", Children: []Node{}}

	// 'static' または 'field' キーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "static") && !expectToken(tokenize, jacktokenizer.KEYWORD, "field") {
		return // エラー
	}
	addTokenNode(tokenize, &classVarDecNode)

	// 型
	if tokenize.GetTokenType() != jacktokenizer.KEYWORD && tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	addTokenNode(tokenize, &classVarDecNode)

	// 変数名
	if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	addTokenNode(tokenize, &classVarDecNode)

	// ',' または ';' を処理
	for tokenize.HasMoreTokens() {
		if expectToken(tokenize, jacktokenizer.SYMBOL, ";") {
			addTokenNode(tokenize, &classVarDecNode)
			break
		}

		if !expectToken(tokenize, jacktokenizer.SYMBOL, ",") {
			return // エラー
		}
		addTokenNode(tokenize, &classVarDecNode)

		if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
			return // エラー
		}
		addTokenNode(tokenize, &classVarDecNode)
	}

	classContentNode.Children = append(classContentNode.Children, classVarDecNode)
}

func CompileSubroutine(tokenize *jacktokenizer.JackTokenizer, classContentNode *ContainerNode) {
	subroutineDecNode := ContainerNode{Name: "subroutineDec", Children: []Node{}}

	// 'constructor', 'function', または 'method' キーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "constructor") &&
		!expectToken(tokenize, jacktokenizer.KEYWORD, "function") &&
		!expectToken(tokenize, jacktokenizer.KEYWORD, "method") {
		return // エラー
	}
	addTokenNode(tokenize, &subroutineDecNode)

	// 戻り値の型
	if tokenize.GetTokenType() != jacktokenizer.KEYWORD && tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	addTokenNode(tokenize, &subroutineDecNode)

	// サブルーチン名
	if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	addTokenNode(tokenize, &subroutineDecNode)

	// '('
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "(") {
		return // エラー
	}
	addTokenNode(tokenize, &subroutineDecNode)

	// パラメータリスト
	CompileParameterList(tokenize, &subroutineDecNode)

	// ')'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, ")") {
		return // エラー
	}
	addTokenNode(tokenize, &subroutineDecNode)

	// サブルーチン本体
	subroutineBodyNode := ContainerNode{Name: "subroutineBody", Children: []Node{}}

	// '{'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "{") {
		return // エラー
	}
	addTokenNode(tokenize, &subroutineBodyNode)

	// var 宣言
	for tokenize.GetTokenType() == jacktokenizer.KEYWORD && tokenize.GetKeyword() == jacktokenizer.VAR {
		CompileVarDec(tokenize, &subroutineBodyNode)
	}

	// ステートメント
	CompileStatements(tokenize, &subroutineBodyNode)

	// '}'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "}") {
		return // エラー
	}
	addTokenNode(tokenize, &subroutineBodyNode)

	// サブルーチン本体をサブルーチン宣言ノードに追加
	subroutineDecNode.Children = append(subroutineDecNode.Children, subroutineBodyNode)

	// クラス内容ノードに追加
	classContentNode.Children = append(classContentNode.Children, subroutineDecNode)

}

func CompileParameterList(tokenize *jacktokenizer.JackTokenizer, subroutineDecNode *ContainerNode) {
	parameterListNode := ContainerNode{Name: "parameterList", Children: []Node{}}

	// 引数ある分繰り返す
	for tokenize.HasMoreTokens() {
		if expectToken(tokenize, jacktokenizer.SYMBOL, ")") {
			break
		}

		// 引数の型
		if tokenize.GetTokenType() != jacktokenizer.KEYWORD &&
			(tokenize.GetKeyword() != jacktokenizer.INT &&
				tokenize.GetKeyword() != jacktokenizer.CHAR &&
				tokenize.GetKeyword() != jacktokenizer.BOOLEAN) &&
			(tokenize.GetTokenType() != jacktokenizer.IDENTIFIER) {
			return // エラー
		}
		addTokenNode(tokenize, &parameterListNode)

		// 引数名
		if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
			return // エラー
		}
		addTokenNode(tokenize, &parameterListNode)

		// ","
		if expectToken(tokenize, jacktokenizer.SYMBOL, ",") {
			addTokenNode(tokenize, &parameterListNode)
		}
	}

	subroutineDecNode.Children = append(subroutineDecNode.Children, parameterListNode)
}

func CompileVarDec(tokenize *jacktokenizer.JackTokenizer, subroutineBodyNode *ContainerNode) {
	varDecNode := ContainerNode{Name: "varDec", Children: []Node{}}

	// varキーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "var") {
		return // エラー
	}
	addTokenNode(tokenize, &varDecNode)

	// 変数の型
	if tokenize.GetTokenType() != jacktokenizer.KEYWORD &&
		(tokenize.GetKeyword() != jacktokenizer.INT &&
			tokenize.GetKeyword() != jacktokenizer.CHAR &&
			tokenize.GetKeyword() != jacktokenizer.BOOLEAN) &&
		(tokenize.GetTokenType() != jacktokenizer.IDENTIFIER) {
		return // エラー
	}
	addTokenNode(tokenize, &varDecNode)

	// 変数名
	if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	addTokenNode(tokenize, &varDecNode)

	// 変数ある分繰り返す
	for tokenize.HasMoreTokens() {
		if expectToken(tokenize, jacktokenizer.SYMBOL, ";") {
			break
		}

		// ","
		if !expectToken(tokenize, jacktokenizer.SYMBOL, ",") {
			return // エラー
		}
		addTokenNode(tokenize, &varDecNode)

		// 変数名
		if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
			return // エラー
		}
		addTokenNode(tokenize, &varDecNode)

	}
	// ";"
	if !expectToken(tokenize, jacktokenizer.SYMBOL, ";") {
		return // エラー
	}
	addTokenNode(tokenize, &varDecNode)

	subroutineBodyNode.Children = append(subroutineBodyNode.Children, varDecNode)

}

func CompileStatements(tokenize *jacktokenizer.JackTokenizer, subroutineBodyNode *ContainerNode) {

	statementsNode := ContainerNode{Name: "statements", Children: []Node{}}

	for tokenize.HasMoreTokens() {
		if expectToken(tokenize, jacktokenizer.SYMBOL, "}") {
			break
		}
		switch tokenize.GetTokenValue() {
		case "let":
			CompileLet(tokenize, &statementsNode)
		case "if":
			CompileIf(tokenize, &statementsNode)
		case "do":
			CompileDo(tokenize, &statementsNode)
		case "while":
			CompileWhile(tokenize, &statementsNode)
		case "return":
			CompileReturn(tokenize, &statementsNode)
		}

	}

	subroutineBodyNode.Children = append(subroutineBodyNode.Children, statementsNode)
}

func CompileDo(tokenize *jacktokenizer.JackTokenizer, statementsNode *ContainerNode) {
	doNode := ContainerNode{Name: "doStatement", Children: []Node{}}

	// 'do' キーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "do") {
		return // エラー
	}
	addTokenNode(tokenize, &doNode)

	// サブルーチン呼び出し
	if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	addTokenNode(tokenize, &doNode)

	// "."
	if expectToken(tokenize, jacktokenizer.SYMBOL, ".") {
		addTokenNode(tokenize, &doNode)

		// subroutine name
		if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
			return // エラー
		}
		addTokenNode(tokenize, &doNode)
	}

	// "("
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "(") {
		return // エラー
	}
	addTokenNode(tokenize, &doNode)

	// expressionList
	CompileExpressionList(tokenize, &doNode)

	// ")"
	if !expectToken(tokenize, jacktokenizer.SYMBOL, ")") {
		return // エラー
	}
	addTokenNode(tokenize, &doNode)

	// ';'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, ";") {
		return // エラー
	}
	addTokenNode(tokenize, &doNode)

	statementsNode.Children = append(statementsNode.Children, doNode)
}

func CompileLet(tokenize *jacktokenizer.JackTokenizer, statementsNode *ContainerNode) {
	letNode := ContainerNode{Name: "letStatement", Children: []Node{}}

	// 'let' キーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "let") {
		return // エラー
	}
	addTokenNode(tokenize, &letNode)

	// 変数名
	if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	addTokenNode(tokenize, &letNode)

	// 配列アクセスの場合
	if expectToken(tokenize, jacktokenizer.SYMBOL, "[") {
		addTokenNode(tokenize, &letNode) // '['
		CompileExpression(tokenize, &letNode)
		if !expectToken(tokenize, jacktokenizer.SYMBOL, "]") {
			return // エラー
		}
		addTokenNode(tokenize, &letNode) // ']'
	}

	// '='
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "=") {
		return // エラー
	}
	addTokenNode(tokenize, &letNode)

	// 式
	CompileExpression(tokenize, &letNode)

	// ';'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, ";") {
		return // エラー
	}
	addTokenNode(tokenize, &letNode)

	// 'let' ノードをステートメントノードに追加
	statementsNode.Children = append(statementsNode.Children, letNode)
}

func CompileWhile(tokenize *jacktokenizer.JackTokenizer, statementsNode *ContainerNode) {
	whileNode := ContainerNode{Name: "whileStatement", Children: []Node{}}

	// 'while' キーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "while") {
		return // エラー
	}
	addTokenNode(tokenize, &whileNode)

	// '('
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "(") {
		return // エラー
	}
	addTokenNode(tokenize, &whileNode)

	// 式
	CompileExpression(tokenize, &whileNode)

	// ')'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, ")") {
		return // エラー
	}
	addTokenNode(tokenize, &whileNode)

	// '{'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "{") {
		return // エラー
	}
	addTokenNode(tokenize, &whileNode)

	// statements
	CompileStatements(tokenize, &whileNode)

	// '}'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "}") {
		return // エラー
	}
	addTokenNode(tokenize, &whileNode)

	statementsNode.Children = append(statementsNode.Children, whileNode)
}

func CompileReturn(tokenize *jacktokenizer.JackTokenizer, statementsNode *ContainerNode) {
	returnNode := ContainerNode{Name: "returnStatement", Children: []Node{}}

	// 'return' キーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "return") {
		return // エラー
	}
	addTokenNode(tokenize, &returnNode)

	// 式 (オプション)
	if tokenize.GetTokenType() != jacktokenizer.SYMBOL || tokenize.GetTokenValue() != ";" {
		CompileExpression(tokenize, &returnNode)
	}

	// ';'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, ";") {
		return // エラー
	}
	addTokenNode(tokenize, &returnNode)

	statementsNode.Children = append(statementsNode.Children, returnNode)
}

func CompileIf(tokenize *jacktokenizer.JackTokenizer, statementsNode *ContainerNode) {
	ifNode := ContainerNode{Name: "ifStatement", Children: []Node{}}

	// 'if' キーワード
	if !expectToken(tokenize, jacktokenizer.KEYWORD, "if") {
		return // エラー
	}
	addTokenNode(tokenize, &ifNode)

	// '('
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "(") {
		return // エラー
	}
	addTokenNode(tokenize, &ifNode)

	// 式
	CompileExpression(tokenize, &ifNode)

	// ')'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, ")") {
		return // エラー
	}
	addTokenNode(tokenize, &ifNode)

	// '{'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "{") {
		return // エラー
	}
	addTokenNode(tokenize, &ifNode)

	// statements
	CompileStatements(tokenize, &ifNode)

	// '}'
	if !expectToken(tokenize, jacktokenizer.SYMBOL, "}") {
		return // エラー
	}
	addTokenNode(tokenize, &ifNode)

	// 'else' (オプション)
	if expectToken(tokenize, jacktokenizer.KEYWORD, "else") {
		addTokenNode(tokenize, &ifNode)

		// '{'
		if !expectToken(tokenize, jacktokenizer.SYMBOL, "{") {
			return // エラー
		}
		addTokenNode(tokenize, &ifNode)

		// statements
		CompileStatements(tokenize, &ifNode)

		// '}'
		if !expectToken(tokenize, jacktokenizer.SYMBOL, "}") {
			return // エラー
		}
		addTokenNode(tokenize, &ifNode)
	}

	statementsNode.Children = append(statementsNode.Children, ifNode)
}

func CompileExpression(tokenize *jacktokenizer.JackTokenizer, parentNode *ContainerNode) {
	expressionNode := ContainerNode{Name: "expression", Children: []Node{}}

	// 最初の term をコンパイル
	CompileTerm(tokenize, &expressionNode)

	// (op term)* を処理
	for tokenize.HasMoreTokens() {
		// 演算子をチェック
		if tokenize.GetTokenType() != jacktokenizer.SYMBOL || !isOperator(tokenize.GetTokenValue()) {
			break
		}

		// 演算子をノードに追加
		addTokenNode(tokenize, &expressionNode)

		// 次の term をコンパイル
		CompileTerm(tokenize, &expressionNode)
	}

	// 完成した expression ノードを親ノードに追加
	parentNode.Children = append(parentNode.Children, expressionNode)
}

// 演算子かどうかを判定するヘルパー関数
func isOperator(token string) bool {
	operators := []string{"+", "-", "*", "/", "&amp;", "|", "&gt;", "&lt;", "="}
	for _, op := range operators {
		if token == op {
			return true
		}
	}
	return false
}

func CompileTerm(tokenize *jacktokenizer.JackTokenizer, parentNode *ContainerNode) {
	termNode := ContainerNode{Name: "term", Children: []Node{}}

	switch tokenize.GetTokenType() {
	case jacktokenizer.INT_CONST:
		// 整数定数
		addTokenNode(tokenize, &termNode)

	case jacktokenizer.STRING_CONST:
		// 文字列定数
		addTokenNode(tokenize, &termNode)

	case jacktokenizer.KEYWORD:
		// キーワード定数 (true, false, null, this)
		if isKeywordConstant(tokenize.GetTokenValue()) {
			addTokenNode(tokenize, &termNode)
		} else {
			return // エラー
		}

	case jacktokenizer.IDENTIFIER:
		// 変数名またはサブルーチン呼び出し
		addTokenNode(tokenize, &termNode)

		// 配列アクセスまたはサブルーチン呼び出しを処理
		if tokenize.GetTokenType() == jacktokenizer.SYMBOL {
			switch tokenize.GetTokenValue() {
			case "[":
				// 配列アクセス
				addTokenNode(tokenize, &termNode) // '['
				CompileExpression(tokenize, &termNode)
				if !expectToken(tokenize, jacktokenizer.SYMBOL, "]") {
					return // エラー
				}
				addTokenNode(tokenize, &termNode) // ']'

			case "(":
				// サブルーチン呼び出し (引数なし)
				addTokenNode(tokenize, &termNode) // '('
				CompileExpressionList(tokenize, &termNode)
				if !expectToken(tokenize, jacktokenizer.SYMBOL, ")") {
					return // エラー
				}
				addTokenNode(tokenize, &termNode) // ')'

			case ".":
				// サブルーチン呼び出し (クラス名または変数名付き)
				addTokenNode(tokenize, &termNode) // '.'
				if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
					return // エラー
				}
				addTokenNode(tokenize, &termNode) // サブルーチン名
				if !expectToken(tokenize, jacktokenizer.SYMBOL, "(") {
					return // エラー
				}
				addTokenNode(tokenize, &termNode) // '('
				CompileExpressionList(tokenize, &termNode)
				if !expectToken(tokenize, jacktokenizer.SYMBOL, ")") {
					return // エラー
				}
				addTokenNode(tokenize, &termNode) // ')'
			}
		}

	case jacktokenizer.SYMBOL:
		switch tokenize.GetTokenValue() {
		case "(":
			// 括弧で囲まれた式
			addTokenNode(tokenize, &termNode) // '('
			CompileExpression(tokenize, &termNode)
			if !expectToken(tokenize, jacktokenizer.SYMBOL, ")") {
				return // エラー
			}
			addTokenNode(tokenize, &termNode) // ')'

		case "-", "~":
			// 単項演算子
			addTokenNode(tokenize, &termNode) // 単項演算子
			CompileTerm(tokenize, &termNode)

		default:
			return // エラー
		}

	default:
		return // エラー
	}

	// 完成した term ノードを親ノードに追加
	parentNode.Children = append(parentNode.Children, termNode)
}

// キーワード定数かどうかを判定するヘルパー関数
func isKeywordConstant(token string) bool {
	keywords := []string{"true", "false", "null", "this"}
	for _, kw := range keywords {
		if token == kw {
			return true
		}
	}
	return false
}

func CompileExpressionList(tokenize *jacktokenizer.JackTokenizer, doNode *ContainerNode) {
	expressionListNode := ContainerNode{Name: "expressionList", Children: []Node{}}

	// 式がある分繰り返す
	for tokenize.HasMoreTokens() {
		if tokenize.GetTokenType() == jacktokenizer.SYMBOL && tokenize.GetTokenValue() == ")" {
			break
		}

		// 式
		CompileExpression(tokenize, &expressionListNode)

		// ","
		if tokenize.GetTokenValue() != "," {
			if tokenize.GetTokenType() == jacktokenizer.SYMBOL && tokenize.GetTokenValue() == ")" {
				break
			} else {
				return // エラー
			}
		}
		commaNode := ParseNode{Type: tokenize.GetTokenType(), Value: tokenize.GetTokenValue()}
		expressionListNode.Children = append(expressionListNode.Children, commaNode)
		tokenize.Advance()
	}

	doNode.Children = append(doNode.Children, expressionListNode)
}

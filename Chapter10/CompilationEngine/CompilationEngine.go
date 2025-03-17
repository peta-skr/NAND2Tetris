package compilationengine

import (
	jacktokenizer "Chapter10/JackTokenizer"
)

type Node interface {
	isNode()
}

type ParseTree struct {
	Nodes []Node
}

type ParseNode struct {
	Type     jacktokenizer.TokenType
	Value    string
	Children []Node
}

func (n ParseNode) isNode() {
}

type ClassNode struct {
	Name     string
	Children []Node
}

func (n ClassNode) isNode() {
}

type ClassVarDecNode struct {
	Name     string
	Children []Node
}

func (n ClassVarDecNode) isNode() {
}

func (tree *ParseTree) AddNode(node Node) {
	tree.Nodes = append(tree.Nodes, node)
}

func Compile(tokenize jacktokenizer.JackTokenizer, outputFile string) {

	parseTree := ParseTree{}

	for tokenize.HasMoreTokens() {
		// token := tokenize.GetTokenValue()
		tokenType := tokenize.GetTokenType()

		switch tokenType {
		case jacktokenizer.KEYWORD:
			switch tokenize.GetKeyword() {
			case jacktokenizer.CLASS:
				parseTree.CompileClass(&tokenize)
			case jacktokenizer.METHOD:
			case jacktokenizer.FUNCTION:
			case jacktokenizer.CONSTRUCTOR:
			case jacktokenizer.INT:
			case jacktokenizer.BOOLEAN:
			case jacktokenizer.CHAR:
			case jacktokenizer.VOID:
			case jacktokenizer.VAR:
			case jacktokenizer.STATIC:
			case jacktokenizer.FIELD:
			case jacktokenizer.LET:
			case jacktokenizer.DO:
			case jacktokenizer.IF:
			case jacktokenizer.ELSE:
			case jacktokenizer.WHILE:
			case jacktokenizer.RETURN:
			case jacktokenizer.TRUE:
			case jacktokenizer.FALSE:
			case jacktokenizer.NULL:
			case jacktokenizer.THIS:
			}
		case jacktokenizer.IDENTIFIER:
		case jacktokenizer.INT_CONST:
		case jacktokenizer.STRING_CONST:
		case jacktokenizer.SYMBOL:

		}

		tokenize.Advance()
	}
}

func (p *ParseTree) CompileClass(tokenize *jacktokenizer.JackTokenizer) {
	classContentNode := ClassNode{Name: "class", Children: []Node{}}

	// 'class' キーワード
	if tokenize.GetTokenType() != jacktokenizer.SYMBOL && tokenize.GetKeyword() != jacktokenizer.CLASS && tokenize.GetTokenValue() != "class" {
		return // エラー
	}
	classNode := ParseNode{Type: jacktokenizer.KEYWORD, Value: "class"}
	classContentNode.Children = append(classContentNode.Children, classNode)
	tokenize.Advance()

	// クラス名
	if tokenize.GetTokenType() != jacktokenizer.IDENTIFIER {
		return // エラー
	}
	classNameNode := ParseNode{Type: jacktokenizer.IDENTIFIER, Value: tokenize.GetTokenValue()}
	classContentNode.Children = append(classContentNode.Children, classNameNode)
	tokenize.Advance()

	// '{' シンボル
	if tokenize.GetTokenType() != jacktokenizer.SYMBOL && tokenize.GetTokenValue() != "{" {
		return // エラー
	}
	openBraceNode := ParseNode{Type: jacktokenizer.SYMBOL, Value: "{"}
	classContentNode.Children = append(classContentNode.Children, openBraceNode)
	tokenize.Advance()

	for tokenize.HasMoreTokens() {
		tokenType := tokenize.GetTokenType()
		if tokenType == jacktokenizer.SYMBOL && tokenize.GetTokenValue() == "}" {
			break
		}

		switch tokenType {
		case jacktokenizer.KEYWORD:
			switch tokenize.GetKeyword() {
			case jacktokenizer.STATIC, jacktokenizer.FIELD:
				CompileClassVarDec(tokenize, &classContentNode)
			case jacktokenizer.METHOD, jacktokenizer.FUNCTION, jacktokenizer.CONSTRUCTOR:
				CompileSubroutine(tokenize, &classContentNode)
			}
		}
		tokenize.Advance()
	}

	p.AddNode(classContentNode)
}

func CompileClassVarDec(tokenize *jacktokenizer.JackTokenizer, classContentNode *ClassNode) {

	classVarDecNode := ClassVarDecNode{Name: "classVarDec", Children: []Node{}}

	// 'static' or 'field' キーワード
	if tokenize.GetTokenType() != jacktokenizer.KEYWORD &&
		(tokenize.GetKeyword() != jacktokenizer.STATIC &&
			tokenize.GetTokenValue() != "static") ||
		(tokenize.GetKeyword() != jacktokenizer.FIELD &&
			tokenize.GetTokenValue() != "field") {
		return // エラー
	}
	classNode := ParseNode{Type: tokenize.GetTokenType(), Value: tokenize.GetTokenValue()}
	classVarDecNode.Children = append(classVarDecNode.Children, classNode)
	tokenize.Advance()

	// 型
	if tokenize.GetTokenType() != jacktokenizer.KEYWORD &&
		(tokenize.GetKeyword() != jacktokenizer.INT &&
			tokenize.GetTokenValue() != "int") ||
		(tokenize.GetKeyword() != jacktokenizer.CHAR &&
			tokenize.GetTokenValue() != "char") ||
		(tokenize.GetKeyword() != jacktokenizer.BOOLEAN &&
			tokenize.GetTokenValue() != "boolean") ||
		(tokenize.GetTokenType() != jacktokenizer.IDENTIFIER) {
		return // エラー
	}

}

func CompileSubroutine(tokenize *jacktokenizer.JackTokenizer, classVarDecNode *ClassNode) {
	node := ParseNode{Type: jacktokenizer.KEYWORD, Value: tokenize.GetTokenValue(), Children: []Node{}}

}

// func CompileParameterList() {

// }

// func CompileVarDec() {

// }

// func CompileStatements() {

// }

// func CompileDo() {

// }

// func CompileLet() {

// }

// func CompileWhile() {

// }

// func CompileReturn() {

// }

// func CompileIf() {

// }

// func CompileExpression() {

// }

// func CompileTerm() {

// }

// func CompileExpressionList() {

// }

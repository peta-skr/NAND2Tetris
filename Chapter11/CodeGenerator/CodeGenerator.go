package codegenerator

import (
	compilationengine "Chapter11/CompilationEngine"
	jacktokenizer "Chapter11/JackTokenizer"
	symboltable "Chapter11/SymbolTable"
	symtable "Chapter11/SymbolTable"
	vmwriter "Chapter11/VmWriter"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var className string
var subroutineName string
var labelCount int
var negFlag bool

func GenerateCode(parseTree compilationengine.ParseTree, symboltable symboltable.SymbolTable, vmFilePath string) {

	vmwriter := vmwriter.Constructor()

	for _, node := range parseTree.Nodes {
		processNode(node, symboltable, vmFilePath, &vmwriter)
	}

	// VMWriterの内容をファイルに書き込む
	vmContent := vmwriter.Content                    // []string 型
	vmContentString := strings.Join(vmContent, "\n") // 改行区切りの文字列に変換

	// ファイルに書き込む
	if err := os.WriteFile(vmFilePath, []byte(vmContentString), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "VMファイルの書き込みに失敗しました: %v\n", err)
	} else {
		fmt.Printf("VMファイル '%s' が正常に作成されました。\n", vmFilePath)
	}
}

func processNode(node compilationengine.Node, symboltable symtable.SymbolTable, vmFilePath string, vmwriter *vmwriter.VMWriter) {
	// ノードの種類に応じて処理を分岐
	switch n := node.(type) {
	case compilationengine.ContainerNode:
		switch n.Name {
		case "class":
			generateClassCode(n, symboltable, vmwriter)
		case "subroutineDec":
		}
	case compilationengine.ParseNode:
		// 必要に応じてシンボルテーブルやVMコード生成処理を追加
	}
}

func generateClassCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// クラスのコード生成処理
	// クラス名を取得

	symboltable.CurrentScope = symtable.CLASS_SCOPE

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "subroutineDec":
				generateSubroutineDecCode(n, symboltable, vmwriter)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER {
				// クラス名を取得
				className = n.Value
			}
		}
	}

	// クラスのシンボルテーブルを取得
	// classSymbolTable := symboltable.ClassSymbolTable[className]

	// VMWriterにクラスのコードを書き込む
	// vmwriter.WriteClass(className, classSymbolTable)
}

func generateSubroutineDecCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// サブルーチンのコード生成処理
	// サブルーチン名を取得
	subroutineName = node.Name

	symboltable.CurrentScope = symtable.SUBROUTINE_SCOPE

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "parameterList":
				// パラメータリストの処理
				generateParameterListCode(n, symboltable, vmwriter)
			case "subroutineBody":
				generateSubroutineBodyCode(n, symboltable, vmwriter)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER {
				// サブルーチン名を取得
				subroutineName = n.Value
				functionName := fmt.Sprintf("%s.%s", className, subroutineName)
				vmwriter.WriteFunction(functionName, symboltable.VarCount(subroutineName, "var"))
			}
		}
	}

	// サブルーチンのシンボルテーブルを取得
	// subroutineSymbolTable := symboltable.SubroutineSymbolTable[subroutineName]

	// VMWriterにサブルーチンのコードを書き込む
	// vmwriter.WriteSubroutine(subroutineName, subroutineSymbolTable)
}

func generateSubroutineBodyCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// サブルーチンボディのコード生成処理
	// ローカル変数の処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "varDec":
				// ローカル変数の処理
				generateVarDecCode(n, symboltable, vmwriter)
			case "statements":
				// ステートメントの処理
				generateStatementsCode(n, symboltable, vmwriter)
			}
		}
	}

	// VMWriterにサブルーチンボディのコードを書き込む
	// vmwriter.WriteSubroutineBody(subroutineName, localSymbolTable)
}

func generateParameterListCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// パラメータリストのコード生成処理
	// パラメータの処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER {
				// vmwriter.WritePush("argument", symboltable.IndexOf(subroutineName, n.Value))
			}
		}
	}

	// VMWriterにパラメータリストのコードを書き込む
	// vmwriter.WriteParameterList(subroutineName, parameterSymbolTable)
}

// 特に何もしない
func generateVarDecCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// ローカル変数のコード生成処理
	// ローカル変数の処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "var":
				// ローカル変数の処理
			}
		}
	}

	// VMWriterにローカル変数のコードを書き込む
	// vmwriter.WriteVarDec(subroutineName, localSymbolTable)
}

func generateStatementsCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// ステートメントのコード生成処理
	// ステートメントの処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "letStatement":
				// let文の処理
				generateLetStatementCode(n, symboltable, vmwriter)
			case "ifStatement":
				// if文の処理
				generateIfStatementCode(n, symboltable, vmwriter)
			case "whileStatement":
				// while文の処理
				generateWhileStatementCode(n, symboltable, vmwriter)
			case "doStatement":
				// do文の処理
				generateDoStatementCode(n, symboltable, vmwriter)
			case "returnStatement":
				// return文の処理
				generateReturnStatementCode(n, symboltable, vmwriter)
			}
		}
	}

	// VMWriterにステートメントのコードを書き込む
	// vmwriter.WriteStatements(subroutineName, statementSymbolTable)
}

func generateStatementCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// ステートメントのコード生成処理
	// ステートメントの処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "letStatement":
				// let文の処理
				generateLetStatementCode(n, symboltable, vmwriter)
			case "ifStatement":
				// if文の処理
				generateIfStatementCode(n, symboltable, vmwriter)
			case "whileStatement":
				// while文の処理
				generateWhileStatementCode(n, symboltable, vmwriter)
			case "doStatement":
				// do文の処理
				generateDoStatementCode(n, symboltable, vmwriter)
			case "returnStatement":
				// return文の処理
				generateReturnStatementCode(n, symboltable, vmwriter)
			}
		}
	}

	// VMWriterにステートメントのコードを書き込む
	// vmwriter.WriteStatement(subroutineName, statementSymbolTable)
}

func generateLetStatementCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {

	varName := ""

	fmt.Println("let:", node)

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expression":
				// 式の処理
				// ここでは式のコード生成処理を呼び出す
				generateExpressionCode(n, symboltable, vmwriter)
			}

		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER {
				// 変数名を取得
				varName = n.Value
				// // 変数のインデックスを取得
				// fmt.Println("varName", varName)
				// index := symboltable.IndexOf(subroutineName, varName)
				// // VMWriterにlet文のコードを書き込む
				// vmwriter.WritePush("local", index)
			} else if n.Type == jacktokenizer.SYMBOL && n.Value == ";" {
				// 変数のインデックスを取得
				index := symboltable.IndexOf(subroutineName, varName)
				// let文の代入処理
				kind := getKind(subroutineName, varName, symboltable)
				vmwriter.WritePop(kind, index)
			}
		}

	}
}

func generateIfStatementCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// if文のコード生成処理
	// 条件式を取得
	// condition := node.Name
	labelCountStr := ""
	ifLabel := ""
	elseLabel := ""

	fmt.Println("if:", node)

	labelCountStr = strconv.Itoa(labelCount)
	elseLabel = className + "_" + labelCountStr
	labelCount++
	labelCountStr = strconv.Itoa(labelCount)
	ifLabel = className + "_" + labelCountStr
	labelCount++

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expression":
				// 式の処理
				// ここでは式のコード生成処理を呼び出す
				generateExpressionCode(n, symboltable, vmwriter)
				vmwriter.WriteArithmetic("not") // 条件式の否定 → goto命令は、条件式がfalseのときに実行されるため、notを使う
				vmwriter.WriteIf(ifLabel)       // 条件式がtrueの場合に実行される

			case "statements":
				// ステートメントの処理
				generateStatementsCode(n, symboltable, vmwriter)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER {
				// 変数名を取得
				// varName := n.Value
				// 変数のインデックスを取得
				// index := symboltable.IndexOf(subroutineName, varName)
			} else if n.Type == jacktokenizer.KEYWORD && n.Value == "else" {
				// else文の処理
				// VMWriterにelse文のコードを書き込む
				// vmwriter.WriteGoto(condition)
				vmwriter.WriteGoto(elseLabel) // 条件式がtrueの場合に実行される
				vmwriter.WriteLabel(ifLabel)  // 条件式がfalseの場合に実行される
			} else if n.Type == jacktokenizer.KEYWORD && n.Value == "if" {
				// VMWriterにif文のコードを書き込む
				// label := fmt.Sprintf("%s.%d", className, labelCount)
				labelCountStr = strconv.Itoa(labelCount)
			}
		}

	}
	vmwriter.WriteLabel(elseLabel) // else文の終了ラベル
}

func generateWhileStatementCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// while文のコード生成処理
	// 条件式を取得

	whileStart := ""
	whileEnd := ""

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expression":
				// 式の処理
				// ここでは式のコード生成処理を呼び出す
				generateExpressionCode(n, symboltable, vmwriter)
				vmwriter.WriteArithmetic("not") // 条件式の否定 → goto命令は、条件式がfalseのときに実行されるため、notを使う
			case "statements":
				// ステートメントの処理
				generateStatementsCode(n, symboltable, vmwriter)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.KEYWORD && n.Value == "while" {
				// while文の開始
				// label := fmt.Sprintf("%s.%d", className, labelCount)
				// VMWriterにwhile文の開始のコードを書き込む
				labelCountStr := strconv.Itoa(labelCount)
				whileStart = className + "_" + labelCountStr
				vmwriter.WriteLabel(whileStart)
				labelCount++
			} else if n.Type == jacktokenizer.SYMBOL && n.Value == ")" {
				// while文の条件式の終了
				labelCountStr := strconv.Itoa(labelCount)
				whileEnd = className + "_" + labelCountStr
				vmwriter.WriteIf(whileEnd)
				labelCount++
			} else if n.Type == jacktokenizer.SYMBOL && n.Value == "}" {
				// while文の終了
				// VMWriterにwhile文の終了のコードを書き込む
				vmwriter.WriteGoto(whileStart)
				vmwriter.WriteLabel(whileEnd)
			}
		}
	}

}

func generateDoStatementCode(node compilationengine.ContainerNode, symboltable symboltable.SymbolTable, vmwriter *vmwriter.VMWriter) {

	doName := ""
	argNum := 0

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expressionList":
				// 引数リストの処理
				// ここでは引数リストのコード生成処理を呼び出す
				argNum = generateExpressionListCode(n, symboltable, vmwriter)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER {
				// 関数名を取得
				doName += n.Value
			} else if n.Type == jacktokenizer.SYMBOL && n.Value == "." {
				// 引数リストの開始
				doName += "."
			} else if n.Type == jacktokenizer.SYMBOL && n.Value == ")" {
				// 引数リストの終了
				vmwriter.WriteCall(doName, argNum)
			}
		}
	}
	vmwriter.WritePop("temp", 0)
}

func generateReturnStatementCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// VMWriterにreturn文のコードを書き込む
	length := len(node.Children)
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expression":
				// 式の処理
				// ここでは式のコード生成処理を呼び出す
				generateExpressionCode(n, symboltable, vmwriter)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.KEYWORD && n.Value == "return" {
				// return文の処理
				// VMWriterにreturn文のコードを書き込む
				if length < 3 {
					vmwriter.WritePush("constant", 0)
					vmwriter.WriteReturn()
				}
			}
		}
	}
	if length > 2 {
		// VMWriterにreturn文のコードを書き込む
		vmwriter.WriteReturn()
	}
}

func generateExpressionListCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) int {

	argNum := 0
	if len(node.Children) != 0 {
		argNum++
	}
	// 引数リストのコード生成処理
	// 引数の処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expression":
				// 式の処理
				// ここでは式のコード生成処理を呼び出す
				generateExpressionCode(n, symboltable, vmwriter)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.SYMBOL && n.Value == "," {
				// 引数の区切り
				argNum++
			}
		}
	}

	// VMWriterに引数リストのコードを書き込む
	// vmwriter.WriteExpressionList(subroutineName, expressionListSymbolTable)

	return argNum
}

func generateExpressionCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// 式のコード生成処理
	// 式の処理

	operations := make([]string, 0)

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "term":
				// 項の処理
				operations = generateTermCode(n, symboltable, vmwriter, operations)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.SYMBOL {
				// 演算子を取得
				operator := n.Value
				operations = append(operations, operator)
				// VMWriterに演算子のコードを書き込む
			}
		}

	}

	for _, op := range operations {
		switch op {
		case "+":
			vmwriter.WriteArithmetic("add")
		case "-":
			vmwriter.WriteArithmetic("sub")
		case "*":
			vmwriter.WriteArithmetic("call Math.multiply 2")
		case "/":
			vmwriter.WriteArithmetic("call Math.divide 2")
		case "&amp;":
			vmwriter.WriteArithmetic("and")
		case "|":
			vmwriter.WriteArithmetic("or")
		case "&lt;":
			vmwriter.WriteArithmetic("lt")
		case "&gt;":
			vmwriter.WriteArithmetic("gt")
		case "=":
			vmwriter.WriteArithmetic("eq")
		case "~":
			// ビット反転演算子の処理
			vmwriter.WriteArithmetic("not")
		}
	}

	// VMWriterに式のコードを書き込む
	// vmwriter.WriteExpression(subroutineName, expressionSymbolTable)
}

func generateTermCode(node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter, operations []string) []string {
	// 項のコード生成処理
	// 項の処理

	methodName := ""
	for i, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ParseNode:
			intValue, _ := strconv.Atoi(n.Value)
			switch n.Type {
			case jacktokenizer.INT_CONST:
				// 整数定数の処理
				vmwriter.WritePush("constant", intValue)
				if negFlag {
					// 負の数を処理
					vmwriter.WriteArithmetic("neg")
					negFlag = false
				}
			case jacktokenizer.STRING_CONST:
				// 文字列定数の処理
				vmwriter.WritePush("constant", intValue)
			case jacktokenizer.KEYWORD:
				// キーワード定数の処理
				switch n.Value {
				case "true":
					// trueの場合は-1をpush
					vmwriter.WritePush("constant", 1)
					vmwriter.WriteArithmetic("neg")
				case "false":
					// falseの場合は0をpush
					vmwriter.WritePush("constant", 0)
				}
			case jacktokenizer.IDENTIFIER:
				// 識別子の処理
				if i+1 < len(node.Children) {
					// 次のノードを確認してメソッド呼び出しかどうかを判定
					if nextNode, ok := node.Children[i+1].(compilationengine.ParseNode); ok && nextNode.Type == jacktokenizer.SYMBOL && nextNode.Value == "(" {
						// メソッド呼び出しの場合
						methodName += n.Value
						generateSubroutineCall(methodName, node, symboltable, vmwriter)
					} else if nextNode, ok := node.Children[i+1].(compilationengine.ParseNode); ok && nextNode.Type == jacktokenizer.SYMBOL && nextNode.Value == "." {
						methodName += n.Value + "."
					} else {
						// 通常の変数として処理
						kind := getKind(subroutineName, n.Value, symboltable)
						vmwriter.WritePush(kind, symboltable.IndexOf(subroutineName, n.Value))
					}
				} else {
					// 通常の変数として処理
					kind := getKind(subroutineName, n.Value, symboltable)
					vmwriter.WritePush(kind, symboltable.IndexOf(subroutineName, n.Value))
				}
			case jacktokenizer.SYMBOL:
				if n.Value == "-" {
					// // 負の数を処理
					// if i+1 < len(node.Children) {
					// 	if nextNode, ok := node.Children[i+1].(compilationengine.ContainerNode); ok && nextNode.Name == "term" {
					// 		// 次のノードが整数定数の場合
					// 		intValue, _ := strconv.Atoi(nextNode.Children[0].(compilationengine.ParseNode).Value)
					// 		vmwriter.WritePush("constant", intValue)
					// 		vmwriter.WriteArithmetic("neg")
					// 		break
					// 	}
					// }
					// // 通常の負号演算子として処理
					// vmwriter.WriteArithmetic("neg")
					negFlag = true

				} else if n.Value == "~" {
					// ビット反転演算子の処理
					operations = append(operations, n.Value)
				}
			}
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expression":
				// 式の処理
				generateExpressionCode(n, symboltable, vmwriter)
			case "expressionList":
				// 引数リストの処理
				// generateExpressionListCode(n, symboltable, vmwriter)
			case "term":
				// 再帰的に項の処理を呼び出す
				generateTermCode(n, symboltable, vmwriter, operations)
			}
		}
	}

	// VMWriterに項のコードを書き込む
	// vmwriter.WriteTerm(subroutineName, termSymbolTable)

	return operations
}

func getKind(subroutineName string, name string, symboltable symtable.SymbolTable) string {
	kind := symboltable.KindOf(subroutineName, name)
	switch kind {
	case "static":
		return "static"
	case "field":
		return "field"
	case "arg":
		return "argument"
	case "var":
		return "local"
	default:
		return ""
	}
	return ""
}

func generateSubroutineCall(methodName string, node compilationengine.ContainerNode, symboltable symtable.SymbolTable, vmwriter *vmwriter.VMWriter) {
	// 引数リストの処理
	argCount := 0
	for _, child := range node.Children {
		if n, ok := child.(compilationengine.ContainerNode); ok && n.Name == "expressionList" {
			argCount = generateExpressionListCode(n, symboltable, vmwriter)
		}
	}

	// メソッド呼び出しのコードを生成
	vmwriter.WriteCall(methodName, argCount)
	// vmwriter.WritePop("local", 0)
}

package codegenerator

import (
	compilationengine "Chapter11/CompilationEngine"
	jacktokenizer "Chapter11/JackTokenizer"
	symboltable "Chapter11/SymbolTable"
	vmwriter "Chapter11/VmWriter"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// コード生成に必要な状態を保持する構造体
type CodeGenerator struct {
	className         string
	subroutineName    string
	labelCount        int
	negFlag           bool
	isMethodCall      bool
	subroutineKindMap map[string]map[string]string // クラス名ごとにサブルーチンの種類を記録
	symboltable       symboltable.SymbolTable
	vmwriter          vmwriter.VMWriter
}

const (
	CLASS_SCOPE      = symboltable.CLASS_SCOPE
	SUBROUTINE_SCOPE = symboltable.SUBROUTINE_SCOPE
)

// 新しいCodeGeneratorを作成
func New(symbolTable symboltable.SymbolTable, subroutineKind map[string]map[string]string) *CodeGenerator {
	return &CodeGenerator{
		labelCount:        0,
		subroutineKindMap: subroutineKind, // クラス名ごとにサブルーチンの種類を記録
		symboltable:       symbolTable,
		vmwriter:          vmwriter.Constructor(),
	}
}

// VMコードを生成して指定されたパスに保存
func (cg *CodeGenerator) Generate(parseTree compilationengine.ParseTree, vmFilePath string) error {
	// パースツリーの処理
	for _, node := range parseTree.Nodes {
		cg.processNode(node)
	}

	return cg.writeToFile(vmFilePath)
}

// VMコードをファイルに書き込む
func (cg *CodeGenerator) writeToFile(filePath string) error {
	vmContent := cg.vmwriter.Content                   // []string 型
	vmContentString := strings.Join(vmContent, "\r\n") // 改行区切りの文字列に変換

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("ファイルを開けません: %v", err)
	}
	defer file.Close()

	// ファイルに書き込む
	if _, err := file.WriteString(vmContentString); err != nil {
		return fmt.Errorf("VMファイルの書き込みに失敗: %v", err)
	}

	fmt.Printf("VMファイル '%s' を作成しました。\n", filePath)
	return nil
}

func (cg *CodeGenerator) processNode(node compilationengine.Node) {
	// ノードの種類に応じて処理を分岐
	switch n := node.(type) {
	case compilationengine.ContainerNode:
		switch n.Name {
		case "class":
			cg.generateClassCode(n)
			// case "subroutineDec":
			// 	cg.generateSubroutineDecCode(n)
		}
	}
}

// クラスのコード生成
func (cg *CodeGenerator) generateClassCode(node compilationengine.ContainerNode) {
	// クラス名を取得

	cg.symboltable.CurrentScope = symboltable.CLASS_SCOPE

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "subroutineDec":
				cg.generateSubroutineDecCode(n)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER {
				// クラス名を取得
				cg.className = n.Value
			}
		}
	}

}

// サブルーチン宣言からVMコードを生成する
func (cg *CodeGenerator) generateSubroutineDecCode(node compilationengine.ContainerNode) {
	// スコープを設定
	cg.symboltable.CurrentScope = symboltable.SUBROUTINE_SCOPE

	// サブルーチンの情報を収集
	var subroutineType, subroutineName string

	// サブルーチン情報を収集する
	cg.collectSubroutineInfo(node, &subroutineType, &subroutineName)

	// サブルーチン情報が揃ったらVMコードを生成
	if subroutineName != "" {
		// サブルーチンの宣言を生成
		cg.generateSubroutineDeclaration(subroutineName)

		// メソッドやコンストラクタの場合の特別な初期化処理
		cg.handleSpecialSubroutineType(subroutineType)
	}

	// サブルーチンボディの処理
	for _, child := range node.Children {
		if n, ok := child.(compilationengine.ContainerNode); ok && n.Name == "subroutineBody" {
			cg.generateSubroutineBodyCode(n)
		}
	}
}

// サブルーチン情報を収集するヘルパー関数
func (cg *CodeGenerator) collectSubroutineInfo(node compilationengine.ContainerNode, subroutineType *string, subroutineName *string) {
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER && cg.className != n.Value {
				// サブルーチン名を取得
				*subroutineName = n.Value
				cg.subroutineName = n.Value
			} else if n.Type == jacktokenizer.KEYWORD {
				switch n.Value {
				case "function", "method", "constructor":
					*subroutineType = n.Value

					// サブルーチンの種類を設定
					if cg.subroutineKindMap[cg.className] == nil {
						cg.subroutineKindMap[cg.className] = make(map[string]string)
					}
					cg.subroutineKindMap[cg.className][*subroutineName] = *subroutineType

					// メソッドフラグ設定
					cg.isMethodCall = (*subroutineType == "method")
				}
			}
		}
	}
}

// サブルーチン宣言のVMコードを生成するヘルパー関数
func (cg *CodeGenerator) generateSubroutineDeclaration(subroutineName string) {
	// サブルーチン名をクラス名と結合
	functionName := fmt.Sprintf("%s.%s", cg.className, subroutineName)

	// ローカル変数の数を取得
	localVarCount := cg.symboltable.VarCount(subroutineName, "var")

	// サブルーチン宣言コード生成
	cg.vmwriter.WriteFunction(functionName, localVarCount)
}

// メソッドやコンストラクタの場合の特別処理
func (cg *CodeGenerator) handleSpecialSubroutineType(subroutineType string) {
	switch subroutineType {
	case "constructor":
		// コンストラクタの場合、オブジェクトのメモリ確保
		cg.symboltable.CurrentScope = symboltable.CLASS_SCOPE
		cg.vmwriter.WritePush("constant", cg.symboltable.VarCount(cg.subroutineName, "field"))
		cg.vmwriter.WriteCall("Memory.alloc", 1)
		cg.vmwriter.WritePop("pointer", 0)
		cg.symboltable.CurrentScope = symboltable.SUBROUTINE_SCOPE
	case "method":
		// メソッドの場合、thisをセット
		cg.vmwriter.WritePush("argument", 0)
		cg.vmwriter.WritePop("pointer", 0)
	}
}

// サブルーチンボディのコード生成
func (cg *CodeGenerator) generateSubroutineBodyCode(node compilationengine.ContainerNode) {
	// ローカル変数の処理
	for _, child := range node.Children {
		if n, ok := child.(compilationengine.ContainerNode); ok && n.Name == "statements" {
			// ステートメントの処理
			cg.generateStatementsCode(n)
		}
	}

}

// ステートメントのコード生成処理
func (cg *CodeGenerator) generateStatementsCode(node compilationengine.ContainerNode) {
	for _, child := range node.Children {
		if n, ok := child.(compilationengine.ContainerNode); ok {
			switch n.Name {
			case "letStatement":
				cg.generateLetStatementCode(n)
			case "ifStatement":
				cg.generateIfStatementCode(n)
			case "whileStatement":
				cg.generateWhileStatementCode(n)
			case "doStatement":
				cg.generateDoStatementCode(n)
			case "returnStatement":
				cg.generateReturnStatementCode(n)
			}
		}
	}

}

// let文のコード生成処理
func (cg *CodeGenerator) generateLetStatementCode(node compilationengine.ContainerNode) {
	// 代入先の変数情報を収集
	var varName string
	isArrayAssignment := false

	// 最初のパスで変数名を取得し、必要ならば配列アクセスを処理する
	for _, child := range node.Children {
		if parseNode, ok := child.(compilationengine.ParseNode); ok {
			if parseNode.Type == jacktokenizer.IDENTIFIER {
				varName = parseNode.Value
			} else if parseNode.Type == jacktokenizer.SYMBOL && parseNode.Value == "[" {
				// 配列のインデックスアクセス開始を検出
				isArrayAssignment = true
			} else if parseNode.Type == jacktokenizer.SYMBOL && parseNode.Value == "]" {
				// 配列のベースアドレスをスタックにプッシュ
				kind := cg.getSegmentName(varName)
				index := cg.symboltable.IndexOf(cg.subroutineName, varName)
				cg.vmwriter.WritePush(kind, index)
				// 配列インデックス式が評価された後
				cg.vmwriter.WriteArithmetic("add") // ベースアドレス + インデックス
			}
		} else if containerNode, ok := child.(compilationengine.ContainerNode); ok && containerNode.Name == "expression" {
			// 式が見つかる前にすでに配列のアドレス計算が終わっていれば、一時的にスタックに保存
			if isArrayAssignment {
				cg.generateExpressionCode(containerNode)
			} else {
				// 右辺の式を評価
				cg.generateExpressionCode(containerNode)
			}
		}
	}

	// 代入の処理
	if isArrayAssignment {
		cg.assignToArrayElement()
	} else {
		cg.assignToVariable(varName)
	}
}

// 配列要素への代入処理
func (cg *CodeGenerator) assignToArrayElement() {
	// 配列要素に対する代入処理
	cg.vmwriter.WritePop("temp", 0)    // 値を一時的に保存
	cg.vmwriter.WritePop("pointer", 1) // 計算されたアドレスをthatポインタにセット
	cg.vmwriter.WritePush("temp", 0)   // 値を再度取得
	cg.vmwriter.WritePop("that", 0)    // 計算されたアドレスに値を格納
}

// 通常変数への代入処理
func (cg *CodeGenerator) assignToVariable(varName string) {
	// 通常変数への代入処理
	kind := cg.getSegmentName(varName)
	index := cg.symboltable.IndexOf(cg.subroutineName, varName)
	cg.vmwriter.WritePop(kind, index) // 変数に値を格納
}

// if文のコード生成処理
func (cg *CodeGenerator) generateIfStatementCode(node compilationengine.ContainerNode) {
	// ラベル名を生成
	elseLabel := cg.generateUniqueLabel()
	ifLabel := cg.generateUniqueLabel()

	hasElse := false

	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expression":
				// 条件式の処理
				cg.generateExpressionCode(n)

				cg.vmwriter.WriteArithmetic("not") // 条件式の否定 → goto命令は、条件式がfalseのときに実行されるため、notを使う
				cg.vmwriter.WriteIf(ifLabel)       // 条件式がtrueの場合に実行される

			case "statements":
				// ステートメントの処理
				cg.generateStatementsCode(n)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.KEYWORD && n.Value == "else" {
				// else文の処理
				// VMWriterにelse文のコードを書き込む
				cg.vmwriter.WriteGoto(elseLabel) // 条件式がtrueの場合に実行される
				cg.vmwriter.WriteLabel(ifLabel)  // 条件式がfalseの場合に実行される
				hasElse = true                   // else文がある場合は、hasElseをtrueにする
			}
		}

	}

	if !hasElse {
		// else文がない場合は、if文の終了ラベルを書き込む
		cg.vmwriter.WriteGoto(elseLabel) // else文がない場合は、if文の終了ラベルを書き込む
		cg.vmwriter.WriteLabel(ifLabel)  // 条件式がfalseの場合に実行される
	}

	cg.vmwriter.WriteLabel(elseLabel) // else文の終了ラベル
}

// while文のコード生成処理
func (cg *CodeGenerator) generateWhileStatementCode(node compilationengine.ContainerNode) {
	// ラベル名を生成
	whileStartLabel := cg.generateUniqueLabel() // ループの先頭ラベル
	whileEndLabel := cg.generateUniqueLabel()   // ループの終了ラベル

	// while文の開始ラベルを設定
	cg.vmwriter.WriteLabel(whileStartLabel)

	// while文の各部分を処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			switch n.Name {
			case "expression":
				// 条件式を評価
				cg.generateExpressionCode(n)

				// 条件が偽の場合、ループ終了へジャンプ
				cg.vmwriter.WriteArithmetic("not") // 条件を反転
				cg.vmwriter.WriteIf(whileEndLabel) // 条件が偽ならループ終了へ

			case "statements":
				// ループ本体のステートメントを処理
				cg.generateStatementsCode(n)
			}
		}
	}

	// ループの先頭へ無条件ジャンプ
	cg.vmwriter.WriteGoto(whileStartLabel)

	// ループ終了ラベル
	cg.vmwriter.WriteLabel(whileEndLabel)
}

// do文のコード生成処理
func (cg *CodeGenerator) generateDoStatementCode(node compilationengine.ContainerNode) {

	// サブルーチン呼び出しの情報を収集
	callerObject := ""    // 呼び出し元オブジェクト（存在する場合）
	callName := ""        // 呼び出すサブルーチン名（完全修飾名）
	argCount := 0         // 引数の数
	isMethodCall := false // メソッド呼び出しかどうか

	// do文の各部分を処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			if n.Name == "expressionList" {
				// 引数リストを処理して引数の数を取得
				argCount += cg.generateExpressionListCode(n)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.IDENTIFIER {
				if callName == "" && callerObject == "" {
					// 最初の識別子はメソッド名またはオブジェクト名として扱う
					callerObject = n.Value
					callName = n.Value
				} else {
					// 2つ目以降の識別子はサブルーチン名
					callName += n.Value
				}
			} else if n.Type == jacktokenizer.SYMBOL {
				switch n.Value {
				case ".":
					callName += "."
					isMethodCall = true
				case "(":
					// 引数リストの開始
					if !isMethodCall {
						cg.vmwriter.WritePush("pointer", 0) // thisをポインタ0にプッシュ
					} else {

						// 外部オブジェクトのメソッド呼び出し
						objectVarName := callerObject

						// メソッド呼び出しの場合、オブジェクト名を取得
						kind := cg.getSegmentName(objectVarName)
						if kind != "" {
							index := cg.symboltable.IndexOf(cg.subroutineName, objectVarName)
							typeName := cg.symboltable.TypeOf(cg.subroutineName, objectVarName)

							// 呼び出すサブルーチン名からクラス名を取得
							idx := strings.LastIndex(callName, ".")
							if idx != -1 && idx+1 < len(callName) {
								callName = callName[idx+1:]
							} else {
								fmt.Println("区切り文字が見つからないか、ピリオドの後に文字がありません。")
							}
							callName = typeName + "." + callName
							cg.vmwriter.WritePush(kind, index)
						}
					}
				case ")":
					if !isMethodCall {
						// メソッド呼び出しの場合
						callName = cg.className + "." + callName
						argCount++ // thisを引数としてカウント
					}

					subroutineClassName := strings.Split(callName, ".")[0]
					callSubroutineName := callName[strings.LastIndex(callName, ".")+1:]
					if argCount == 0 && cg.subroutineKindMap[subroutineClassName][callSubroutineName] == "method" {
						// 引数がない場合は1を設定
						argCount = 1
					} else if cg.className != subroutineClassName && cg.subroutineKindMap[subroutineClassName][callSubroutineName] == "method" {
						// メソッド呼び出しの場合、thisを引数としてカウント
						argCount++
					}
					callName = strings.ToUpper(string(callName[0])) + callName[1:]
					cg.vmwriter.WriteCall(callName, argCount)
				}
			}
		}
	}
	cg.vmwriter.WritePop("temp", 0)
}

// return文のコード生成処理
func (cg *CodeGenerator) generateReturnStatementCode(node compilationengine.ContainerNode) {
	// return文が式を含むかどうか調べる
	hasExpression := false

	for _, child := range node.Children {
		if containerNode, ok := child.(compilationengine.ContainerNode); ok && containerNode.Name == "expression" {
			cg.generateExpressionCode(containerNode)
			hasExpression = true
			break // returnは最大1つの式しか持たないので早期脱出
		}
	}

	// 式がない場合（void戻り値の関数）は0を返す
	// Jack言語では全ての関数が値を返す必要がある
	if !hasExpression {
		cg.vmwriter.WritePush("constant", 0)
	}

	cg.vmwriter.WriteReturn()
}

// 引数リストのコード生成と引数数のカウント
func (cg *CodeGenerator) generateExpressionListCode(node compilationengine.ContainerNode) int {
	// 空の式リストの場合は引数なし
	if len(node.Children) == 0 {
		return 0
	}

	// 式の数をカウント（最初の式は確定で1つ）
	argCount := 1

	// 子ノードを走査して式とカンマを処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			if n.Name == "expression" {
				// 式を評価してその結果をスタックに積む
				cg.generateExpressionCode(n)
			}

		case compilationengine.ParseNode:
			// カンマはそれに続く別の式の存在を意味する
			if n.Type == jacktokenizer.SYMBOL && n.Value == "," {
				argCount++
			}
		}
	}

	return argCount
}

// 式のコード生成処理
func (cg *CodeGenerator) generateExpressionCode(node compilationengine.ContainerNode) {
	// 演算子とその順序を記録する配列
	operations := make([]string, 0)

	// 式の各要素を処理
	for _, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ContainerNode:
			if n.Name == "term" {
				// 項を評価してVM命令を生成し、演算子リストを更新
				operations = cg.generateTermCode(n, operations)
			}
		case compilationengine.ParseNode:
			if n.Type == jacktokenizer.SYMBOL {
				// 演算子をリストに追加
				operations = append(operations, n.Value)
			}
		}
	}

	// 演算子に基づいてVM命令を生成
	for _, op := range operations {
		switch op {
		case "+":
			cg.vmwriter.WriteArithmetic("add")
		case "-":
			cg.vmwriter.WriteArithmetic("sub")
		case "*":
			cg.vmwriter.WriteArithmetic("call Math.multiply 2")
		case "/":
			cg.vmwriter.WriteArithmetic("call Math.divide 2")
		case "&amp;":
			cg.vmwriter.WriteArithmetic("and")
		case "|":
			cg.vmwriter.WriteArithmetic("or")
		case "&lt;":
			cg.vmwriter.WriteArithmetic("lt")
		case "&gt;":
			cg.vmwriter.WriteArithmetic("gt")
		case "=":
			cg.vmwriter.WriteArithmetic("eq")
		case "~":
			// ビット反転演算子の処理
			cg.vmwriter.WriteArithmetic("not")
		default:
			fmt.Printf("未対応の演算子: %s\n", op)
		}
	}

}

// 項のコード生成処理
func (cg *CodeGenerator) generateTermCode(node compilationengine.ContainerNode, operations []string) []string {
	// 前処理：必要な状態変数の初期化
	methodName := ""
	varName := ""

	// 項の要素を順番に処理
	for i, child := range node.Children {
		switch n := child.(type) {
		case compilationengine.ParseNode:
			// パースノードの場合はタイプに応じた処理
			operations = cg.processParseNodeInTerm(n, i, node, methodName, varName, operations)
			// 識別子の場合は追加情報を収集
			if n.Type == jacktokenizer.IDENTIFIER {
				methodName, varName = cg.handleIdentifierInTerm(n, i, node, methodName, varName)
			}
		case compilationengine.ContainerNode:
			// コンテナノードの場合はコンテナタイプに応じた処理
			operations = cg.processContainerNodeInTerm(n, operations)
		}
	}

	return operations
}

// 項内のパースノード処理
func (cg *CodeGenerator) processParseNodeInTerm(node compilationengine.ParseNode, index int,
	parent compilationengine.ContainerNode,
	methodName string, varName string,
	operations []string) []string {
	switch node.Type {
	case jacktokenizer.INT_CONST:
		// 整数定数の処理
		cg.processIntConstant(node.Value)
	case jacktokenizer.STRING_CONST:
		// 文字列定数の処理
		cg.processStringConstant(node.Value)
	case jacktokenizer.KEYWORD:
		// キーワードの処理（true, false, this, null）
		cg.processKeywordConstant(node.Value)
	case jacktokenizer.SYMBOL:
		// 記号の処理（-, ~, [, ], など）
		return cg.processSymbolInTerm(node.Value, varName, operations)
	}

	return operations
}

// 項内のコンテナノード処理
func (cg *CodeGenerator) processContainerNodeInTerm(node compilationengine.ContainerNode, operations []string) []string {
	switch node.Name {
	case "expression":
		// 式の処理
		cg.generateExpressionCode(node)
	case "term":
		// 再帰的に項の処理
		return cg.generateTermCode(node, operations)
	}

	return operations
}

// 整数定数の処理
func (cg *CodeGenerator) processIntConstant(value string) {
	intValue, _ := strconv.Atoi(value)
	cg.vmwriter.WritePush("constant", intValue)

	// 負数フラグがセットされていれば、値を負数に変換
	if cg.negFlag {
		cg.vmwriter.WriteArithmetic("neg")
		cg.negFlag = false
	}
}

// 文字列定数の処理
func (cg *CodeGenerator) processStringConstant(value string) {
	// 文字列長をプッシュして新しい文字列オブジェクトを作成
	cg.vmwriter.WritePush("constant", len(value))
	cg.vmwriter.WriteCall("String.new", 1)

	// 各文字を追加
	for _, char := range value {
		cg.vmwriter.WritePush("constant", int(char))
		cg.vmwriter.WriteCall("String.appendChar", 2)
	}
}

// キーワード定数の処理
func (cg *CodeGenerator) processKeywordConstant(keyword string) {
	switch keyword {
	case "true":
		// true: 1をプッシュして-1に変換（Jack言語の仕様）
		cg.vmwriter.WritePush("constant", 1)
		cg.vmwriter.WriteArithmetic("neg")
	case "false", "null":
		// false/null: 0をプッシュ
		cg.vmwriter.WritePush("constant", 0)
	case "this":
		// this: 現在のオブジェクト参照（pointer 0）をプッシュ
		cg.vmwriter.WritePush("pointer", 0)
	}
}

// 識別子の処理と情報収集
func (cg *CodeGenerator) handleIdentifierInTerm(node compilationengine.ParseNode, index int,
	parent compilationengine.ContainerNode,
	currentMethodName string, currentVarName string) (string, string) {
	// 次のノードを確認
	if index+1 < len(parent.Children) {
		if nextNode, ok := parent.Children[index+1].(compilationengine.ParseNode); ok {
			if nextNode.Type == jacktokenizer.SYMBOL {
				switch nextNode.Value {
				case "(":
					// メソッド呼び出しの場合
					methodName := currentMethodName + node.Value
					cg.generateSubroutineCall(methodName, parent)
					return methodName, currentVarName
				case ".":
					// クラス名またはオブジェクト名
					return currentMethodName + node.Value + ".", currentVarName
				case "[":
					// 配列アクセスの場合
					return currentMethodName, node.Value
				}
			}
		}
	}

	// 単純な変数参照の場合
	kind := cg.getSegmentName(node.Value)
	idx := cg.getVarIndex(node.Value, cg.isMethodCall)
	cg.vmwriter.WritePush(kind, idx)

	return currentMethodName, currentVarName
}

// 記号の処理
func (cg *CodeGenerator) processSymbolInTerm(symbol string, varName string, operations []string) []string {
	switch symbol {
	case "-":
		// 単項マイナス演算子
		cg.negFlag = true
	case "~":
		// ビット反転演算子
		operations = append(operations, symbol)
	case "]":
		// 配列アクセス処理
		cg.processArrayAccess(varName)
	}

	return operations
}

// 配列アクセスの処理
func (cg *CodeGenerator) processArrayAccess(varName string) {
	// 変数の種類とインデックスを取得
	kind := cg.getSegmentName(varName)
	index := cg.symboltable.IndexOf(cg.subroutineName, varName)

	// 配列ベースアドレスをプッシュ
	cg.vmwriter.WritePush(kind, index)

	// アドレス計算：ベース + インデックス
	cg.vmwriter.WriteArithmetic("add")

	// 計算したアドレスをポインタにセット
	cg.vmwriter.WritePop("pointer", 1) // that = arr + i

	// 配列要素の値を取得
	cg.vmwriter.WritePush("that", 0) // push arr[i]
}

// 変数のインデックスを取得するヘルパー関数
func (cg *CodeGenerator) getVarIndex(name string, isMethod bool) int {
	index := cg.symboltable.IndexOf(cg.subroutineName, name)
	kind := cg.symboltable.KindOf(cg.subroutineName, name)
	if kind == "arg" && isMethod {
		return index + 1 // methodのときはthis分ずらす
	}
	return index
}

// 変数の種類を取得するヘルパー関数
func (cg *CodeGenerator) getSegmentName(name string) string {
	kind := cg.symboltable.KindOf(cg.subroutineName, name)
	switch kind {
	case "static":
		return "static"
	case "field":
		return "this"
	case "arg":
		return "argument"
	case "var":
		return "local"
	default:
		return ""
	}
}

// 一意のラベル名を生成
func (cg *CodeGenerator) generateUniqueLabel() string {
	label := fmt.Sprintf("%s_%d", cg.className, cg.labelCount)
	cg.labelCount++
	return label
}

func (cg *CodeGenerator) generateSubroutineCall(methodName string, node compilationengine.ContainerNode) {
	// 引数リストの処理
	argCount := 0
	for _, child := range node.Children {
		if n, ok := child.(compilationengine.ContainerNode); ok && n.Name == "expressionList" {
			argCount = cg.generateExpressionListCode(n)
		}
	}

	subroutineClassName := strings.Split(methodName, ".")[0]
	subroutineClassName = strings.ToUpper(string(subroutineClassName[0])) + subroutineClassName[1:] // クラス名の最初の文字を大文字に変換
	callSubroutineName := methodName[strings.LastIndex(methodName, ".")+1:]
	if argCount == 0 && cg.subroutineKindMap[subroutineClassName][callSubroutineName] == "method" {
		// 引数がない場合は1を設定
		argCount = 1
	}

	// メソッド呼び出しのコードを生成
	objectName := strings.Split(methodName, ".")[0]
	methodName = strings.ToUpper(string(methodName[0])) + methodName[1:] // メソッド名の最初の文字を大文字に変換
	kind := cg.getSegmentName(objectName)
	if kind != "" {
		index := cg.symboltable.IndexOf(cg.subroutineName, objectName)

		cg.vmwriter.WritePush(kind, index)
	}
	cg.vmwriter.WriteCall(methodName, argCount)
}

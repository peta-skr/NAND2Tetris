package vmwriter

import "strconv"

// SEGMENT
const (
	CONST   = "constant"
	ARG     = "argument"
	LOCAL   = "local"
	STATIC  = "static"
	THIS    = "this"
	THAT    = "that"
	POINTER = "pointer"
	TEMP    = "temp"
)

// COMMAND
const (
	ADD = "add"
	SUB = "sub"
	EQ  = "eq"
	LT  = "lt"
	GT  = "gt"
	NEG = "neg"
	AND = "and"
	OR  = "or"
	NOT = "not"
)

type VMWriter struct {
	Content []string // VM命令を格納するスライス
}

func Constructor() VMWriter {
	vmWriter := VMWriter{}

	return vmWriter
}

func (v *VMWriter) WritePush(segment string, index int) {
	// VMWriterにpush命令を書き込む
	// segment: セグメント名
	// index: インデックス
	v.Content = append(v.Content, "push "+segment+" "+strconv.Itoa(index))
}

func (v *VMWriter) WritePop(segment string, index int) {
	// VMWriterにpop命令を書き込む
	// segment: セグメント名
	// index: インデックス
	v.Content = append(v.Content, "pop "+segment+" "+strconv.Itoa(index))
}

func (v *VMWriter) WriteArithmetic(command string) {
	// VMWriterに演算命令を書き込む
	// command: 演算命令
	v.Content = append(v.Content, command)
}

func (v *VMWriter) WriteLabel(label string) {
	// VMWriterにラベル命令を書き込む
	// label: ラベル名
	v.Content = append(v.Content, "label "+label)
}

func (v *VMWriter) WriteGoto(label string) {
	// VMWriterにgoto命令を書き込む
	// label: ラベル名
	v.Content = append(v.Content, "goto "+label)
}

func (v *VMWriter) WriteIf(label string) {
	// VMWriterにif命令を書き込む
	// label: ラベル名
	v.Content = append(v.Content, "if-goto "+label)
}

func (v *VMWriter) WriteCall(name string, nArgs int) {
	// VMWriterにcall命令を書き込む
	// name: 関数名
	// nArgs: 引数の数
	v.Content = append(v.Content, "call "+name+" "+strconv.Itoa(nArgs))
}

func (v *VMWriter) WriteFunction(name string, nLocals int) {
	// VMWriterにfunction命令を書き込む
	// name: 関数名
	// nLocals: ローカル変数の数
	v.Content = append(v.Content, "function "+name+" "+strconv.Itoa(nLocals))
}

func (v *VMWriter) WriteReturn() {
	// VMWriterにreturn命令を書き込む
	v.Content = append(v.Content, "return")
}

func (v *VMWriter) Close() {
	// VMWriterを閉じる
	// ファイルを閉じる
}

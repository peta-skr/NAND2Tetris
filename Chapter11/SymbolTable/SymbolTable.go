package symboltable

import (
	"fmt"
	"strconv"
)

type Table map[string]map[string]string

type ScopeName int

const (
	CLASS_SCOPE ScopeName = iota
	SUBROUTINE_SCOPE
)

type SymbolTable struct {
	// クラススコープのシンボルテーブル
	ClassSymbolTable      Table
	SubroutineSymbolTable Table
	// サブルーチンスコープのシンボルテーブル
	SubroutineSymbolTables map[string]Table
	CurrentScope           ScopeName
}

// var ClassSymbolTable map[string]map[string]string
// var SubroutineSymbolTable map[string]map[string]string

const (
	STATIC = "static"
	ARG    = "arg"
	FIELD  = "field"
	VAR    = "var"
	NONE   = "none"
)

var class_staticCount int
var class_fieldCount int
var class_argCount int
var class_varCount int
var class_localCount int

var subroutine_staticCount int
var subroutine_fieldCount int
var subroutine_argCount int
var subroutine_varCount int
var subroutine_localCount int

func Constructor() SymbolTable {
	SymbolTable := SymbolTable{
		ClassSymbolTable:       make(map[string]map[string]string),
		SubroutineSymbolTables: make(map[string]Table),
		CurrentScope:           CLASS_SCOPE,
	}
	class_staticCount = 0
	class_fieldCount = 0
	class_argCount = 0
	class_varCount = 0
	class_localCount = 0

	return SymbolTable
}

func (s *SymbolTable) StartSubroutine(subroutineName string) {
	s.SubroutineSymbolTables[subroutineName] = make(map[string]map[string]string)
	s.SubroutineSymbolTable = s.SubroutineSymbolTables[subroutineName]
	s.CurrentScope = SUBROUTINE_SCOPE
	subroutine_staticCount = 0
	subroutine_fieldCount = 0
	subroutine_argCount = 0
	subroutine_varCount = 0
	subroutine_localCount = 0
}

func (s *SymbolTable) EndSubroutine() {
	s.CurrentScope = CLASS_SCOPE
}

func (s *SymbolTable) Define(name string, varType string, kind string) {
	// 変数を定義する
	if s.CurrentScope == CLASS_SCOPE {
		// クラススコープの場合
		if _, ok := s.ClassSymbolTable[name]; !ok {
			s.ClassSymbolTable[name] = make(map[string]string)
		}
		s.ClassSymbolTable[name]["type"] = varType
		s.ClassSymbolTable[name]["kind"] = kind
	} else {
		// サブルーチンスコープの場合
		if _, ok := s.SubroutineSymbolTable[name]; !ok {
			s.SubroutineSymbolTable[name] = make(map[string]string)
		}
		s.SubroutineSymbolTable[name]["type"] = varType
		s.SubroutineSymbolTable[name]["kind"] = kind
	}

	switch kind {
	case STATIC:
		if s.CurrentScope == CLASS_SCOPE {
			s.ClassSymbolTable[name]["index"] = strconv.Itoa(class_staticCount)
			class_staticCount++
		} else {
			s.SubroutineSymbolTable[name]["index"] = strconv.Itoa(subroutine_staticCount)
			subroutine_staticCount++
		}
	case ARG:
		if s.CurrentScope == CLASS_SCOPE {
			s.ClassSymbolTable[name]["index"] = strconv.Itoa(class_argCount)
			class_argCount++
		} else {
			s.SubroutineSymbolTable[name]["index"] = strconv.Itoa(subroutine_argCount)
			subroutine_argCount++
		}
	case FIELD:
		if s.CurrentScope == CLASS_SCOPE {
			s.ClassSymbolTable[name]["index"] = strconv.Itoa(class_fieldCount)
			class_fieldCount++
		} else {
			s.SubroutineSymbolTable[name]["index"] = strconv.Itoa(subroutine_fieldCount)
			subroutine_fieldCount++
		}
	case VAR:
		if s.CurrentScope == CLASS_SCOPE {
			s.ClassSymbolTable[name]["index"] = strconv.Itoa(class_varCount)
			class_varCount++
		} else {
			s.SubroutineSymbolTable[name]["index"] = strconv.Itoa(subroutine_varCount)
			subroutine_varCount++
		}
	default:
		fmt.Println("エラー: 定義できない変数の種類です")
		return // エラー処理
	}

}

func (s *SymbolTable) VarCount(tableName string, kind string) int {
	if s.CurrentScope == CLASS_SCOPE {
		switch kind {
		case STATIC:
			return class_staticCount
		case FIELD:
			return class_fieldCount
		default:
			return -1 // エラー処理
		}
	} else {
		switch kind {
		case ARG:
			table := s.SubroutineSymbolTables[tableName]
			// count := 0
			// for _, i := range table {
			// 	fmt.Println("this: ", i)
			// 	if i["kind"] == ARG {
			// 		count++
			// 	}
			// }
			return len(table)
		case VAR:
			table := s.SubroutineSymbolTables[tableName]
			count := 0
			for _, i := range table {
				fmt.Println("this: ", i)
				if i["kind"] == VAR {
					count++
				}
			}
			return count
		default:
			return -1 // エラー処理
		}
	}
}

func (s *SymbolTable) KindOf(scopeName string, name string) string {

	if s.CurrentScope == CLASS_SCOPE {
		if _, ok := s.ClassSymbolTable[name]; ok {
			return s.ClassSymbolTable[name]["kind"]
		} else {
			return NONE
		}
	} else {
		if table, ok := s.SubroutineSymbolTables[scopeName]; ok {
			return table[name]["kind"]
		} else {
			return NONE
		}
	}
}

func (s *SymbolTable) TypeOf(scopeName string, name string) string {
	if s.CurrentScope == CLASS_SCOPE {
		if _, ok := s.ClassSymbolTable[name]; ok {
			return s.ClassSymbolTable[name]["type"]
		} else {
			return NONE
		}
	} else {
		if table, ok := s.SubroutineSymbolTables[scopeName]; ok {
			return table[name]["type"]
		} else {
			return NONE
		}
	}
}

func (s *SymbolTable) IndexOf(scopeName string, name string) int {
	if s.CurrentScope == CLASS_SCOPE {
		if _, ok := s.ClassSymbolTable[name]; ok {
			index, _ := strconv.Atoi(s.ClassSymbolTable[name]["index"])
			return index
		} else {
			return -1
		}
	} else {
		if table, ok := s.SubroutineSymbolTables[scopeName]; ok {
			index, _ := strconv.Atoi(table[name]["index"])
			return index
		} else {
			return -1
		}
	}
}

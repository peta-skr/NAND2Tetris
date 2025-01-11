package symbolTable

type SymbolTable map[string]int

func Initialize() SymbolTable {
	symbolTable := make(SymbolTable)

	return symbolTable
}

func (s *SymbolTable) AddEntry(symbol string, address int) {
	(*s)[symbol] = address
}

func (s *SymbolTable) contains(symbol string) bool {
	_, ok := (*s)[symbol]

	return ok
}

func (s *SymbolTable) getAddress(symbol string) int {
	return (*s)[symbol]
}
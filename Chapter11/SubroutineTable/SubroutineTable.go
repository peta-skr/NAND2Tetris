package subroutinetable

type SubroutineTable map[string]string

const (
	METHOD      = "method"
	FUNCTION    = "function"
	CONSTRUCTOR = "constructor"
)

func Constructor() SubroutineTable {
	subroutineTable := make(map[string]string)
	return subroutineTable
}

func (s *SubroutineTable) Define(subroutineName string, subroutineType string) {
	(*s)[subroutineName] = subroutineType

}

package tree_paser

import (
	"github.com/TobiasYin/go-lsp/logs"
	"github.com/TobiasYin/go-lsp/lsp/defines"
)

const (
	global = iota
	pakeji
	unda
)

var TOP_SCOPE = NewScope()
var KEYWORDS = []string{"unda", "pakeji", "tumia", "kama",
	"au", "sivyo", "andika", "wakati", "sukuma", "kwa", "ktk",
	"jaza"}

// Beggining and end of scope
type Range struct {
	// Starting index
	Start uint
	// Ending index
	End uint
}

type Variable struct {
	Name string
	Type uint
	Definition string
}

// Hold scope between a range with their variables
type Scope struct {
	File           string
	Type           uint
	RangeOfScope   Range
	VariableScopes []Variable
	ChildScope     []Scope
}

func NewScope() *Scope {
	return &Scope{
		Type: global,
	}
}

func (*Scope) parse(content *defines.TextDocumentContentChangeEvent){
	if content.Range.Start.Line == content.Range.End.Line {
		//meaning we recieved the whole document, so parse it from
		//top to bottom
	}else{
		logs.Println("RANGE NOT 0, didn't get full file")
	}
}

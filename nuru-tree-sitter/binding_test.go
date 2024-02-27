package nuru_tree_sitter

import (
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
)

func TestSimpleProgramIsParsed(t *testing.T) {
	parser := sitter.NewParser()
	parser.SetLanguage(GetLanguage())
}

package nuru_tree_sitter

// #include "src/parser.c"
// const TSLanguage *tree_sitter_nuru();
import "C"
import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	//println("NURU")
	ptr := unsafe.Pointer(C.tree_sitter_nuru())
	return sitter.NewLanguage(ptr)
}

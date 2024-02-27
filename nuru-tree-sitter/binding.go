package nuru_tree_sitter

// #cgo CFLAGS: -Isrc/ -Isrc/tree_sitter/
// #include <parser.h>
//const TSLanguage * tree_sitter_nuru()
import "C"
import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	ptr := unsafe.Pointer(C.tree_sitter_nuru())
	return sitter.NewLanguage(ptr)
}

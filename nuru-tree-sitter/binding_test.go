package nuru_tree_sitter

import (
	"context"
	"fmt"
	"os"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
)

func TestSimpleProgramIsParsed(t *testing.T) {
	parser := sitter.NewParser()
	parser.SetLanguage(GetLanguage())
	file, err := os.ReadFile("../a.nr")
	if err != nil {
		t.Fatalf("Gotten error opening file: %s", err)
	}

	tree, err := parser.ParseCtx(context.Background(), nil, file)
	if err != nil {
		t.Fatalf("Error reading tree: %s", err)
	}

	children := tree.RootNode().NamedChildCount()

	fmt.Printf("Named children are: %d", children)
	//if children == 4 {
	t.Fatalf("Named children can't be nil %d", children)
	//}
}

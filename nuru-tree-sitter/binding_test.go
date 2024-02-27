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

	firstChildFunction := tree.RootNode().ChildByFieldName("functionname")
	if firstChildFunction == nil {
		t.Fatal("First function not found")
	}

	starts := firstChildFunction.StartByte()
	ends := firstChildFunction.EndByte()

	fmt.Printf("FIrst function variable Name is: %s", file[starts:ends])
}

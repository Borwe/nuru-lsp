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

	q, err := sitter.NewQuery([]byte("(declaration_statement)"), GetLanguage())
	query := sitter.NewQueryCursor()
	query.Exec(q, tree.RootNode())

	m, ok := query.NextMatch()
	if !ok {
		fmt.Println("No match found")
	}

	for _, c := range m.Captures {
		fmt.Printf("function variable Name is: %s", c.Node.Content(file))
	}
}

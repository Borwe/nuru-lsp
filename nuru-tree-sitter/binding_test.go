package nuru_tree_sitter

import (
	"context"
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

	//get pakeji statements should be 1
	q, err := sitter.NewQuery([]byte("(pakeji_tumia_statement) @pakeji"), GetLanguage())
	if err != nil {
		t.Fatalf("Failed to create query with error %s", err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, tree.RootNode())

	matches := 0
	for {

		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			matches += 1
			t.Log(c.Node.Content(file))
		}
	}

	if matches != 1 {
		t.Fatalf("Failed to get 1 match of pakeji, got %d", matches)
	}
}

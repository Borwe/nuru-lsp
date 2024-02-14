package tests

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"nuru-lsp/data"
	"os"
	"testing"
)

func TestBuildingTreePakeji(t *testing.T) {
	file, err := os.Open("full_pakeji.nr")
	if err != nil {
		t.Fatal("Error opening:", err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	full_pakeji := &data.Top{
		Pakeji: data.Pakeji{
			Items: []data.FuncVar{{
				Items: []data.FuncVar{{
					Line:      4,
					Name:      "@.zolovar",
					StartDecl: 9,
					EndDecl:   26,
					IsScope:   false,
				},
				},
				Line:      2,
				Name:      "andaa",
				StartDecl: 5,
				EndDecl:   19,
				IsScope:   true,
			}, {
				Items: []data.FuncVar{{
					Line:      8,
					Name:      "mti",
					StartDecl: 9,
					EndDecl:   21,
					IsScope:   false,
				}},
				Line:      7,
				Name:      "chora",
				StartDecl: 5,
				EndDecl:   19,
				IsScope:   true,
			},
			},
		},
		Items: []data.FuncVar{},
	}

	parsedTree := data.ParseTree(lines)

	marshalled_full_pakeji_tree, err := json.Marshal(full_pakeji.Tree)
	marshalled_parsedTree, err := json.Marshal(parsedTree)

	fmt.Println("full_pakeji:", string(marshalled_full_pakeji_tree))
	fmt.Println("full_pakeji:", string(marshalled_parsedTree))

	if !bytes.Equal(marshalled_full_pakeji_tree, marshalled_parsedTree) {
		t.Fatal("Parsed fulle_pakeji_tree not matching as expected")
	}
}

func TestBuildingTree2(t *testing.T) {
	t.Fatalf("NOT IMPLEMENTED YET")
}

package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"nuru-lsp/data"
	"os"
	"testing"
)

func TestBuildingTree1(t *testing.T) {
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

	full_pakeji := data.Data{
		File:    "full_pakeji",
		Version: 0,
		Content: lines,
		Tree: &data.Top{
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
		},
	}

	marshalled_full_pakeji_tree, err := json.Marshal(full_pakeji.Tree)

	fmt.Println("full_pakeji:", string(marshalled_full_pakeji_tree))

	t.Fatalf("NOT IMPLEMENTED YET")
}

func TestBuildingTree2(t *testing.T) {
	t.Fatalf("NOT IMPLEMENTED YET")
}

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
	if err!= nil {
		t.Fatal("Error opening:",err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	full_pakeji := data.Data{
		File: "full_pakeji",
		Version: 0,
		Content: lines,
	}

	marshalled_full_pakeji, err := json.Marshal(full_pakeji)

	fmt.Println("full_pakeji:", string(marshalled_full_pakeji))

	t.Fatalf("NOT IMPLEMENTED YET")
}

func TestBuildingTree2(t *testing.T) {
	t.Fatalf("NOT IMPLEMENTED YET")
}

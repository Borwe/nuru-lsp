package main

import (
	"context"
	"log"
	"nuru-lsp/completions"
	"nuru-lsp/data"
	"os"

	"github.com/TobiasYin/go-lsp/logs"
	"github.com/TobiasYin/go-lsp/lsp"
	"github.com/TobiasYin/go-lsp/lsp/defines"
)

func setupLog() {
	foundFile := true
	if len(os.Args) == 2 {
		file := os.Args[1]
		f, err := os.Open(file)
		if err != nil {
			f, err = os.Create(file)
			if err != nil {
				foundFile = false
			}
		}
		logs.Init(log.New(f, "nuru-lsp:=> ", 0))
	} else {
		foundFile = false
	}
	if foundFile {
		return
	}
	logs.Init(log.New(os.Stderr, "nuru-lsp:=> ", 0))
}

func main() {
	setupLog()

	server := lsp.NewServer(&lsp.Options{
		CompletionProvider: &defines.CompletionOptions{
			TriggerCharacters: &[]string{"."},
		},
	})
	server.OnHover(func(ctx context.Context,
		req *defines.HoverParams) (*defines.Hover, error) {
		logs.Println("Hover:", req)
		return &defines.Hover{
			Contents: defines.MarkupContent{
				Kind:  defines.MarkupKindPlainText,
				Value: "OnHover Testing",
			},
		}, nil
	})

	server.OnCompletion(completions.CompletionFunc)

	server.OnDidChangeTextDocument(data.OnDataChange)

	server.Run()
}

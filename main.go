package main

import (
	"context"
	"log"
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

	server.OnCompletion(func(ctx context.Context,
		req *defines.CompletionParams) (*[]defines.CompletionItem, error) {
		logs.Println("Completion:", req)
		return nil, nil
	})

	server.OnDidChangeTextDocument(func(ctx context.Context, req *defines.DidChangeTextDocumentParams) error {
		logs.Println("DocChange: ", req)
		logs.Println("Version: ", req.TextDocument.Version)
		logs.Println("URI: ", req.TextDocument.Uri)
		for i, v := range req.ContentChanges {
			logs.Println("Range ", i, ": ", v.Range)
			logs.Println("ContentChange ", i, ": ", v.Text)
		}
		return nil
	})

	server.Run()
}

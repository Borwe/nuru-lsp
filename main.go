package main

import (
	"context"
	"log"
	"nuru-lsp/completions"
	"nuru-lsp/data"
	"nuru-lsp/server"
	"os"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
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
		logs.Init(log.New(f, "", 0))
	} else {
		foundFile = false
	}
	if foundFile {
		return
	}
	logs.Init(log.New(os.Stderr, "", 0))
}

func main() {
	setupLog()

	server.Server.OnHover(func(ctx context.Context,
		req *defines.HoverParams) (*defines.Hover, error) {
		logs.Println("Hover:", req)
		return &defines.Hover{
			Contents: defines.MarkupContent{
				Kind:  defines.MarkupKindPlainText,
				Value: "OnHover Testing",
			},
		}, nil
	})

	server.Server.OnCompletion(completions.CompletionFunc)

	server.Server.OnDidChangeTextDocument(data.OnDataChange)

	server.Server.Run()
}

package main

import (
	"context"
	"fmt"
	"nuru-lsp/completions"
	"nuru-lsp/data"
	"nuru-lsp/server"
	"nuru-lsp/setup"
	"os"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
)

const Version = "0.0.06"

func main() {

	if len(os.Args) >= 2 {
		if os.Args[1] == "--version" {
			fmt.Println("VERSION:",Version)
			return
		}
	}

	setup.SetupLog()

	server.Server.OnInitialized(func(ctx context.Context, req *defines.InitializeParams) (err error) {
		return nil
	})

	server.Server.OnHover(func(ctx context.Context,
		req *defines.HoverParams) (*defines.Hover, error) {
		logs.Println("Hover:", req)
		return &defines.Hover{
			Contents: defines.MarkupContent{
				Kind:  defines.MarkupKindPlainText,
				Value: "OnHover Not implemented yet",
			},
		}, nil
	})

	server.Server.OnDidSaveTextDocument(func(ctx context.Context, req *defines.DidSaveTextDocumentParams) (err error) {
		return nil
	})
	server.Server.OnDidOpenTextDocument(data.OnDocOpen)
	server.Server.OnDidCloseTextDocument(data.OnDidClose)
	server.Server.OnDidChangeTextDocument(data.OnDataChange)
	server.Server.OnCompletion(completions.CompletionFunc)

	server.Server.Run()
}

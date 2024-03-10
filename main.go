package main

import (
	"context"
	"nuru-lsp/completions"
	"nuru-lsp/data"
	"nuru-lsp/server"
	"nuru-lsp/setup"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
)

func main() {
	
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
				Value: "OnHover Testing",
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

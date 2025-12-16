package main

import (
	"context"
	"fmt"
	"nuru-lsp/completions"
	"nuru-lsp/data"
	"nuru-lsp/hover"
	"nuru-lsp/definition"
	"nuru-lsp/server"
	"nuru-lsp/setup"
	"os"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
)

const Version = "0.0.10"

func main() {

	if len(os.Args) >= 2 {
		if os.Args[1] == "--version" {
			fmt.Printf("VERSION: v%s\n", Version)
			return
		}
	}

	setup.SetupLog()

	server.Server.OnInitialized(func(ctx context.Context, req *defines.InitializeParams) (err error) {
		return nil
	})

	server.Server.OnHover(hover.HoverFunc)

	server.Server.OnDefinition(definition.DefinitionFunc)

	server.Server.OnDidSaveTextDocument(func(ctx context.Context, req *defines.DidSaveTextDocumentParams) (err error) {
		return nil
	})
	server.Server.OnDidChangeWatchedFiles(func(ctx context.Context,
		req *defines.DidChangeWatchedFilesParams) (err error) {
		return nil
	})
	server.Server.OnDidOpenTextDocument(data.OnDocOpen)
	server.Server.OnDidCloseTextDocument(data.OnDidClose)
	server.Server.OnDidChangeTextDocument(data.OnDataChange)
	server.Server.OnCompletion(completions.CompletionFunc)
	server.Server.OnExit(func(ctx context.Context, req *interface{}) (err error) {
		logs.Println("EXIT VARIABLE:",req)
		return nil
	})
	server.Server.OnShutdown(func(ctx context.Context, req *interface{}) (err error) {
		logs.Println("SHUTDOWN VARIABLE:",req)
		return nil
	})

	server.Server.Run()
}

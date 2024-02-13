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

func getLogFile() (bool,*string) {
		if os.Args[1]== "--stdio" {
			if len(os.Args) == 2 {
				return false, nil
			}
			return true, &os.Args[2]
		}else{
			return true, &os.Args[1]
		}
}

func setupLog() {
	foundFile, file := getLogFile()
	if foundFile {
		f, err := os.Open(*file)
		if err != nil {
			f, err = os.Create(*file)
			if err != nil {
				foundFile = false
			}
		}
		logs.Init(log.New(f, "", 0))
		return
	}
	logs.Init(log.New(os.Stderr, "", 0))
}

func main() {
	setupLog()

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

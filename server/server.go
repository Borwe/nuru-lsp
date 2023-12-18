package server

import (
	"github.com/Borwe/go-lsp/lsp"
	"github.com/Borwe/go-lsp/lsp/defines"
)

var Server = lsp.NewServer(&lsp.Options{
	CompletionProvider: &defines.CompletionOptions{
		TriggerCharacters: &[]string{"."},
	},
})

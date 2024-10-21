package hovers

import (
	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
)

func GetHover(req *defines.HoverParams) (*defines.Hover, error) {
	logs.Println("Hover:", req)
	return &defines.Hover{
		Contents: defines.MarkupContent{
			Kind:  defines.MarkupKindPlainText,
			Value: "OnHover Not implemented yet",
		},
	}, nil
}

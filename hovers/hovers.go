package hovers

import (
	"fmt"
	"nuru-lsp/data"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/NuruProgramming/Nuru/module"
	"github.com/NuruProgramming/Nuru/object"
)

func getAllIdengitiers() {}

var StdTumiasInfo = map[string]string{}

func init() {
	StdTumiasInfo["os"] = "Pakeji Yenye Operashani wa System"
	StdTumiasInfo["muda"] = "Pakeji yenye unaweza tumia kuapata muda"
	StdTumiasInfo["mtandao"] = "Pakeji yenye inaweza tumika kwa kufanya networking na mengine"
	StdTumiasInfo["jsoni"] = "Pakeji wa kufungua na kufunga jsoni"
	StdTumiasInfo["hisabati"] = "Pakeji yenye inasaidia kufanya hesabu na hisabati mingi"
}

func getPackageAndMethodFromWord(word *string) (pkg *string,
	parameter *string, isFunction bool, found bool) {
	pkg = nil
	parameter = nil
	isFunction = true
	found = false
	return
}

func GetHover(req *defines.HoverParams) (*defines.Hover, error) {
	d, ok := data.Pages[string(req.TextDocument.Uri)]
	if !ok {
		return nil, nil
	}

	wordHovered, ok := d.GetWord(req.Position)
	if !ok {
		return nil, nil
	}

	logs.Println("WORD:", *wordHovered)

	stdTumia, ok := StdTumiasInfo[*wordHovered]
	if ok {
		//meaning this is an stdtumia being hovered
		return &defines.Hover{
			Contents: defines.MarkupContent{
				Kind:  defines.MarkupKindPlainText,
				Value: stdTumia,
			},
		}, nil
	}

	//slip word with . and remove any ()
	pkg, property, isFunction, found := getPackageAndMethodFromWord(wordHovered)
	if mod, ok := module.Mapper[*pkg]; ok && found {
		//only hisabati has constants,
		if *pkg == "hisabati" {
			obj, ok := module.Constants[*property]
			if floatValue, floatValid := obj.(*object.Float); ok && floatValid && !isFunction {
				return &defines.Hover{
					Contents: defines.MarkupContent{
						Kind:  defines.MarkupKindPlainText,
						Value: fmt.Sprintf("Kutoka pakeji wa %s kuna value wa %d", *pkg, floatValue.Value),
					},
				}, nil
			}
		} else {

		}
	}

	return &defines.Hover{
		Contents: defines.MarkupContent{
			Kind:  defines.MarkupKindPlainText,
			Value: "OnHover Not implemented yet",
		},
	}, nil
}

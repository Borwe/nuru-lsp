package hover

import (
	"context"
	"fmt"
	"nuru-lsp/completions"
	"nuru-lsp/data"
	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/NuruProgramming/Nuru/ast"
)

// maelezo ya function za nuru
var KeywordDocs = map[string]string{
	"fanya":   "Kutangaza kigezo kipya\nMfano: fanya x = 5",
	"unda":    "Kuunda function mpya\nMfano: unda(a, b) { rudisha a + b }",
	"kama":    "Kauli ya masharti (if statement)\nMfano: kama (x > 5) { andika(x) }",
	"sivyo":   "Sehemu ya sivyo ya kauli ya masharti (else)\nMfano: kama (x > 5) { ... } sivyo { ... }",
	"au":      "Masharti mengine (else if)\nMfano: kama (x > 5) { ... } au (x > 3) { ... }",
	"wakati":  "Kitanzi cha wakati (while loop)\nMfano: wakati (x < 10) { x = x + 1 }",
	"kwa":     "Kitanzi cha kwa (for loop)\nMfano: kwa i ktk orodha { andika(i) }",
	"ktk":     "Ndani ya (in) - kutumika na kitanzi\nMfano: kwa i ktk [1,2,3] { ... }",
	"rudisha": "Kurudisha thamani kutoka kwa function\nMfano: rudisha x + y",
	"vunja":   "Kuvunja kitanzi (break)\nMfano: vunja",
	"endelea": "Kuendelea na iteration inayofuata (continue)\nMfano: endelea",
	"kweli":   "Thamani ya kweli (true)",
	"sikweli": "Thamani ya sikweli (false)",
	"tupu":    "Thamani tupu (null)",
	"badili":  "Kauli ya kubadili (switch statement)\nMfano: badili x { ikiwa 1 { ... } }",
	"ikiwa":   "Kesi katika kauli ya badili (case)\nMfano: ikiwa 1 { andika(\"moja\") }",
	"kawaida": "Kesi ya kawaida katika badili (default)\nMfano: kawaida { andika(\"nyingine\") }",
	"tumia":   "Kuingiza pakeji/moduli\nMfano: tumia hisabati",
	"pakeji":  "Kutangaza pakeji\nMfano: pakeji jina { ... }",
	"@":       "Rejeleo la sasa (this/self reference)",
}


func getWordAtPosition(content []string, line, char uint) string {
	if int(line) >= len(content) {
		return ""
	}
	
	lineContent := content[line]
	if int(char) > len(lineContent) {
		return ""
	}

	// tafuta mwanzo na mwisho wa neno
	start := int(char)
	end := int(char)

	// kama uko mwisho, rudi nyuma kutafuta mwanzo wa neno
	for start > 0 && isWordChar(rune(lineContent[start-1])) {
		start--
	}

	// kama uko mwanzo, nenda mbele
	for end < len(lineContent) && isWordChar(rune(lineContent[end])) {
		end++
	}

	if start == end {
		return ""
	}

	return lineContent[start:end]
}

func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || 
		   (r >= 'A' && r <= 'Z') || 
		   (r >= '0' && r <= '9') || 
		   r == '_' || 
		   r == '@'
}

func findDefinition(node ast.Node, name string) string {
	if node == nil {
		return ""
	}

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			if result := findDefinition(stmt, name); result != "" {
				return result
			}
		}
	case *ast.LetStatement:
		if n.Name != nil && n.Name.Value == name {
			if n.Value != nil {
				return fmt.Sprintf("fanya %s = %s", name, n.Value.String())
			}
			return fmt.Sprintf("fanya %s", name)
		}
	case *ast.Assign:
		if n.Name != nil && n.Name.String() == name {
			if n.Value != nil {
				return fmt.Sprintf("%s = %s", name, n.Value.String())
			}
		}
	case *ast.FunctionLiteral:
		if n.Body != nil {
			return findDefinition(n.Body, name)
		}
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			if result := findDefinition(stmt, name); result != "" {
				return result
			}
		}
	case *ast.ExpressionStatement:
		return findDefinition(n.Expression, name)
	}

	return ""
}

func HoverFunc(ctx context.Context, req *defines.HoverParams) (*defines.Hover, error) {
	file := string(req.TextDocument.Uri)

	data.PagesMutext.Lock()
	defer data.PagesMutext.Unlock()

	doc, found := data.Pages[file]
	if !found {
		logs.Println("Hover: document not found:", file)
		return nil, nil
	}

	word := getWordAtPosition(doc.Content, req.Position.Line, req.Position.Character)
	if word == "" {
		return nil, nil
	}

	logs.Println("Hover word:", word)

	var hoverText string

	if desc, ok := KeywordDocs[word]; ok {
		hoverText = desc
	}

	if hoverText == "" {
		if desc, ok := completions.Functions[word]; ok {
			hoverText = desc
		}
	}

	if hoverText == "" && doc.RootTree != nil {
		if def := findDefinition(*doc.RootTree, word); def != "" {
			hoverText = def
		}
	}

	if hoverText == "" {
		for _, tumia := range data.TUMIAS {
			if tumia == word {
				hoverText = fmt.Sprintf("Pakeji: %s\nTumia: tumia %s", word, word)
				break
			}
		}
	}

	if hoverText == "" {
		return nil, nil
	}

	return &defines.Hover{
		Contents: defines.MarkupContent{
			Kind:  defines.MarkupKindMarkdown,
			Value: hoverText,
		},
	}, nil
}

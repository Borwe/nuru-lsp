package definition

import (
	"context"
	"nuru-lsp/data"
	"strings"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/NuruProgramming/Nuru/ast"
)

type DefLocation struct {
	Line   uint
	Column uint
	Found  bool
}

func getWordAtPosition(content []string, line, char uint) string {
	if int(line) >= len(content) {
		return ""
	}

	lineContent := content[line]
	if int(char) > len(lineContent) {
		return ""
	}

	start := int(char)
	end := int(char)

	for start > 0 && isWordChar(rune(lineContent[start-1])) {
		start--
	}

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
		r == '_'
}

func findDefinitionLocation(node ast.Node, name string) DefLocation {
	if node == nil {
		return DefLocation{Found: false}
	}

	switch n := node.(type) {
	case *ast.Program:
		for _, stmt := range n.Statements {
			if loc := findDefinitionLocation(stmt, name); loc.Found {
				return loc
			}
		}
	case *ast.LetStatement:
		if n.Name != nil && n.Name.Value == name {
			return DefLocation{
				Line:   uint(n.Token.Line - 1), 
				Column: 0,
				Found:  true,
			}
		}
		if n.Value != nil {
			if loc := findDefinitionLocation(n.Value, name); loc.Found {
				return loc
			}
		}
	case *ast.Assign:
		if n.Name != nil && n.Name.String() == name {
			return DefLocation{
				Line:   uint(n.Token.Line - 1),
				Column: 0,
				Found:  true,
			}
		}
	case *ast.FunctionLiteral:
		for _, param := range n.Parameters {
			if param.Value == name {
				return DefLocation{
					Line:   uint(param.Token.Line - 1),
					Column: 0,
					Found:  true,
				}
			}
		}
		if n.Body != nil {
			if loc := findDefinitionLocation(n.Body, name); loc.Found {
				return loc
			}
		}
	case *ast.BlockStatement:
		for _, stmt := range n.Statements {
			if loc := findDefinitionLocation(stmt, name); loc.Found {
				return loc
			}
		}
	case *ast.ExpressionStatement:
		return findDefinitionLocation(n.Expression, name)
	case *ast.Package:
		if n.Name != nil && n.Name.Value == name {
			return DefLocation{
				Line:   uint(n.Token.Line - 1),
				Column: 0,
				Found:  true,
			}
		}
		if n.Block != nil {
			if loc := findDefinitionLocation(n.Block, name); loc.Found {
				return loc
			}
		}
	}

	return DefLocation{Found: false}
}

func DefinitionFunc(ctx context.Context, req *defines.DefinitionParams) (*[]defines.LocationLink, error) {
	file := string(req.TextDocument.Uri)

	data.PagesMutext.Lock()
	defer data.PagesMutext.Unlock()

	doc, found := data.Pages[file]
	if !found {
		logs.Println("Definition: document not found:", file)
		return nil, nil
	}

	word := getWordAtPosition(doc.Content, req.Position.Line, req.Position.Character)
	if word == "" {
		return nil, nil
	}

	logs.Println("Definition word:", word)

	if doc.RootTree != nil {
		loc := findDefinitionLocation(*doc.RootTree, word)
		if loc.Found {
			uri := file
			if !strings.HasPrefix(uri, "file://") {
				uri = "file://" + uri
			}

			return &[]defines.LocationLink{
				{
					TargetUri: defines.DocumentUri(uri),
					TargetRange: defines.Range{
						Start: defines.Position{Line: loc.Line, Character: loc.Column},
						End:   defines.Position{Line: loc.Line, Character: loc.Column + uint(len(word))},
					},
					TargetSelectionRange: defines.Range{
						Start: defines.Position{Line: loc.Line, Character: loc.Column},
						End:   defines.Position{Line: loc.Line, Character: loc.Column + uint(len(word))},
					},
				},
			}, nil
		}
	}

	return nil, nil
}

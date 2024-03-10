package data

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	nuru_tree_sitter "nuru-lsp/nuru-tree-sitter"
	"nuru-lsp/server"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/AvicennaJr/Nuru/lexer"
	"github.com/AvicennaJr/Nuru/parser"
	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
	sitter "github.com/smacker/go-tree-sitter"
)

var TUMIAS []string = []string{
	"os", "muda", "mtandao", "jsoni", "hisabati",
}

type ErrorMapLineNumbers = map[uint][]string

// Hold information on .nr file
type Data struct {
	File     string
	Version  uint64
	Errors   ErrorMapLineNumbers
	Content  []string
	RootTree *sitter.Node
}

var Pages = make(map[string]Data)
var PagesMutext = sync.Mutex{}

func ParseTree(lines []string) (*sitter.Node, error) {
	data := strings.Join(lines, "\n")
	data = data + "\n"
	node, err := sitter.ParseCtx(context.Background(),
		[]byte(data),
		nuru_tree_sitter.GetLanguage())
	return node, err
}

type ClosesNodeNotFound string

func (e ClosesNodeNotFound) Error() string {
	return string(e)
}

func traverseTreeToClosestNamedNode(node *sitter.Node, position defines.Position) (*sitter.Node, error) {
	row := position.Line
	//check if start happens ahead
	if node.StartPoint().Row > uint32(row) || node.EndPoint().Row < uint32(row) {
		fmt.Println("NOT WITHIN ROW", node.StartPoint(), node.EndPoint(), position, node.String())
		return nil, ClosesNodeNotFound("node out happening ahead position")
	}

	count := node.ChildCount()

	//meaning we have reached the last node
	if count == 0 {
		return node, nil
	}

	var resultNode *sitter.Node = nil
	for i := uint32(0); i < count; i++ {
		result, err := traverseTreeToClosestNamedNode(node.Child(int(i)), position)
		if err != nil {
			continue
		}
		if resultNode == nil {
			resultNode = result
			continue
		}
		distance := (result.StartPoint().Column - uint32(position.Character)) +
			(uint32(position.Character) - result.EndPoint().Column)
		distanceResultNode := (resultNode.StartPoint().Column - uint32(position.Character)) +
			(uint32(position.Character) - resultNode.EndPoint().Column)

		fmt.Println("CURRENT", distance, "LAST RESULT", distanceResultNode)
		fmt.Println("CURRENT: ", result.Type(),
			"LAST RESULT", resultNode.Type())

		if distance < distanceResultNode {
			resultNode = result
		}
	}

	if resultNode == nil {
		return nil, ClosesNodeNotFound("Closes Node not found")
	}
	return resultNode, nil
}

func (d *Data) GetModulesInDir() []string {
	modules := []string{}
	filepath.Walk(filepath.Dir(d.File), func(path string, info os.FileInfo, _err error) error {
		name := strings.Split(info.Name(), ".")
		ending := name[len(name)-1]
		if !info.IsDir() && (ending == "nr" || ending == "sr") {
			var data_to_use *Data
			if data, found := Pages[path]; found == false {
				bytes, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				data := strings.Split(string(bytes), "\n")
				data_to_use, err = NewData(path, 0, data)
				if err != nil {
					return nil
				}
			} else {
				data_to_use = &data
			}

			root := data_to_use.RootTree
			q, err := sitter.NewQuery([]byte(HII_NI_PAKEJI), nuru_tree_sitter.GetLanguage())
			if err != nil {
				fmt.Println("WTF", err)
				return nil
			}
			qc := sitter.NewQueryCursor()
			qc.Exec(q, root)

			if _, ok := qc.NextMatch(); ok == true {
				modules = append(modules, name[0])
			}
		}
		return nil
	})

	return modules
}

func (d *Data) getAllVariablesAndFunctions() *[]defines.CompletionItem {
	return nil
}

func (d *Data) TreeSitterCompletions(params *defines.CompletionParams) (*[]defines.CompletionItem, error) {
	fmt.Println("GOT ", d.RootTree.String())

	//query for possible node type at position of completions
	closestNode, err := traverseTreeToClosestNamedNode(d.RootTree, params.Position)
	if err != nil {
		fmt.Printf("SHIT %s", err)
		return nil, err
	}

	parent := closestNode.Parent()
	if parent != nil {
		fmt.Println("CLOSEST:", closestNode.Type(), parent.Type())
	}
	//do for tumia completions
	if parent != nil && parent.Type() == "pakeji_tumia_statement" {
		//get the identifier if any, then search files in same directory
		// if they match value given, also check default tumias available
		// by the nuru native.
		if closestNode.Type() == "identifier" {
			identifier := closestNode.Content([]byte(strings.Join(d.Content, "")))
			completionItems := []defines.CompletionItem{}
			for _, c := range TUMIAS {
				if strings.Contains(c, identifier) {
					detail := fmt.Sprintf("pakeji %s", c)
					kind := defines.CompletionItemKind(defines.CompletionItemKindModule)
					completionItems = append(completionItems,
						defines.CompletionItem{
							Label:      c,
							Detail:     &detail,
							InsertText: &c,
							Kind:       &kind,
						})
				}
			}

			for _, module := range d.GetModulesInDir() {
				if strings.Contains(module, identifier) {
					detail := fmt.Sprintf("pakeji %s", module)
					kind := defines.CompletionItemKind(defines.CompletionItemKindModule)
					completionItems = append(completionItems,
						defines.CompletionItem{
							Label:      module,
							Detail:     &detail,
							InsertText: &module,
							Kind:       &kind,
						})
				}
			}

			return &completionItems, nil
		}

		// If empty, then show all default tumias, and any file that is a pakeji in
		// the same directory then return
		if closestNode.Type() == "tumia" {
			completionItems := []defines.CompletionItem{}
			for _, c := range TUMIAS {
				detail := fmt.Sprintf("pakeji %s", c)
				kind := defines.CompletionItemKind(defines.CompletionItemKindModule)
				completionItems = append(completionItems,
					defines.CompletionItem{
						Label:      c,
						Detail:     &detail,
						InsertText: &c,
						Kind:       &kind,
					})
			}

			for _, module := range d.GetModulesInDir() {
				detail := fmt.Sprintf("pakeji %s", module)
				kind := defines.CompletionItemKind(defines.CompletionItemKindModule)
				completionItems = append(completionItems,
					defines.CompletionItem{
						Label:      module,
						Detail:     &detail,
						InsertText: &module,
						Kind:       &kind,
					})
			}

			return &completionItems, nil
		}
	}
	closestType := closestNode.Type()
	switch closestType {
	case "ending":
		completions := d.getAllVariablesAndFunctions()
		return completions, nil
	default:
		fmt.Println("SHHHIIIT")
	}
	return nil, ClosesNodeNotFound("Failed to find element")
}

func NewData(file string, version uint64, content []string) (*Data, error) {
	root, err := ParseTree(content)
	if err != nil {
		return nil, err
	}

	Pages[file] = Data{
		File:     file,
		Version:  version,
		Errors:   make(ErrorMapLineNumbers, 0),
		Content:  content,
		RootTree: root,
	}

	data := Pages[file]

	return &data, nil
}

func parseErrorFromParser(error string) (uint, *string, *error) {

	startPos := 0
	numPosStart := 0
	started := false
	intString := ""
	lineString := ""

	logs.Printf("Debug %s", error)

	for startPos < len(error) {
		if unicode.IsDigit(rune(error[startPos])) && started == false {
			started = true
			numPosStart = startPos
		}
		if !unicode.IsDigit((rune(error[startPos]))) && started == true {
			intString = error[numPosStart:startPos]
			startPos += 2
			break
		}
		startPos += 1
	}

	lineString = error[startPos:]
	pos, er := strconv.Atoi(intString)
	if er != nil {
		err := errors.New("Failed to parse number")
		return 0, nil, &err
	}
	if len(lineString) == 0 {
		err := errors.New("Failed to parse error message")
		return 0, nil, &err
	}

	return uint(pos), &lineString, nil
}

func OnDidClose(ctx context.Context, req *defines.DidCloseTextDocumentParams) (err error) {
	return nil
}

func parseAndNotifyErrors(doc *Data, uri defines.DocumentUri) {
	//remove all previous errors
	doc.Errors = make(ErrorMapLineNumbers, 0)
	fileData := strings.Join(doc.Content, "\n")
	l := lexer.New(fileData)
	p := parser.New(l)
	p.ParseProgram()
	//if errors, update doc
	if len(p.Errors()) > 0 {
		for _, e := range p.Errors() {
			pos, line, err := parseErrorFromParser(e)
			if err != nil {
				logs.Printf("Error parsing errors: %s\n", *err)
				return
			}
			errorsList := doc.Errors[pos]
			doc.Errors[pos] = append(errorsList, *line)
		}
	}
	//if errors not empty, now send them over to client
	diagnostics := make([]defines.Diagnostic, 0)
	for k, v := range doc.Errors {
		for _, e := range v {
			var endChar uint = 0
			if k < uint(len(doc.Content)) {
				endChar = uint(len(doc.Content[k-1]))
			}
			diagnostics = append(diagnostics, defines.Diagnostic{
				Message: e,
				Range: defines.Range{
					Start: defines.Position{
						Line:      k,
						Character: 0,
					},
					End: defines.Position{
						Line:      k,
						Character: endChar,
					},
				},
			})
		}
	}
	publishDiag := defines.PublishDiagnosticsParams{
		Uri:         uri,
		Diagnostics: diagnostics,
	}
	server.Notify(server.Server, "textDocument/publishDiagnostics", publishDiag)
	//now go line after line adding variables to scope
}

func OnDocOpen(ctx context.Context, req *defines.DidOpenTextDocumentParams) (err error) {
	PagesMutext.Lock()
	defer PagesMutext.Unlock()

	file := string(req.TextDocument.Uri)
	parsed, err := url.Parse(file)
	if err != nil {
		return nil
	}

	//check if it exists
	_, err = os.Stat(parsed.Path)
	if os.IsNotExist(err) {
		return nil
	}

	//we reach here means it exists, so open file and read it line by line

	//read content of file line by line
	content := strings.Split(req.TextDocument.Text, "\n")

	if len(content) == 0 {
		//empty file
		return nil
	}
	doc, err := NewData(parsed.Path, 0, content)
	if err != nil {
		return nil
	}

	//do diagnostics here on the file
	parseAndNotifyErrors(doc, req.TextDocument.Uri)

	logs.Printf("NURULSP DONE Opened file-> %s\n", parsed.Path)
	return nil
}

func OnDataChange(ctx context.Context, req *defines.DidChangeTextDocumentParams) error {
	PagesMutext.Lock()
	defer PagesMutext.Unlock()

	file := string(req.TextDocument.Uri)

	doc, found := Pages[file]
	if !found {
		content := []string{}
		for _, v := range req.ContentChanges {
			content = append(content, fmt.Sprint(v.Text))
		}
		docnew, _ := NewData(string(req.TextDocument.Uri), uint64(req.TextDocument.Version), content)
		doc = *docnew

	} else {
		if doc.Version < uint64(req.TextDocument.Version) {
			doc.Version = uint64(req.TextDocument.Version)
			content := []string{}
			for _, v := range req.ContentChanges {
				content = append(content, fmt.Sprint(v.Text))
			}
			doc.Content = content
		}
	}

	parseAndNotifyErrors(&doc, req.TextDocument.Uri)

	Pages[file] = doc

	return nil
}

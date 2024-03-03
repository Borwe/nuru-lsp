package data

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	nuru_tree_sitter "nuru-lsp/nuru-tree-sitter"
	"nuru-lsp/server"
	"os"
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
	node, err := sitter.ParseCtx(context.Background(),
		[]byte(strings.Join(lines, "")),
		nuru_tree_sitter.GetLanguage())
	return node, err
}

type ClosesNodeNotFound string

func (e ClosesNodeNotFound) Error() string {
	return string(e)
}

func traverseTreeToClosestNamedNode(node *sitter.Node, position defines.Position) (*sitter.Node, error) {
	row := position.Line
	collumn := position.Character
	//check if start happens ahead
	if node.StartPoint().Row > uint32(row) || node.EndPoint().Row < uint32(row) {
		fmt.Println("FUCK AHEAD", node.StartPoint(), node.EndPoint())
		return nil, ClosesNodeNotFound("node out happening ahead position")
	}
	// check if end happens before 
	if node.StartPoint().Column > uint32(collumn) || node.EndPoint().Column < uint32(collumn) {
		fmt.Println("FUCK BEFORE", node.StartPoint(), node.EndPoint())
		return nil, ClosesNodeNotFound("node out happening before position")
	}


	count := node.ChildCount()

	//meaning we have reached the last node
	if count == 0 {
		return node, nil
	}

	for i := uint32(0); i <count; i++ {
		resultNode, _ := traverseTreeToClosestNamedNode(node.Child(int(i)), position)
		if resultNode != nil {
			return resultNode, nil
		}
	}

	return nil, ClosesNodeNotFound("Closes Node not found")
}

func (d *Data) TreeSitterCompletions(params *defines.CompletionParams) (*[]defines.CompletionItem, error) {
	node, err := ParseTree(d.Content)
	if err != nil {
		return nil, err
	}
	fmt.Println("GOT ",node.String())
	d.RootTree = node

	//query for possible node type at position of completions
	closestNode, err := traverseTreeToClosestNamedNode(node, params.Position)
	fmt.Printf("SHIT %s", err)
	if err != nil {
		return nil, err
	}
	//do for tumia completions
	if closestNode.Type() == "pakeji_tumia_statement" {
		//get the identifier if any, then search files in same directory
		// if they match value given, also check default tumias available
		// by the nuru native.
		// If empty, then show all default tumias, and any file that is a pakeji in
		// the same directory
		//then return
	}
	return nil, nil
}

func NewData(file string, version uint64, content []string) Data {
	return Data{
		File:    file,
		Version: version,
		Errors:  make(ErrorMapLineNumbers, 0),
		Content: content,
	}
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
	doc := NewData(parsed.Path, 0, content)
	//do diagnostics here on the file
	parseAndNotifyErrors(&doc, req.TextDocument.Uri)

	//store
	Pages[parsed.Path] = doc

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
		doc = NewData(string(req.TextDocument.Uri), uint64(req.TextDocument.Version), content)

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

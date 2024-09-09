package data

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"nuru-lsp/server"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/AvicennaJr/Nuru/ast"
	"github.com/AvicennaJr/Nuru/lexer"
	"github.com/AvicennaJr/Nuru/parser"
	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
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
	RootTree *ast.Node
}

var Pages = make(map[string]Data)
var PagesMutext = sync.Mutex{}

func ParseTree(parser *parser.Parser) (ast.Node, []string) {
	ast := parser.ParseProgram()
	errorsList := parser.Errors()
	return ast, errorsList
}

type ClosesNodeNotFound string

func (e ClosesNodeNotFound) Error() string {
	return string(e)
}

func parseWordBackFromPosition(line []rune, pos int) *string{
	word := []rune{}
	startedGettingWords := false
	for i:=pos; i>=0; i-- {
		if line[i] == ' ' && startedGettingWords {
			break;
		}

		if !startedGettingWords && line[i] != ' '{
			startedGettingWords = true
		}
		word = append([]rune{line[i]},word...)
	}
	tmp := string(word)
	return &tmp
}

func (d *Data) Completions(completeParams *defines.CompletionParams) (*[]defines.CompletionItem, error){
	//get current word, otherwise get previous
	var word *string = nil

	for no, line := range d.Content{
		if no != int(completeParams.Position.Line-1){
			continue
		}
		//-1 because files consider column 1 as index 0
		charPos := completeParams.Position.Character-2
		runed := []rune(line)
		if charPos>0 {
			//parse the full word,
			word = parseWordBackFromPosition(runed, int(charPos))
		}
	}

	if(word==nil){
		return nil, errors.New("NOT IMPLEMENTED BASIC COMPLETION") 
	}
	return nil, errors.New("NOT IMPLEMENTED")
}

func (d *Data) getAllVariablesAndFunctions() *[]defines.CompletionItem {
	result := make([]defines.CompletionItem, 1)
	return &result
}

func NewData(file string, version uint64, content []string) (*Data, error) {
	lexer := lexer.New(strings.Join(content, ""))
	parser := parser.New(lexer)
	root, errs := ParseTree(parser)

	data := Data{
		File:     file,
		Version:  version,
		Errors:   make(ErrorMapLineNumbers, 0),
		Content:  content,
		RootTree: &root,
	}

	if len(errs) > 0 {
		notifyErrors(&data, errs)
	}

	Pages[file] = data
	data = Pages[file]

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

func notifyErrors(doc *Data, errors []string) {
	if len(errors) > 0 {
		for _, e := range errors {
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
		Uri:         defines.DocumentUri(doc.File),
		Diagnostics: diagnostics,
	}
	server.Notify(server.Server, "textDocument/publishDiagnostics", publishDiag)
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

func ReadLine(in string) string {
	out := strings.Split(in, "\r\n")
	out = strings.Split(out[0], "\n")
	return out[0]
}

func ReadContents(text string) []string {
	content := []string{}
	lines := strings.Split(text, "\r\n")
	for _, line := range lines {
		sublines := strings.Split(line, "\n")
		for _, l := range sublines {
			content = append(content, l)
		}
	}
	return content
}

func OnDocOpen(ctx context.Context, req *defines.DidOpenTextDocumentParams) (err error) {
	PagesMutext.Lock()
	defer PagesMutext.Unlock()

	file := string(req.TextDocument.Uri)
	parsed, err := url.Parse(file)
	if err != nil {
		return nil
	}

	//we reach here means it exists, so open file and read it line by line
	//read content of file line by line
	content := ReadContents(req.TextDocument.Text)

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
		logs.Println("WEWEWEWE", ReadContents(fmt.Sprint(req.ContentChanges[0].Text)))
		logs.Println("RERERERE", len(req.ContentChanges))
		content = append(content, ReadContents(fmt.Sprint(req.ContentChanges[0].Text))...)
		docnew, _ := NewData(string(req.TextDocument.Uri), uint64(req.TextDocument.Version), content)
		doc = *docnew

	} else {
		if doc.Version < uint64(req.TextDocument.Version) {
			doc.Version = uint64(req.TextDocument.Version)
			content := []string{}
			content = append(content, ReadContents(fmt.Sprint(req.ContentChanges[0].Text))...)
			doc.Content = content
		}
	}

	parseAndNotifyErrors(&doc, req.TextDocument.Uri)

	Pages[file] = doc

	return nil
}

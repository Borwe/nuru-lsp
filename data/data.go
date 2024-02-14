package data

import (
	"context"
	"errors"
	"fmt"
	"net/url"
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
)

type ErrorMapLineNumbers = map[uint][]string

/*
  - Hold variable information, for lookup
    when invoking onComplete
*/
type VariableOrFunction struct {
	Line       uint
	Name       string
	DocString  *string
	File       string
	ScopeStart *uint
	ScopeEnd   *uint
	IsScope    bool
	IsPakeji   bool
	IsFunction bool
}

func addDocString(varOrFunc *VariableOrFunction) {
	content := Pages[varOrFunc.File].Content
	if varOrFunc.Line-1 >= 0 {
		trimed := strings.Trim(content[varOrFunc.Line-1], " ")
		if len(trimed) >= 2 {
			//single line comment
			if trimed[0:2] == "//" {
				comment := trimed[2:]
				varOrFunc.DocString = &comment
				return
			}

			//check if multi line comment
			if trimed[len(trimed)-2:] == "*/" {
				if trimed[0:2] == "/*" {
					//means double line quotes is in same line
					comment := trimed[2 : len(trimed)-2]
					varOrFunc.DocString = &comment
					return
				}

				comment := make([]string, 0)
				comment = append(comment, trimed[:2])
				startPos := varOrFunc.Line - 2
				for startPos >= 0 {
					trimed = strings.Trim(content[startPos], " ")
					if trimed[0:2] == "/*" {
						comment = append([]string{trimed[2:]}, comment...)
						final := strings.Join(comment, "\n")
						varOrFunc.DocString = &final
						return
					}
					startPos -= 1
				}
			}
		}
	}
}

// Hold top of tree
type Top struct {
	Pakeji Pakeji
	Items  []FuncVar
}

type FuncVar struct {
	Items     []FuncVar
	Line      int64
	Name      string
	StartDecl int64
	EndDecl   int64
	IsScope   bool
}

type Pakeji struct {
	Items []FuncVar
}

// Hold information on .nr file
type Data struct {
	File    string
	Version uint64
	Errors  ErrorMapLineNumbers
	Content []string
	Tree    *Top
}

var Pages = make(map[string]Data)
var PagesMutext = sync.Mutex{}

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

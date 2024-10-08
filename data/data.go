package data

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"nuru-lsp/server"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/NuruProgramming/Nuru/ast"
	"github.com/NuruProgramming/Nuru/lexer"
	"github.com/NuruProgramming/Nuru/parser"
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

func parseWordBackFromPosition(line []rune, pos int) *string {
	word := []rune{}
	startedGettingWords := false
	for i := pos; i >= 0; i-- {
		if line[i] == ' ' && startedGettingWords {
			break
		}

		if !startedGettingWords && line[i] != ' ' {
			startedGettingWords = true
		} else if !startedGettingWords {
			continue
		}
		word = append([]rune{line[i]}, word...)
	}
	tmp := string(word)
	return &tmp
}

func checkFileIsPackage(dir string, file fs.DirEntry) bool {
	fpath := path.Join(dir, file.Name())
	logs.Println("DIRFPATH:", fpath)
	data, ok := Pages[fpath]
	if !ok {
		linesBytes, err := os.ReadFile(fpath)
		if err != nil {
			return false
		}
		lines := strings.Split(string(linesBytes), "\n")
		tmp, err, errs := NewData(fpath, 0, lines)
		if err != nil || len(errs) > 0 {
			return false
		}
		Pages[fpath] = *tmp
		data = Pages[fpath]
	}
	pakejiAsts := &[]*ast.Package{}
	getAsts(*data.RootTree, &pakejiAsts)

	if len(*pakejiAsts) > 0 {
		return true
	}

	return false
}

func getAsts[T ast.Node](node ast.Node, result **[]T) {
	switch node := node.(type) {
	case T:
		tmp := append(**result, node)
		*result = &tmp
		break
	case *ast.Program:
		for _, stmt := range node.Statements {
			getAsts(stmt, result)
		}
		break

	case *ast.ExpressionStatement:
		getAsts(node.Expression, result)
		break
	case *ast.IntegerLiteral:
	case *ast.FloatLiteral:
	case *ast.Boolean:
		break
	case *ast.PrefixExpression:
		getAsts(node.Right, result)
		break
	case *ast.InfixExpression:
		getAsts(node.Left, result)
		getAsts(node.Right, result)
		break
	case *ast.PostfixExpression:
		break
	case *ast.BlockStatement:
		for _, stmt := range node.Statements {
			getAsts(stmt, result)
		}
		break
	case *ast.IfExpression:
		getAsts(node.Condition, result)
		getAsts(node.Alternative, result)
		getAsts(node.Consequence, result)
		break
	case *ast.ReturnStatement:
		getAsts(node.ReturnValue, result)
		break
	case *ast.LetStatement:
		getAsts(node.Value, result)
		break
	case *ast.Identifier:
		break
	case *ast.FunctionLiteral:
		getAsts(node.Body, result)
		for _, stmt := range node.Parameters {
			getAsts(stmt, result)
		}
		for _, stmt := range node.Defaults {
			getAsts(stmt, result)
		}
		break
	case *ast.MethodExpression:
		getAsts(node.Object, result)
		getAsts(node.Method, result)
		for _, stmt := range node.Defaults {
			getAsts(stmt, result)
		}
		for _, stmt := range node.Arguments {
			getAsts(stmt, result)
		}
		break
	case *ast.Import:
		for _, stmt := range node.Identifiers {
			getAsts(stmt, result)
		}
		break
	case *ast.CallExpression:
		getAsts(node.Function, result)
		for _, stmt := range node.Arguments {
			getAsts(stmt, result)
		}
		break
	case *ast.StringLiteral:
		break
	case *ast.At:
		break
	case *ast.ArrayLiteral:
		for _, stmt := range node.Elements {
			getAsts(stmt, result)
		}
		break
	case *ast.IndexExpression:
		getAsts(node.Left, result)
		getAsts(node.Index, result)
		break
	case *ast.DictLiteral:
		for stmt1, stmt2 := range node.Pairs {
			getAsts(stmt1, result)
			getAsts(stmt2, result)
		}
		break
	case *ast.WhileExpression:
		getAsts(node.Condition, result)
		getAsts(node.Consequence, result)
		break
	case *ast.Break:
	case *ast.Continue:
		break
	case *ast.SwitchExpression:
		getAsts(node.Value, result)
		for _, stmt := range node.Choices {
			getAsts(stmt, result)
		}
		break
	case *ast.Null:
		break
	case *ast.ForIn:
		getAsts(node.Iterable, result)
		getAsts(node.Block, result)
		break
	case *ast.Package:
		getAsts(node.Name, result)
		getAsts(node.Block, result)
		break
	case *ast.PropertyExpression:
		getAsts(node.Object, result)
		getAsts(node.Property, result)
		break
	case *ast.PropertyAssignment:
		getAsts(node.Name, result)
		getAsts(node.Value, result)
		break
	case *ast.Assign:
		getAsts(node.Name, result)
		getAsts(node.Value, result)
		break
	case *ast.AssignEqual:
		getAsts(node.Left, result)
		getAsts(node.Value, result)
		break
	case *ast.AssignmentExpression:
		getAsts(node.Left, result)
		getAsts(node.Value, result)
		break
	}
}

func getNuruFilesInDir(dir string) []fs.DirEntry {
	result := []fs.DirEntry{}
	files, error := os.ReadDir(dir)
	if error == nil {
		//meaning we have files
		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				continue
			}
			if info.IsDir() {
				continue
			}

			fName := file.Name()
			if len(fName) < 3 {
				continue
			}

			extention := fName[len(fName)-2:]
			if extention == "nr" {
				result = append(result, file)
			}
		}
	}
	return result
}

func (d *Data) Completions(completeParams *defines.CompletionParams) (*[]defines.CompletionItem, error) {
	//get current word, otherwise get previous
	var word *string = nil
	var prevWord *string = nil

	for no, line := range d.Content {
		if no != int(completeParams.Position.Line) {
			continue
		}
		//-1 because files consider column 1 as index 0
		charPos := completeParams.Position.Character - 1
		runed := []rune(line)
		if charPos > 0 {
			//parse the full word,
			word = parseWordBackFromPosition(runed, int(charPos))
			if word != nil {
				prevWord = parseWordBackFromPosition(runed, int(charPos)-len(*word))
			}
		}
	}

	if word == nil {
		return nil, errors.New("NOT IMPLEMENTED BASIC COMPLETION")
	}

	switch *word {
	case "tumia":
		//get all files in directory of current data
		logs.Println("FILE COMPLETING:", d.File)
		packajiFiles := []string{}
		dir := path.Dir(d.File)
		files := getNuruFilesInDir(dir)
		for _, file := range files {
			if checkFileIsPackage(dir, file) {
				name := file.Name()
				packajiFiles = append(packajiFiles, name[:len(name)-3])
			}
		}
		logs.Println("PAKEJIS:", packajiFiles)
		completions := []defines.CompletionItem{}
		pakejiInfo := "Ni pakeji"
		pakejiKind := defines.CompletionItemKindFile
		for _, pakeji := range packajiFiles {
			completions = append(completions, defines.CompletionItem{
				Label:        pakeji,
				Kind:         &pakejiKind,
				LabelDetails: &defines.CompletionItemLabelDetails{Detail: &pakejiInfo},
			})
		}
		for _, tumia := range TUMIAS {
			completions = append(completions, defines.CompletionItem{
				Label:        tumia,
				LabelDetails: &defines.CompletionItemLabelDetails{Detail: &pakejiInfo},
				Kind:         &pakejiKind,
			})
		}
		for file, page := range Pages {
			fileDir := path.Dir(file)
			if fileDir == dir {
				checks := &[]*ast.Package{}
				getAsts(*page.RootTree, &checks)
				if len(*checks) > 0 {
					name := filepath.Base(file)
					name = name[0 : len(name)-3]
					same := false
					for _, completion := range completions {
						if completion.Label == name {
							same = true
							break
						}
					}

					if !same {
						completions = append(completions, defines.CompletionItem{
							Label:        name,
							LabelDetails: &defines.CompletionItemLabelDetails{Detail: &pakejiInfo},
							Kind:         &pakejiKind,
						})
					}
				}
			}
		}
		return &completions, nil
	default:
		if *prevWord == "tumia" {
			files := getNuruFilesInDir(path.Dir(d.File))
			completions := []defines.CompletionItem{}
			pakejiKind := defines.CompletionItemKindFile
			pakejiInfo := "Ni pakeji"
			for _, file := range files {
				if checkFileIsPackage(path.Dir(d.File), file) {
					name := file.Name()
					if strings.Contains(name, *word) {
						completions = append(completions, defines.CompletionItem{
							Label:        name[0 : len(name)-3],
							Kind:         &pakejiKind,
							LabelDetails: &defines.CompletionItemLabelDetails{Detail: &pakejiInfo},
						})
					}
				}
			}
			for _, pakeji := range TUMIAS {
				if strings.Contains(pakeji, *word) {
					completions = append(completions, defines.CompletionItem{
						Label:        pakeji,
						Kind:         &pakejiKind,
						LabelDetails: &defines.CompletionItemLabelDetails{Detail: &pakejiInfo},
					})
				}
			}
			dir := path.Dir(d.File)
			for file, page := range Pages {
				fileDir := path.Dir(file)
				name := filepath.Base(file)
				name = name[0 : len(name)-3]
				if fileDir == dir && strings.Contains(name, *word) {
					checks := &[]*ast.Package{}
					getAsts(*page.RootTree, &checks)
					if len(*checks) > 0 {
						same := false
						for _, completion := range completions {
							if completion.Label == name {
								same = true
								break
							}
						}

						if !same {
							completions = append(completions, defines.CompletionItem{
								Label:        name,
								LabelDetails: &defines.CompletionItemLabelDetails{Detail: &pakejiInfo},
								Kind:         &pakejiKind,
							})
						}
					}
				}
			}
			return &completions, nil
		}
		return nil, errors.New("NOT IMPLEMENTED")
	}

}

func (d *Data) getAllVariablesAndFunctions() *[]defines.CompletionItem {
	result := make([]defines.CompletionItem, 1)
	return &result
}

func NewData(file string, version uint64, content []string) (*Data, error, []string) {
	lexer := lexer.New(strings.Join(content, ""))
	parser := parser.New(lexer)
	root, errs := ParseTree(parser)

	fileUrl, err := url.Parse(file)
	if err != nil {
		return nil, err, errs
	}

	file = fileUrl.Path
	if strings.HasPrefix(file, "/") && filepath.IsAbs(file[1:]) {
		file = file[1:]
	}

	data := Data{
		File:     file,
		Version:  version,
		Errors:   make(ErrorMapLineNumbers, 0),
		Content:  content,
		RootTree: &root,
	}

	Pages[file] = data
	data = Pages[file]

	return &data, nil, errs
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
	logs.Println("DIAGS:", len(diagnostics), len(errors))
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

	//if already previously opened by other methods just return here
	if _, ok := Pages[file]; ok {
		return nil
	}

	//we reach here means it exists, so open file and read it line by line
	//read content of file line by line
	content := ReadContents(req.TextDocument.Text)

	doc, err, errs := NewData(file, 0, content)
	if err != nil {
		return nil
	}

	//do diagnostics here on the file
	notifyErrors(doc, errs)

	logs.Printf("NURULSP DONE Opened file-> %s\n", file)
	return nil
}

func OnDataChange(ctx context.Context, req *defines.DidChangeTextDocumentParams) error {
	PagesMutext.Lock()
	defer PagesMutext.Unlock()

	file := string(req.TextDocument.Uri)

	errs := []string{}

	doc, found := Pages[file]
	if !found {
		content := []string{}
		logs.Println("WEWEWEWE", ReadContents(fmt.Sprint(req.ContentChanges[0].Text)))
		logs.Println("RERERERE", len(req.ContentChanges))
		content = append(content, ReadContents(fmt.Sprint(req.ContentChanges[0].Text))...)
		docnew, _, errsDoc := NewData(string(req.TextDocument.Uri), uint64(req.TextDocument.Version), content)
		errs = errsDoc
		doc = *docnew

	} else {
		if doc.Version < uint64(req.TextDocument.Version) {
			content := []string{}
			content = append(content, ReadContents(fmt.Sprint(req.ContentChanges[0].Text))...)
			docnew, _, errsDoc := NewData(string(req.TextDocument.Uri), uint64(req.TextDocument.Version), content)
			errs = errsDoc
			doc = *docnew
		}
	}

	notifyErrors(&doc, errs)

	Pages[file] = doc

	return nil
}

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
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/NuruProgramming/Nuru/ast"
	"github.com/NuruProgramming/Nuru/lexer"
	"github.com/NuruProgramming/Nuru/module"
	"github.com/NuruProgramming/Nuru/parser"
)

var TUMIAS []string = []string{}

func init() {
	// Get all default tumias exported by NURU
	for tumia := range module.Mapper {
		TUMIAS = append(TUMIAS, tumia)
	}
}

type ErrorMapLineNumbers = map[uint][]string

// Hold information on .nr file
type Data struct {
	File        string
	Version     uint64
	Errors      ErrorMapLineNumbers
	Content     []string
	RootTree    *ast.Node
	WorkingTree *ast.Node
}

var Pages = map[string]Data{}
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
		data = *tmp
	}
	pakejiAsts := []*ast.Package{}
	getAsts(*data.RootTree, &pakejiAsts)

	return len(pakejiAsts) > 0
}

func getAsts[T ast.Node](node ast.Node, result *[]T) {
	if v := reflect.ValueOf(node); !v.IsValid() || v.IsNil() {
		return
	}
	switch node := node.(type) {
	case T:
		tmp := append(*result, node)
		*result = tmp
		logs.Println("YOOOOOOOOOOOOOOOOOOOOOOOO")
		break
	case *ast.Program:
		logs.Println("Pogram")
		for _, stmt := range node.Statements {
			getAsts(stmt, result)
		}
		break
	case *ast.ExpressionStatement:
		logs.Println("ExpressionStatement")
		getAsts(node.Expression, result)
		break
	case *ast.IntegerLiteral:
	case *ast.FloatLiteral:
	case *ast.Boolean:
		logs.Println("Literal")
		break
	case *ast.PrefixExpression:
		logs.Println("PrefixExpression")
		getAsts(node.Right, result)
		break
	case *ast.InfixExpression:
		logs.Println("InfixExpression")
		getAsts(node.Left, result)
		getAsts(node.Right, result)
		break
	case *ast.PostfixExpression:
		logs.Println("PostfixExpression")
		break
	case *ast.BlockStatement:
		logs.Println("BlockStatement")
		for _, stmt := range node.Statements {
			getAsts(stmt, result)
		}
		break
	case *ast.IfExpression:
		logs.Println("IfExpression")
		getAsts(node.Condition, result)
		getAsts(node.Alternative, result)
		getAsts(node.Consequence, result)
		break
	case *ast.ReturnStatement:
		logs.Println("ReturnStatement")
		getAsts(node.ReturnValue, result)
		break
	case *ast.LetStatement:
		logs.Println("LetStatement")
		getAsts(node.Value, result)
		break
	case *ast.FunctionLiteral:
		logs.Println("FunctionLiteral")
		getAsts(node.Body, result)
		for _, stmt := range node.Parameters {
			getAsts(stmt, result)
		}
		for _, stmt := range node.Defaults {
			getAsts(stmt, result)
		}
		break
	case *ast.PropertyExpression:
		logs.Println("PropertyExpression")
		getAsts(node.Object, result)
		getAsts(node.Property, result)
		break
	case *ast.PropertyAssignment:
		logs.Println("PropertyAssignment")
		getAsts(node.Name, result)
		getAsts(node.Value, result)
		break
	case *ast.Assign:
		logs.Println("Assign")
		getAsts(node.Name, result)
		getAsts(node.Value, result)
		break
	case *ast.AssignEqual:
		logs.Println("AssignEqual")
		getAsts(node.Left, result)
		getAsts(node.Value, result)
		break
	case *ast.AssignmentExpression:
		logs.Println("AssignmentExpression")
		getAsts(node.Left, result)
		getAsts(node.Value, result)
		break
	case *ast.MethodExpression:
		logs.Println("MethodExpression")
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
		logs.Println("Import")
		for _, stmt := range node.Identifiers {
			var nd ast.Node = stmt
			getAsts(nd, result)
		}
		break
	case *ast.CallExpression:
		logs.Println("CallExpression")
		getAsts(node.Function, result)
		for _, stmt := range node.Arguments {
			getAsts(stmt, result)
		}
		break
	case *ast.StringLiteral:
		logs.Println("StringLiteral")
		break
	case *ast.At:
		logs.Println("At")
		break
	case *ast.ArrayLiteral:
		logs.Println("ArrayLiteral")
		for _, stmt := range node.Elements {
			getAsts(stmt, result)
		}
		break
	case *ast.IndexExpression:
		logs.Println("IndexExpression")
		getAsts(node.Left, result)
		getAsts(node.Index, result)
		break
	case *ast.DictLiteral:
		logs.Println("DictLiteral")
		for stmt1, stmt2 := range node.Pairs {
			getAsts(stmt1, result)
			getAsts(stmt2, result)
		}
		break
	case *ast.WhileExpression:
		logs.Println("WhileExpression")
		getAsts(node.Condition, result)
		getAsts(node.Consequence, result)
		break
	case *ast.Break:
	case *ast.Continue:
		logs.Println("BreakOrContinue")
		break
	case *ast.SwitchExpression:
		logs.Println("SwitchExpression")
		getAsts(node.Value, result)
		for _, stmt := range node.Choices {
			getAsts(stmt, result)
		}
		break
	case *ast.Null:
		logs.Println("Null")
		break
	case *ast.ForIn:
		logs.Println("ForIn")
		getAsts(node.Iterable, result)
		getAsts(node.Block, result)
		break
	case *ast.Package:
		logs.Println("Package")
		getAsts(node.Name, result)
		getAsts(node.Block, result)
		break
	case *ast.Identifier:
		logs.Println("Identifier")
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

func getTumiaIdentifiers(node *ast.Node) []*ast.Identifier {
	tumiaIdentifiers := []*ast.Identifier{}
	tumiaLists := []*ast.Import{}
	getAsts(*node, &tumiaLists)
	for _, tumias := range tumiaLists {
		getAsts(tumias, &tumiaIdentifiers)
	}
	return tumiaIdentifiers
}

func (d *Data) getCompletions(word *string) (*[]defines.CompletionItem, error) {
	completions := []defines.CompletionItem{}
	//get all the tumias
	tumiaIdentifiers := getTumiaIdentifiers(d.RootTree)
	for _, val := range tumiaIdentifiers {
		kind := defines.CompletionItemKindFile
		label := val.String()
		detail := "Ni pakeji"
		logs.Println("TUMIAS NAMED:", label, "VAL", detail)

		completions = append(completions, defines.CompletionItem{
			Kind:  &kind,
			Label: label,
			LabelDetails: &defines.CompletionItemLabelDetails{
				Detail: &detail,
			},
		})
	}

	//get all Identifiers in current file
	letEquals := []*ast.LetStatement{}
	getAsts(*d.RootTree, &letEquals)
	assignmentEquals := []*ast.Assign{}
	getAsts(*d.RootTree, &assignmentEquals)

	//get variables
	for _, val := range letEquals {

		funcKind := defines.CompletionItemKindFunction
		label := val.Name.String()
		detail := ""
		if val.Value != nil {
			detail = val.String()
		}
		logs.Println("NAMED:", label, "VAL", detail)

		completions = append(completions, defines.CompletionItem{
			Kind:  &funcKind,
			Label: label,
			LabelDetails: &defines.CompletionItemLabelDetails{
				Detail: &detail,
			},
		})
	}
	for _, val := range assignmentEquals {

		kind := defines.CompletionItemKindField
		label := val.Name.String()
		detail := ""
		if val.Value != nil {
			detail = val.String()
		}
		logs.Println("NAMED:", label, "VAL", detail)

		completions = append(completions, defines.CompletionItem{
			Kind:  &kind,
			Label: label,
			LabelDetails: &defines.CompletionItemLabelDetails{
				Detail: &detail,
			},
		})
	}

	if word != nil && *word == "" {
		return nil, errors.New(fmt.Sprint("Passed an empy string for completions:", *word, "As word"))
	} else if word == nil {
		return &completions, nil
	}

	finalCompletion := []defines.CompletionItem{}
	for _, completion := range completions {
		if strings.Contains(completion.Label, *word) {
			finalCompletion = append(finalCompletion, completion)
		}
	}
	return &finalCompletion, nil
}

func combineCompletions(completions []defines.CompletionItem,
	toAdd *[]defines.CompletionItem, filter *string) *[]defines.CompletionItem {
	if toAdd != nil {
		for _, item := range *toAdd {
			if filter == nil || *filter == "=" || item.Label == *filter {
				completions = append(completions, item)
			}
		}
	}
	//logs.Println("COMPLETION ITEMS:",completions)
	return &completions
}

func (d *Data) getPackageCompletions(word *string) (*[]defines.CompletionItem, error) {
	completions := []defines.CompletionItem{}
	//word should not be empty
	if len(*word) == 0 {
		return &completions, nil
	}

	pkg := (*word)[:len(*word)-1]
	found := false
	tumiaIdentifiers := getTumiaIdentifiers(d.RootTree)
	for _, tm := range tumiaIdentifiers {
		if tm.Value == pkg {
			found = true
		}
	}

	if !found {
		//meaning not an imported tumia, so no need to compelte
		return &completions, nil
	}

	found = false
	for _, std := range TUMIAS{
		if std == pkg {
			found = true
		}
	}

	if !found {
		//meaning not an std package completions
		getAndFillPakejiPages(path.Dir(d.File))
		var dToUse *Data = nil
		for _, data := range Pages {
			if strings.Contains(data.File, pkg+".nr") {
				dToUse = &data
				break
			}
		}

		pakejis := []*ast.Package{}
		getAsts(*dToUse.RootTree, &pakejis)
		if len(pakejis) == 0 {
			//meaning not a pakeji
			return &completions, nil
		}

		letEquals := []*ast.LetStatement{}
		getAsts(pakejis[0], &letEquals)
		assignmentEquals := []*ast.Assign{}
		getAsts(pakejis[0], &assignmentEquals)

		methodKind := defines.CompletionItemKindMethod
		propertyKind := defines.CompletionItemKindProperty
		for _, val := range letEquals {
			kind := propertyKind
			label := val.Name.String()
			detail := ""
			if val.Value != nil {
				detail = val.String()
				if _, ok := val.Value.(*ast.FunctionLiteral); ok {
					kind = methodKind
				}
			}
			completions = append(completions, defines.CompletionItem{
				Kind:  &kind,
				Label: label,
				LabelDetails: &defines.CompletionItemLabelDetails{
					Detail: &detail,
				},
			})
		}
		for _, val := range assignmentEquals {
			kind := propertyKind
			label := val.Name.String()
			detail := ""
			if val.Value != nil {
				detail = val.String()
				if _, ok := val.Value.(*ast.FunctionLiteral); ok {
					kind = methodKind
				}
			}
			completions = append(completions, defines.CompletionItem{
				Kind:  &kind,
				Label: label,
				LabelDetails: &defines.CompletionItemLabelDetails{
					Detail: &detail,
				},
			})
		}

		return &completions, nil
	}

	mod, found := module.Mapper[pkg]

	if found {
		funcKind := defines.CompletionItemKindMethod
		for methods := range mod.Functions {
			completions = append(completions, defines.CompletionItem{
				Label: methods,
				Kind:  &funcKind,
			})
		}
	}

	if pkg == "hisabati" {
		varKind := defines.CompletionItemKindProperty
		//hisabati is special as it has const variables too
		for cnst, val := range module.Constants {
			valstr := val.Inspect()
			completions = append(completions, defines.CompletionItem{
				Label: cnst,
				Kind:  &varKind,
				LabelDetails: &defines.CompletionItemLabelDetails{
					Detail: &valstr,
				},
			})
		}
	}

	return &completions, nil
}

func getAndFillPakejiPages(dir string) []string {
	packajiFiles := []string{}
	files := getNuruFilesInDir(dir)
	for _, file := range files {
		if checkFileIsPackage(dir, file) {
			name := file.Name()
			packajiFiles = append(packajiFiles, name[:len(name)-3])
		}
	}
	return packajiFiles
}

func (d *Data) Completions(completeParams *defines.CompletionParams,
	defaultCompletions *[]defines.CompletionItem) (*[]defines.CompletionItem,
	error) {
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

	//meaning we have no input from user to go by
	//so just get all idenfitiers available
	if (prevWord == nil && word == nil) || (*prevWord == "" && *word == "") {
		completes, err := d.getCompletions(nil)
		if err != nil {
			return defaultCompletions, nil
		}
		return combineCompletions(*completes, defaultCompletions, nil), nil
	}

	switch *word {
	case "tumia":
		//get all files in directory of current data
		logs.Println("TUMIA FILE COMPLETING:", d.File)
		dir := path.Dir(d.File)
		packajiFiles := getAndFillPakejiPages(dir)
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
			fileDir := path.Dir(page.File)
			if fileDir == dir && file != string(completeParams.TextDocument.Uri) {
				checks := []*ast.Package{}
				getAsts(*page.RootTree, &checks)
				if len(checks) > 0 {
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
			for _, page := range Pages {
				fileDir := path.Dir(page.File)
				name := filepath.Base(page.File)
				name = name[0 : len(name)-3]
				if fileDir == dir && strings.Contains(name, *word) {
					checks := []*ast.Package{}
					getAsts(*page.RootTree, &checks)
					if len(checks) > 0 {
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
		} else if word != nil && strings.Contains(*word, ".") {
			//pakeji completions
			return d.getPackageCompletions(word)
		} else if word != nil && *word != "" && !(prevWord != nil && *prevWord == "fanya") {
			completions, err := d.getCompletions(word)
			if err != nil {
				return defaultCompletions, err
			}
			logs.Println("PREVWORD:", *prevWord, "WORD:", *word)
			return combineCompletions(*completions, defaultCompletions, nil), nil
		}
		return nil, errors.New(fmt.Sprintf("%s prev->%s word->%s", "NOT IMPLEMENTED", *prevWord, *word))
	}
}

func (d *Data) getAllVariablesAndFunctions() *[]defines.CompletionItem {
	result := make([]defines.CompletionItem, 1)
	return &result
}

func NewData(file string, version uint64, content []string) (*Data, error, []string) {
	lexer := lexer.New(strings.Join(content, "\n"))
	parser := parser.New(lexer)
	root, errs := ParseTree(parser)

	fileUrl, err := url.Parse(file)
	if err != nil {
		return nil, err, errs
	}

	logs.Println("FILEOPENBEFORE:", file)
	filePath := fileUrl.Path
	if strings.HasPrefix(filePath, "/") && filepath.IsAbs(filePath[1:]) {
		filePath = filePath[1:]
	}

	logs.Println("FILEOPENAFTER:", filePath)
	data := Data{
		File:     filePath,
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
	
	// Fix URI scheme for LSP clients
	uri := doc.File
	if !strings.HasPrefix(uri, "file://") {
		uri = "file://" + uri
	}
	publishDiag := defines.PublishDiagnosticsParams{
		Uri:         defines.DocumentUri(uri),
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
	
	// Fix URI scheme for LSP clients
	uriStr := string(uri)
	if !strings.HasPrefix(uriStr, "file://") {
		uriStr = "file://" + uriStr
	}
	publishDiag := defines.PublishDiagnosticsParams{
		Uri:         defines.DocumentUri(uriStr),
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
		return err
	}

	//do diagnostics here on the file
	notifyErrors(doc, errs)

	logs.Printf("NURULSP DONE Opened file-> %s\n", file)
	return nil
}

func OnDataChange(ctx context.Context,
	req *defines.DidChangeTextDocumentParams) error {
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

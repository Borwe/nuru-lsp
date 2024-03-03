package completions

import (
	"context"
	"nuru-lsp/data"
	"strings"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
)

var functions = map[string]string{
	"andika": `Inatumika kuandika mistari kwa terminali
		mfano: andika(1,2,3) 
		itaandika: 1, 2, 3 
		katika terminali unachotumia`,
	"jaza": `Inatumika kupata mistari kutuko kwa mtu
		mfano: fanya jina = jaza("Andika Jina");
			andika(jina)
		mwelezo: ukijaza "Brian" na kufinya enter itaweka hiyo
		kwa variabu inaitwa jina, ukiiandika utaona imeandika
		"Brian" katika terminali`,
	"fungua": `Inatumika kufunugua file
		mfano: f = funugua("./kitu.txt")`,
	"aina": `kinatumika kutambua aina ya kitu
		mfano: aina(2)
		itaandika: "NAMBA"`,
}

var keywords = []string{
	"unda",
	"fanya",
	"kweli",
	"sikweli",
	"kama",
	"au",
	"sivyo",
	"wakati",
	"rudisha",
	"vunja",
	"endelea",
	"tupu",
	"ktk",
	"kwa",
	"badili",
	"ikiwa",
	"kawaida",
	"tumia",
	"pakeji",
	"@",
}

var Candidates = new(map[string]uint64)

func defaultCompletionGenerator() (*[]defines.CompletionItem, error) {
	result := make([]defines.CompletionItem, 0)

	funcsKind := defines.CompletionItemKindFunction
	for k, v := range functions {
		result = append(result, defines.CompletionItem{
			Kind:          &funcsKind,
			Label:         k,
			Documentation: v,
		})
	}

	keyWordCompletion := defines.CompletionItemKindKeyword
	for _, v := range keywords {
		completion := defines.CompletionItem{
			Kind:  &keyWordCompletion,
			Label: v,
		}
		result = append(result, completion)
	}

	return &result, nil
}

func CompletionFunc(ctx context.Context,
	req *defines.CompletionParams) (*[]defines.CompletionItem, error) {
	logs.Println("CompletionShow:", req)

	file := string(req.TextDocument.Uri)
	position := req.TextDocumentPositionParams.Position

	defaultCompletion, _ := defaultCompletionGenerator()

	data.PagesMutext.Lock()
	defer data.PagesMutext.Unlock()

	doc, found := data.Pages[file]
	//check if such a doc was already included, if not just skip to do
	//default evaluation with hints
	if found == false {
		return defaultCompletion, nil
	}
	if position.Line >= uint(len(doc.Content)) {
		logs.Printf("Error: position  %d > file %s of lines %d \n",
			position.Line, file, len(doc.Content))
		return defaultCompletion, nil
	}

	//get the word to be completed
	wordToCompelte := ""
	line := doc.Content[position.Line]
	startPosition := position.Character - 1
	for startPosition >= 0 && startPosition < uint(len(line)) {
		//get space symbolizing start of new word
		if line[startPosition] == ' ' {
			startPosition += 1
			break
		}
		//get space symbolizing start of new word after a dot
		if line[startPosition] == '.' {
			startPosition += 1
			break
		}
		if startPosition == 0 {
			break
		}
		startPosition -= 1
	}
	wordToCompelte = line[startPosition:position.Character]

	if len(wordToCompelte) == 0 {
		//return using all keywods
		return defaultCompletion, nil
	} else {
		//filter the data to be sent
		finalData := make([]defines.CompletionItem, 0)
		for _, v := range *defaultCompletion {
			if strings.Contains(v.Label, wordToCompelte) {
				finalData = append(finalData, v)
			}
		}
		return &finalData, nil
	}
}

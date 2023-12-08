package completions

import (
	"context"
	"nuru-lsp/data"
	"strings"

	"github.com/TobiasYin/go-lsp/logs"
	"github.com/TobiasYin/go-lsp/lsp/defines"
)

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
	for _, v := range keywords {
		completion := defines.CompletionItem{
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
	if position.Line > uint(len(doc.Content)) {
		logs.Printf("Error: position  %d > file %s of lines %d \n",
			position.Line, file, len(doc.Content))
		return defaultCompletion, nil
	}

	//get the word to be completed
	wordToCompelte := ""
	for i, v := range doc.Content {
		if i == int(position.Line) {
			startPosition := position.Character - 1
			for startPosition >= 0 {
				//get space symbolizing start of new word
				if v[startPosition] == ' ' {
					startPosition += 1
					break
				}
				//get space symbolizing start of new word after a dot
				if v[startPosition] == '.' {
					startPosition += 1
					break
				}
				startPosition -= 1
			}
			wordToCompelte = v[startPosition:func() uint {
				if position.Character > 0 {
					return position.Character + 1
				} else {
					return position.Character
				}
			}()]
		}
	}

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

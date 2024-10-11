package completions

import (
	"context"
	"errors"
	"nuru-lsp/data"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
)

var Functions = map[string]string{
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

var Keywords = []string{
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

func DefaultCompletionGenerator() (*[]defines.CompletionItem, error) {
	result := make([]defines.CompletionItem, 0)

	funcsKind := defines.CompletionItemKindFunction
	for k, v := range Functions {
		result = append(result, defines.CompletionItem{
			Kind:          &funcsKind,
			Label:         k,
			Documentation: v,
		})
	}

	keyWordCompletion := defines.CompletionItemKindKeyword
	for _, v := range Keywords {
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

	file := string(req.TextDocument.Uri)

	defaultCompletion, _ := DefaultCompletionGenerator()

	data.PagesMutext.Lock()
	defer data.PagesMutext.Unlock()

	logs.Println("FILE IS:",file, "DOCS LENGTH:",len(data.Pages))
	doc, found := data.Pages[file]

	if !found {
		//This should technically never run, as all docs must exist
		return nil, errors.New("DOC NOT FOUND")
	}

	logs.Println("POSITIONS:", req.Position)
	logs.Println("CONTENT:",doc.Content)
	for _,l := range doc.Content{
		logs.Println(l)
	}

	docCompletions, err := doc.Completions(req, defaultCompletion)
	if err!=nil || docCompletions == nil{
		logs.Println("WTF? GANI TENA?", err, docCompletions)
		return defaultCompletion, nil
	}

	return docCompletions, nil
}

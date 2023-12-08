package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/TobiasYin/go-lsp/lsp/defines"
)

type ErrorMapLineNumbers = map[uint64][]string

// Hold information on .nr file
type Data struct {
	File    string
	Version uint64
	Errors  *ErrorMapLineNumbers
	Content []string
}

var Pages = make(map[string]Data)
var PagesMutext = sync.Mutex{}

func NewData(file string, version uint64, content []string) Data {
	return Data{
		File:    file,
		Version: version,
		Errors:  new(ErrorMapLineNumbers),
		Content: content,
	}
}

func OnDataChange(ctx context.Context, req *defines.DidChangeTextDocumentParams) error {
	file := string(req.TextDocument.Uri)

	PagesMutext.Lock()
	defer PagesMutext.Unlock()

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

	Pages[file] = doc

	return nil
}

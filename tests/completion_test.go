package tests

import (
	"nuru-lsp/completions"
	"nuru-lsp/data"
	"runtime"
	"testing"

	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/stretchr/testify/assert"
)

func createCompletionParams(t *testing.T,
	position defines.Position,
	docInput []string, path *string) (data.Data, defines.CompletionParams) {
	var file *string = nil
	if path == nil {
		if _, file_loc, _, ok := runtime.Caller(0); ok == false {
			t.Fatal("Failed to get dir of current package")
		} else {
			file = &file_loc
		}
	} else {
		file = path
	}

	assert.NotNil(t, file)
	assert.NotEqual(t, 0, len(*file))
	t.Logf("File is: %s", *file)

	return data.Data{
			File:    *file,
			Version: 1,
			Errors:  data.ErrorMapLineNumbers{},
			Content: docInput,
		}, defines.CompletionParams{
			TextDocumentPositionParams: defines.TextDocumentPositionParams{
				TextDocument: defines.TextDocumentIdentifier{
					Uri: defines.DocumentUri(*file),
				},
				Position: position,
			},
		}
}

func TestTumiaCompletion(t *testing.T) {
	//create a completions params
	data, completionParams := createCompletionParams(t, defines.Position{
		Line:      0,
		Character: 5,
	}, []string{"tumia "}, nil)

	items, err := data.TreeSitterCompletions(&completionParams)
	assert.Nil(t, err)
	assert.Nil(t, items)
	tumias := append(completions.TUMIAS, "test")

	for _, item := range *items {
		assert.Contains(t, tumias, item.Label)
	}
	assert.Equal(t, nil, err, "TreeSitterCompletions shouldn't return error here")
}

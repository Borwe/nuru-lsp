package tests

import (
	"nuru-lsp/completions"
	data_mod "nuru-lsp/data"
	"runtime"
	"strings"
	"testing"

	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/stretchr/testify/assert"
)

func createCompletionParams(t *testing.T,
	position defines.Position,
	docInput []string, path *string) (data_mod.Data, defines.CompletionParams) {
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

	data, _ := data_mod.NewData(*file, 0, docInput)

	return *data, defines.CompletionParams{
		TextDocumentPositionParams: defines.TextDocumentPositionParams{
			TextDocument: defines.TextDocumentIdentifier{
				Uri: defines.DocumentUri(*file),
			},
			Position: position,
		},
	}
}

func TestTumiaCompletionNoIdentifier(t *testing.T) {
	//create a completions params
	data, completionParams := createCompletionParams(t, defines.Position{
		Line:      0,
		Character: 5,
	}, []string{"tumia "}, nil)

	items, err := data.TreeSitterCompletions(&completionParams)
	assert.Nil(t, err)
	tumias := append(data_mod.TUMIAS, "test", "full_pakeji")

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	for _, item := range tumias {
		assert.Contains(t, itemsLabels, item)
	}
}

func TestTumiaCompletionWithIdentifier(t *testing.T) {
	//create a completions params
	data, completionParams := createCompletionParams(t, defines.Position{
		Line:      0,
		Character: 6,
	}, []string{"tumia t"}, nil)

	items, err := data.TreeSitterCompletions(&completionParams)
	assert.Nil(t, err)

	//fill tumias
	tumias := []string{"test"}
	for _, mod := range data_mod.TUMIAS {
		if strings.Contains(mod, "t") {
			tumias = append(tumias, mod)
		}
	}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Equal(t, len(tumias), len(itemsLabels), "More items in completion than expected")

	for _, item := range tumias {
		assert.Contains(t, itemsLabels, item)
	}
}

func TestVariableFunctionCompletionWithoutIdentifier(t *testing.T) {
	//create a completions params
	data, completionParams := createCompletionParams(t, defines.Position{
		Line:      3,
		Character: 0,
	}, []string{"tumia test",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"yolo = 123",
		"chora = \"50 Cent\"",
		"",
	}, nil)

	items, err := data.TreeSitterCompletions(&completionParams)
	assert.Nil(t, err)

	//fill completions expected
	completions_expected := []string{}
	completions_expected = append(completions_expected, data_mod.TUMIAS...)
	completions_expected = append(completions_expected, completions.Keywords...)
	for k := range completions.Functions {
		completions_expected = append(completions_expected, k)
	}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Equal(t, len(completions_expected), len(itemsLabels), "More items in completion than expected")

	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}

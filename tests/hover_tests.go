package tests

import (
	data_mod "nuru-lsp/data"
	"nuru-lsp/setup"
	"testing"

	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/stretchr/testify/assert"
)


func TestLocalFileHoverOnStdTumias(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 8,
	}, []string{"tumia jsoni"}, nil)

	items, err := data.Completions(&completionParams, nil)
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

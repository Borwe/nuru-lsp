package tests

import (
	"fmt"
	data_mod "nuru-lsp/data"
	"strings"
	"testing"

	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/stretchr/testify/assert"
)

func TestTumiaCompletionNoIdentifier(t *testing.T) {
	//create a completions params
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 6,
	}, []string{"tumia "}, nil)

	items, err := data.Completions(&completionParams)
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
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 7,
	}, []string{"tumia t"}, nil)

	items, err := data.Completions(&completionParams)
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

	fmt.Println("LABELS: ", itemsLabels)
	assert.Equal(t, len(tumias), len(itemsLabels), "More items in completion than expected")

	for _, item := range tumias {
		assert.Contains(t, itemsLabels, item)
	}
}

func TestCompleteTumiaHeaderCompletionToContainNewFileCretedAfter(t *testing.T){
	firstEdit := "test123.nr"
	secondFile := "testhead123.nr"
	secondFilePakejiname := secondFile[0:len(secondFile)-3]
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 7,
	}, []string{"tumia t"}, &firstEdit)
	items, err := data.Completions(&completionParams)
	assert.Nil(t, err)
	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}
	assert.NotContains(t, itemsLabels,secondFilePakejiname)

	//add secondfile
	_, _, errs := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 7,
	}, []string{ fmt.Sprintf(
			"pakeji %s { checka = unda(){ andika (\"HAHA\")}}",
			secondFilePakejiname)}, &secondFile)
	assert.Equal(t, 0, len(errs), errs)
	items, err = data.Completions(&completionParams)
	assert.Nil(t, err)
	itemsLabels = []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}
	assert.Contains(t, itemsLabels,secondFilePakejiname)
}

//func TestVariableFunctionCompletionWithoutIdentifier(t *testing.T) {
//	//create a completions params
//	data, completionParams := createCompletionParams(t, defines.Position{
//		Line:      5,
//		Character: 0,
//	}, []string{"tumia test",
//		"fanya checka = unda(){ andika(\"Yolo\");}",
//		"wewe = unda(){ andika(\"WEWE\");}",
//		"yolo = 123",
//		"chora = \"50 Cent\"",
//		"",
//	}, nil)
//
//	items, err := data.Completions(&completionParams)
//	assert.Nil(t, err)
//
//	//fill completions expected
//	completions_expected := []string{"test", "checka", "wewe", "yolo", "chora"}
//
//	itemsLabels := []string{}
//	for _, item := range *items {
//		itemsLabels = append(itemsLabels, item.Label)
//	}
//
//	assert.Equal(t, len(completions_expected), len(itemsLabels),
//		"Not same number of items in completion to the expected")
//
//	for _, item := range completions_expected {
//		assert.Contains(t, itemsLabels, item)
//	}
//}

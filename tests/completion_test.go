package tests

import (
	"fmt"
	data_mod "nuru-lsp/data"
	"nuru-lsp/setup"
	"strings"
	"testing"

	"github.com/Borwe/go-lsp/logs"
	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/NuruProgramming/Nuru/module"
	"github.com/stretchr/testify/assert"
)

func TestTumiaCompletionNoIdentifier(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 6,
	}, []string{"tumia "}, nil)

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

func TestTumiaCompletionWithIdentifier(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 7,
	}, []string{"tumia t"}, nil)

	items, err := data.Completions(&completionParams, nil)
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

func TestCompleteTumiaHeaderCompletionToContainNewFileCretedAfter(t *testing.T) {
	setup.SetupLog()
	firstEdit := "test123.nr"
	secondFile := "testhead123.nr"
	secondFilePakejiname := secondFile[0 : len(secondFile)-3]
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 7,
	}, []string{"tumia t"}, &firstEdit)
	items, err := data.Completions(&completionParams, nil)
	assert.Nil(t, err)
	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}
	assert.NotContains(t, itemsLabels, secondFilePakejiname)

	//add secondfile
	_, _, errs := CreateCompletionParams(t, defines.Position{
		Line:      0,
		Character: 7,
	}, []string{fmt.Sprintf(
		"pakeji %s { checka = unda(){ andika (\"HAHA\")}}",
		secondFilePakejiname)}, &secondFile)
	assert.Equal(t, 0, len(errs), errs)
	items, err = data.Completions(&completionParams, nil)
	assert.Nil(t, err)
	itemsLabels = []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}
	assert.Contains(t, itemsLabels, secondFilePakejiname)
}

func TestVariableFunctionCompletionWithoutIdentifierOnNewLine(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, errs := CreateCompletionParams(t, defines.Position{
		Line:      5,
		Character: 0,
	}, []string{"tumia test",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"wewe = unda(){ andika(\"WEWE\");}",
		"yolo = 123",
		"chora = \"50 Cent\"",
		"",
	}, nil)
	assert.Equal(t, 0, len(errs))

	items, err := data.Completions(&completionParams, nil)
	assert.Nil(t, err)

	//fill completions expected
	completions_expected := []string{"test", "checka", "wewe", "yolo", "chora"}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Greater(t, len(itemsLabels), 0, items)
	t.Log("ITEMS: ", itemsLabels)
	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}

func TestVariableFunctionCompletionWithoutIdentifierOnNotNewLine(t *testing.T) {
	setup.SetupLog()
	logs.Println("INDENT TESTING STARTED")
	//create a completions params
	data, completionParams, errs := CreateCompletionParams(t, defines.Position{
		Line:      6,
		Character: 4,
	}, []string{"tumia test",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"wewe = unda(){ andika(\"WEWE\");}",
		"yolo = 123",
		"chora = \"50 Cent\"",
		"fanya hotdamn = unda(){",
		"    ",
		"}",
	}, nil)
	assert.Equal(t, 0, len(errs))

	items, err := data.Completions(&completionParams, nil)
	assert.Nil(t, err)

	//fill completions expected
	completions_expected := []string{"test", "checka", "wewe", "yolo", "chora"}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Greater(t, len(itemsLabels), 0)
	t.Log("ITEMS: ", itemsLabels)
	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}

func TestVariableFunctionCompletionWithIdentifierOnNewLine(t *testing.T) {
	setup.SetupLog()
	logs.Println("INDENT TESTING STARTED")
	//create a completions params
	data, completionParams, errs := CreateCompletionParams(t, defines.Position{
		Line:      5,
		Character: 2,
	}, []string{"tumia test",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"wewe = unda(){ andika(\"WEWE\");}",
		"yolo = 123",
		"chora = \"50 Cent\"",
		"ch",
	}, nil)
	assert.Equal(t, 0, len(errs), errs)

	items, err := data.Completions(&completionParams, nil)
	assert.Nil(t, err, err)

	//fill completions expected
	completions_expected := []string{"checka", "chora"}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Greater(t, len(itemsLabels), 0)
	assert.NotContains(t, itemsLabels, "wewe", "wewe does not contain a ch character, so shouldn't be a completion")
	t.Log("ITEMS: ", itemsLabels)
	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}






func TestVariableFunctionCompletionOfStdPackageOnLastLine(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, errs := CreateCompletionParams(t, defines.Position{
		Line:      5,
		Character: 9,
	}, []string{"tumia hisabati",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"wewe = unda(){ andika(\"WEWE\");}",
		"yolo = 123",
		"chora = \"50 Cent\"",
		"hisabati.",
	}, nil)
	assert.Equal(t, 0, len(errs))

	items, err := data.Completions(&completionParams, nil)
	assert.Nil(t, err)

	//fill completions expected
	completions_expected := []string{}
	hisabati_functions, ok := module.Mapper["hisabati"]
	assert.True(t, ok, "Couldn't get histabati module")
	hisabati_consts := module.Constants
	for fn := range hisabati_functions.Functions{
		completions_expected = append(completions_expected, fn)
	}
	for cs := range hisabati_consts {
		completions_expected = append(completions_expected, cs)
	}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Greater(t, len(itemsLabels), 0)
	t.Log("ITEMS: ", itemsLabels)
	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}


func TestVariableFunctionCompletionOfNonStdPackageOnLastLine(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, errs := CreateCompletionParams(t, defines.Position{
		Line:      5,
		Character: 5,
	}, []string{"tumia test",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"wewe = unda(){ andika(\"WEWE\");}",
		"yolo = 123",
		"chora = \"50 Cent\"",
		"test.",
	}, nil)
	assert.Equal(t, 0, len(errs))

	items, err := data.Completions(&completionParams, nil)
	assert.Nil(t, err)

	//fill completions expected
	completions_expected := []string{"yo","cheka"}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Greater(t, len(itemsLabels), 0)
	t.Log("ITEMS: ", itemsLabels)
	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}


func TestVariableFunctionCompletionOfNonStdPackageInside(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      5,
		Character: 9,
	}, []string{"tumia test",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"wewe = unda(){ andika(\"WEWE\");}",
		"yolo = 123",
		"chora = unda(){",
		"    test.",
		"}",
		"bolo = \"sohk\"",
	}, nil)

	items, err := data.Completions(&completionParams, nil)
	assert.Nil(t, err)

	//fill completions expected
	completions_expected := []string{"yo","cheka"}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Greater(t, len(itemsLabels), 0)
	t.Log("ITEMS: ", itemsLabels)
	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}


func TestVariableFunctionCompletionOfNonStdPackageInsideMethodCompletion(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      5,
		Character: 11,
	}, []string{"tumia test",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"wewe = unda(){ andika(\"WEWE\");}",
		"yolo = 123",
		"chora = unda(){",
		"    test.ch",
		"}",
		"bolo = \"sohk\"",
	}, nil)

	items, _ := data.Completions(&completionParams, nil)

	//fill completions expected
	completions_expected := []string{"cheka"}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Equal(t, 1 , len(itemsLabels), itemsLabels)
	t.Log("ITEMS: ", itemsLabels)
	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}


func TestVariableFunctionCompletionOfStdPackageInsideMethodCompletion(t *testing.T) {
	setup.SetupLog()
	//create a completions params
	data, completionParams, _ := CreateCompletionParams(t, defines.Position{
		Line:      5,
		Character: 11,
	}, []string{"tumia jsoni",
		"fanya checka = unda(){ andika(\"Yolo\");}",
		"wewe = unda(){ andika(\"WEWE\");}",
		"yolo = 123",
		"chora = unda(){",
		"    jsoni.enk",
		"}",
		"bolo = \"sohk\"",
	}, nil)

	items, _ := data.Completions(&completionParams, nil)

	//fill completions expected
	completions_expected := []string{"enkodi"}

	itemsLabels := []string{}
	for _, item := range *items {
		itemsLabels = append(itemsLabels, item.Label)
	}

	assert.Equal(t, 1 , len(itemsLabels), itemsLabels)
	t.Log("ITEMS: ", itemsLabels)
	for _, item := range completions_expected {
		assert.Contains(t, itemsLabels, item)
	}
}

package tests

import (
	"errors"
	"net/url"
	data_mod "nuru-lsp/data"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/stretchr/testify/assert"
)

func CreateImaginaryFilePath(file *string) (*string, error) {
	_, file_loc, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New("Failed to get package dir")
	}

	if file == nil {
		return &file_loc, nil
	}

	dir := filepath.Dir(file_loc)
	abspath := path.Join(dir, *file)

	return &abspath, nil
}

func CreateCompletionParams(t *testing.T,
	position defines.Position,
	docInput []string, path *string) (data_mod.Data, defines.CompletionParams, []string) {

	path, err := CreateImaginaryFilePath(path)
	assert.Nil(t, err)

	file := &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(*path),
	}

	assert.NotEqual(t, 0, len(file.Path))
	data, _, errs := data_mod.NewData(file.String(), 0, docInput)

	return *data, defines.CompletionParams{
		TextDocumentPositionParams: defines.TextDocumentPositionParams{
			TextDocument: defines.TextDocumentIdentifier{
				Uri: defines.DocumentUri(file.String()),
			},
			Position: position,
		},
	}, errs
}

func CreateHoverParam(t *testing.T,
	position defines.Position,
	docInput []string, path *string) defines.HoverParams {

	d, _, _ := CreateCompletionParams(t, position, docInput, path)
	return defines.HoverParams{
		TextDocumentPositionParams: defines.TextDocumentPositionParams{
			TextDocument: defines.TextDocumentIdentifier{
				Uri: defines.DocumentUri(d.FileUri),
			},
			Position: position,
		},
	}
}

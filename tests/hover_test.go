package tests

import (
	//data_mod "nuru-lsp/data"
	"nuru-lsp/hovers"
	"nuru-lsp/setup"
	"testing"

	"github.com/Borwe/go-lsp/lsp/defines"
	"github.com/stretchr/testify/assert"
)

func TestLocalFileHoverOnStdTumias(t *testing.T) {
	setup.SetupLog()
	hoverParams := CreateHoverParam(t, defines.Position{
		Line:      0,
		Character: 8,
	}, []string{"tumia jsoni"}, nil)

	jsoniMessage, ok := hovers.StdTumiasInfo["jsoni"]
	assert.True(t,ok)

	hoverResult, err := hovers.GetHover(&hoverParams)
	assert.Nil(t, err)
	assert.NotNil(t, hoverResult)
	n, ok := hoverResult.Contents.(defines.MarkupContent)
	assert.True(t, ok)
	assert.Equal(t, n.Value, jsoniMessage)
}

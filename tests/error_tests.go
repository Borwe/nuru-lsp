package tests

import (
	datamod "nuru-lsp/data"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldFail(t *testing.T) {
	_, err, errs := datamod.NewData("failtest.nr", 0, []string{"sadasdasd"})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(errs))
}

package gentest_test

import (
	"bytes"
	"testing"

	"github.com/oribe1115/gentest"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	buffer := &bytes.Buffer{}
	gentest.SetWriter(buffer)

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, gentest.Analyzer, "a")

	expected := "hoge\n"
	assert.Equal(t, expected, buffer.String())
}

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
	gentest.Analyzer.Flags.Set("offset", "30")

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, gentest.Analyzer, "a")

	expected := "\n\tfunc TestF(){t *testing.T}\n\t"
	assert.Equal(t, expected, buffer.String())
}

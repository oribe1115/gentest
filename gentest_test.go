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
	testdata := analysistest.TestData()

	tests := []struct {
		Label    string
		TestDir  string
		Offset   int
		Expected string
	}{
		{
			Label:   "simple func",
			TestDir: "a",
			Offset:  30,
			Expected: `
func TestF(t *testing.T) {
	tests := []struct{}{}
	for _, test := range tests {
	}
}`,
		},
		{
			Label:   "simple int func",
			TestDir: "a",
			Offset:  142,
			Expected: `
func TestReturnInt(t *testing.T) {
	tests := []struct{}{}
	for _, test := range tests {
	}
}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			gentest.SetWriter(buffer)
			gentest.SetOffset(test.Offset)
			analysistest.Run(t, testdata, gentest.Analyzer, test.TestDir)
			assert.Equal(t, test.Expected, buffer.String())
		})

	}
}
